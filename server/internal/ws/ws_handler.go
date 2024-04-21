package ws

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/internal/user"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub         *Hub
	UserHandler *user.Handler
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type Permission uint

const (
	CanModifyRoles Permission = 1 << iota
	CanKick
	CanBan
)

func checkPermission(handler *Handler, roomId uint64, userId uint64, permission Permission) bool {
	// check whether user is member of room
	var count uint32
	err := handler.hub.Database.GetDB().QueryRow("SELECT COUNT(*) FROM public.room_members WHERE user_id = $1 AND room_id = $2", userId, roomId).Scan(&count)
	if err != nil {
		return false
	}

	// check whether user is owner of room
	var creatorId uint64
	err = handler.hub.Database.GetDB().QueryRow("SELECT creator_id FROM public.rooms WHERE id = $1", roomId).Scan(&creatorId)
	if err != nil {
		return false
	}

	if creatorId == userId {
		return true
	}

	// check whether user has the permission
	var hasPermission bool
	// loop over all the roles the user has
	rows, err := handler.hub.Database.GetDB().Query("SELECT role_id FROM public.role_members WHERE user_id = $1", userId)
	if err != nil {
		log.Fatalf("could not get roles: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roleId uint64
		err = rows.Scan(&roleId)
		if err != nil {
			log.Fatalf("could not scan role: %s", err)
		}

		switch permission {
		case CanModifyRoles:
			err = handler.hub.Database.GetDB().QueryRow("SELECT can_modify_roles FROM public.roles WHERE room_id = $1 AND role_id = $2", roomId, roleId).Scan(&hasPermission)
		case CanKick:
			err = handler.hub.Database.GetDB().QueryRow("SELECT can_kick FROM public.roles WHERE room_id = $1 AND role_id = $2", roomId, roleId).Scan(&hasPermission)
		case CanBan:
			err = handler.hub.Database.GetDB().QueryRow("SELECT can_ban FROM public.roles WHERE room_id = $1 AND role_id = $2", roomId, roleId).Scan(&hasPermission)
		}

		if hasPermission {
			break
		}

		if err != nil {
			log.Fatalf("could not get permission: %s", err)
		}
	}

	return hasPermission
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

func (handler *Handler) CreateRoom(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req CreateRoomRequest
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var insertedId uint64
	err = handler.hub.Database.GetDB().QueryRow("INSERT INTO public.rooms (creator_id, name) VALUES ($1, $2) RETURNING id", user.Id, req.Name).Scan(&insertedId)
	if err != nil {
		log.Fatalf("could not insert room: %s", err)
	}
	handler.hub.Rooms[insertedId] = &Room{
		Id:      uint64(insertedId),
		Clients: make(map[uint64]*Client),
	}

	context.JSON(http.StatusOK, gin.H{
		"id": uint64(insertedId),
	})

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (handler *Handler) JoinRoom(context *gin.Context) {

	token, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		fmt.Println("neco 1")
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(token)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		fmt.Println("neco 2")
		return
	}

	username := user.Username
	clientId := user.Id

	channelIdInt, err := strconv.Atoi(context.Param("roomId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("neco 3")
		return
	}

	channelId := uint64(channelIdInt)

	var creatorId uint64
	var active bool

	var channelRoomId uint64

	err = handler.hub.Database.GetDB().QueryRow("SELECT room_id FROM public.channels WHERE id = $1", channelId).Scan(&channelRoomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("neco 4")
		return
	}

	var isVoice bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT active, creator_id FROM public.rooms WHERE id = $1", channelRoomId).Scan(&active, &creatorId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("neco 5")
		return
	}

	err = handler.hub.Database.GetDB().QueryRow("SELECT is_voice FROM public.channels where id = $1", channelId).Scan(&isVoice)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error with is_voice": err.Error()})
		fmt.Println("neco 6")
	}

	if !active {
		context.JSON(http.StatusBadRequest, gin.H{"error": "room is not active"})
		fmt.Println("neco 7")
		return

	}

	// is banned?
	var isBanned bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM public.bans WHERE room_id = $1 AND user_id = $2)", channelRoomId, user.Id).Scan(&isBanned)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("neco 11")
		return
	}

	if isBanned {
		context.JSON(http.StatusBadRequest, gin.H{"error": "user is banned from the room"})
		fmt.Println("neco 12")
		return
	}

	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("neco 13")
		return
	}

	if !isVoice {
		client := &Client{
			Conn:     conn,
			Message:  make(chan *Message, 10),
			Id:       clientId,
			RoomID:   channelId,
			Username: username,
		}

		message := &Message{
			Content: "A new user has joined the room",
			RoomId:  channelId,
			UserId:  client.Id,
		}

		handler.hub.Register <- client
		handler.hub.Broadcast <- message

		go client.writeMessage()
		client.readMessage(handler.hub)
		fmt.Println("neco 13")
	} else {
		fmt.Println("neco 14")
		go func() {

		}()
	}
}

type RoomRes struct {
	Id uint64 `json:"id"`
}

func (handler *Handler) GetRooms(context *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, room := range handler.hub.Rooms {
		rooms = append(rooms, RoomRes{
			Id: room.Id,
		})
	}

	context.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}

func (handler *Handler) GetClients(context *gin.Context) {
	var clients []ClientRes
	channelIdInt, err := strconv.Atoi(context.Param("roomId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channelId := uint64(channelIdInt)

	if _, ok := handler.hub.Rooms[channelId]; !ok {
		clients = make([]ClientRes, 0)
		context.JSON(http.StatusOK, clients)
	}

	for _, client := range handler.hub.Rooms[channelId].Clients {
		clients = append(clients, ClientRes{
			Id:       client.Id,
			Username: client.Username,
		})
	}

	context.JSON(http.StatusOK, clients)
}

type CreateInviteReq struct {
	RoomId uint64 `json:"roomId"`
}

func (handler *Handler) InviteToRoom(context *gin.Context) {
	var req CreateInviteReq
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var creatorId uint64
	var friendRoom bool
	var active bool
	// invitable
	err = handler.hub.Database.GetDB().QueryRow("SELECT active, friendship_room, creator_id FROM public.rooms WHERE id = $1", req.RoomId).Scan(&active, &friendRoom, &creatorId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !active {
		context.JSON(http.StatusBadRequest, gin.H{"error": "room is not active"})
		return
	}

	if user.Id != creatorId {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you are not the creator of the room"})
		return
	}

	if friendRoom {
		context.JSON(http.StatusBadRequest, gin.H{"error": "room is not invitable"})
		return
	}

	var outId uint64
	err = handler.hub.Database.GetDB().QueryRow("INSERT INTO public.invites (room_id, inviter_id) VALUES ($1, $2) returning id", req.RoomId, user.Id).Scan(&outId)
	if err != nil {
		log.Fatalf("could not insert invite: %s", err)
	}

	context.JSON(http.StatusOK, gin.H{"invite": outId})
}

func (handler *Handler) CreateInvite(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var creatorId uint64
	var friendRoom bool

	_, err = handler.hub.Database.GetDB().Query("SELECT friendship_room, creator_id FROM public.rooms WHERE id = $1", user.Id, &friendRoom, &creatorId)
	if err != nil {
		log.Fatalf("could not get room: %s", err)
	}

	if uint64(user.Id) != creatorId {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you are not the creator of the room"})
		return
	}

	if friendRoom {
		context.JSON(http.StatusBadRequest, gin.H{"error": "room is not invitable"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.invites (room_id, inviter_id) VALUES ($1, $2)", user.Id, user.Id)
	if err != nil {
		log.Fatalf("could not insert invite: %s", err)
	}

	context.JSON(http.StatusOK, gin.H{"message": "invite added"})
}

type CreateFriendReq struct {
	FriendId uint64 `json:"id"`
}

func (handler *Handler) CreateFriend(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var friendId CreateFriendReq
	err = context.ShouldBindJSON(&friendId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if friendId.FriendId == user.Id {
		context.JSON(http.StatusBadRequest, gin.H{"error": "cannot add yourself as a friend"})
		return
	}

	var insertedId uint64
	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.friend_requests (sender_id, receiver_id) VALUES ($1, $2)", user.Id, friendId.FriendId)
	if err != nil {
		log.Fatalf("could not create friend request: %s", err)
	}

	context.JSON(http.StatusOK, gin.H{"id": insertedId})
}

const MAX_UPLOAD_SIZE = 8 * 1024 * 1024 // 8MiB

func (handler *Handler) UploadFile(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["file"]
	if len(files) != 1 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "file length bad"})
		return
	}

	file, err := files[0].Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	if files[0].Size > MAX_UPLOAD_SIZE {
		context.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	data := make([]byte, files[0].Size)
	file.Read(data)

	sumArray := sha256.Sum256(data)
	sum := binary.LittleEndian.Uint32(sumArray[:])

	var count uint32
	err = handler.hub.Database.GetDB().QueryRow("SELECT COUNT(*) FROM public.files WHERE hash = $1", sum).Scan(&count)
	if err != nil {
		log.Fatalf("could not get count: %s", err)
	}
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Println("shit")
	fmt.Printf("count = %v\n", count)
	if count == 0 {
		_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.files (hash) VALUES ($1)", sum)
		if err != nil {
			log.Fatalf("could not insert file: %s", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current working directory: %s", err)
		}
		newFile, err := os.Create(fmt.Sprintf("%v\\files\\%v", cwd, sum))
		if err != nil {
			log.Fatalf("could not create file: %s", err)
		}

		_, err = newFile.Write(data)
		if err != nil {
			log.Fatalf("could not write file: %s", err)
		}

		newFile.Close()
	}

	context.JSON(http.StatusOK, gin.H{"hash": sum})
}

func (handler *Handler) GetUpload(context *gin.Context) {
	hash := context.Param("hash")
	file, err := os.Open("files/" + hash)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	context.File("files/" + hash)
}

func (handler *Handler) UpdateUsername(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var username string
	err = context.ShouldBindJSON(&username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("UPDATE public.users SET username = $1 WHERE id = $2", username, user.Id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, room := range handler.hub.Rooms {
		for _, client := range room.Clients {
			if client.Id == user.Id {
				client.ChangeUsername(username)
			}
		}
	}

	context.JSON(http.StatusOK, gin.H{"message": "username updated"})
}

func (handler *Handler) JoinInvite(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var invite uint64
	err = context.ShouldBindJSON(&invite)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var roomId uint64
	err = handler.hub.Database.GetDB().QueryRow("SELECT room_id FROM public.invites WHERE id = $1", invite).Scan(&roomId)
	if err != nil {
		fmt.Printf("abc\n")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var active bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT active FROM public.rooms WHERE id = $1", roomId).Scan(&active)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !active {
		context.JSON(http.StatusBadRequest, gin.H{"error": "room is not active"})
		return
	}

	var count uint64
	err = handler.hub.Database.GetDB().QueryRow("SELECT COUNT(*) FROM public.bans WHERE user_id = $1 AND room_id = $2", user.Id, roomId).Scan(&count)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if count != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you are banned from the room"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.room_members (user_id, room_id) VALUES ($1, $2)", user.Id, roomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "joined room"})
}

func (handler *Handler) LeaveRoom(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var roomId uint64
	err = context.ShouldBindJSON(&roomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.room_members WHERE user_id = $1 AND room_id = $2", user.Id, roomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "left room"})
}

func (handler *Handler) DeleteRoom(context *gin.Context) {

	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var roomId uint64
	err = context.ShouldBindJSON(&roomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var creatorId uint64
	err = handler.hub.Database.GetDB().QueryRow("SELECT creator_id FROM public.rooms WHERE id = $1", roomId).Scan(&creatorId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if creatorId != user.Id {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you are not the creator of the room"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("UPDATE public.rooms SET active = false WHERE id = $1", roomId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "deleted room"})

}

func (handler *Handler) ShowFriendReqs(context *gin.Context) {
	token, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(token)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return

	}

	var requests []uint64

	rows, err := handler.hub.Database.GetDB().Query("SELECT sender_id FROM public.friend_requests WHERE receiver_id = $1", user.Id)
	if err != nil {
		log.Fatalf("could not get friend requests: %s", err)
	}

	defer rows.Close()

	var senderId uint64
	rows.Scan(&senderId)
	requests = append(requests, senderId)

	for rows.Next() {
		var senderId uint64
		err = rows.Scan(&senderId)
		if err != nil {
			log.Fatalf("could not scan friend request: %s", err)
		}

		requests = append(requests, senderId)
	}

	context.JSON(http.StatusOK, gin.H{"requests": requests})

}

type AcceptFriendReq struct {
	Id uint64 `json:"id"`
}

func (handler *Handler) AcceptFriend(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var friendId AcceptFriendReq
	err = context.ShouldBindJSON(&friendId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	handler.hub.Database.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM public.friend_requests WHERE sender_id = $1 AND receiver_id = $2)", friendId.Id, user.Id).Scan(&exists)
	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{"error": "friend request does not exist"})
		return
	}

	// todo: error checking
	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.friend_requests WHERE sender_id = $1 AND receiver_id = $2", friendId.Id, user.Id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// handler.hub.Database.GetDB().Exec("INSERT INTO public.rooms (creator_id, invitable) VALUES ($1, $2)", friendId.Id, false)
	// handler.hub.Database.GetDB().Exec("INSERT INTO public.room_members (user_id, room_id) VALUES ($1, $2)", user.Id, friendId.Id)
	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.friends (user_id, friend_id) VALUES ($1, $2)", user.Id, friendId.Id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id uint64

	var username string
	err = handler.hub.Database.GetDB().QueryRow("SELECT username FROM public.users WHERE id = $1", friendId.Id).Scan(&username)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room_name := fmt.Sprintf("%v-%v", user.Username, username)

	err = handler.hub.Database.GetDB().QueryRow("INSERT INTO public.rooms (creator_id, name, friendship_room) VALUES ($1, $2, $3) RETURNING id", user.Id, room_name, true).Scan(&id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var voice_name string = fmt.Sprintf("%v-%v voice", user.Username, username)
	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.channels (name, room_id, is_voice) VALUES ($1, $2, $3)", voice_name, id, true)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "accepted friend request"})
}

type RejectFriendReq struct {
	Id uint64 `json:"id"`
}

func (handler *Handler) RejectFriend(context *gin.Context) {
	var friendId RejectFriendReq
	err := context.ShouldBindJSON(&friendId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var exists bool
	handler.hub.Database.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM public.friend_requests WHERE sender_id = $1 AND receiver_id = $2)", friendId.Id, user.Id).Scan(&exists)
	if !exists {
		context.JSON(http.StatusBadRequest, gin.H{"error": "friend request does not exist"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.friend_requests WHERE sender_id = $1 AND receiver_id = $2", friendId.Id, user.Id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "rejected friend request"})
}

func (handler *Handler) GetFriends(context *gin.Context) {
	token, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(token)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var friends []uint64

	rows, err := handler.hub.Database.GetDB().Query("SELECT friend_id FROM public.friends WHERE user_id = $1", user.Id)
	if err != nil {
		log.Fatalf("could not get friends: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var friendId uint64
		err = rows.Scan(&friendId)
		if err != nil {
			log.Fatalf("could not scan friend: %s", err)
		}

		friends = append(friends, friendId)
	}

	rows, err = handler.hub.Database.GetDB().Query("SELECT user_id FROM public.friends WHERE friend_id = $1", user.Id)
	if err != nil {
		log.Fatalf("could not get friends: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var friendId uint64
		err = rows.Scan(&friendId)
		if err != nil {
			log.Fatalf("could not scan friend: %s", err)
		}
		friends = append(friends, friendId)
	}

	uniqueFriends := make([]uint64, 0)
	friendMap := make(map[uint64]bool)

	for _, id := range friends {
		if _, exists := friendMap[id]; !exists {
			friendMap[id] = true
			uniqueFriends = append(uniqueFriends, id)
		}
	}

	context.JSON(http.StatusOK, gin.H{"friends": uniqueFriends})

}

type MemberInteractionRequest struct {
	UserId uint64 `json:"userId"`
	RoomId uint64 `json:"roomId"`
}

func (handler *Handler) Kick(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req MemberInteractionRequest
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId := req.RoomId

	var friendRoom bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT friendship_room FROM public.rooms WHERE id = $1", roomId).Scan(&friendRoom)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if friendRoom {
		context.JSON(http.StatusBadRequest, gin.H{"error": "cannot kick from friendship room"})
		return
	}

	if !checkPermission(handler, uint64(roomId), user.Id, CanKick) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to kick users"})
		return
	}

	if req.UserId == user.Id {
		context.JSON(http.StatusBadRequest, gin.H{"error": "cannot kick yourself"})
		return
	}

	handler.hub.Database.GetDB().Exec("DELETE FROM public.room_members WHERE user_id = $1 AND room_id = $2", req.UserId, roomId)
}

func (handler *Handler) Ban(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req MemberInteractionRequest
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId := req.RoomId

	var friendRoom bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT friendship_room FROM public.rooms WHERE id = $1", roomId).Scan(&friendRoom)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if friendRoom {
		context.JSON(http.StatusBadRequest, gin.H{"error": "cannot kick from friendship room"})
		return
	}

	if !checkPermission(handler, uint64(roomId), user.Id, CanBan) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to ban users"})
		return
	}
	handler.hub.Database.GetDB().Exec("INSERT INTO public.bans (room_id, user_id) VALUES ($1, $2)", roomId, req.UserId)
	handler.hub.Database.GetDB().Exec("DELETE FROM public.room_members WHERE user_id = $1 AND room_id = $2", req.UserId, roomId)
}

func (handler *Handler) Unban(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req MemberInteractionRequest
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId := req.RoomId

	if !checkPermission(handler, uint64(roomId), user.Id, CanBan) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to unban users"})
		return
	}

	var userId uint64
	err = context.ShouldBindJSON(&userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.hub.Database.GetDB().Exec("DELETE FROM public.bans WHERE room_id = $1 AND user_id = $2", roomId, userId)
}

// role changes
type CreateRoleRequest struct {
	RoomId uint   `json:"roomId"`
	Name   string `json:"name"`
}

func (handler *Handler) CreateRole(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req CreateRoleRequest
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var friendRoom bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT friendship_room FROM public.rooms WHERE id = $1", req.RoomId).Scan(&friendRoom)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if friendRoom {
		context.JSON(http.StatusBadRequest, gin.H{"error": "cannot kick from friendship room"})
		return
	}

	if !checkPermission(handler, uint64(req.RoomId), user.Id, CanModifyRoles) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to create roles"})
		return
	}

	var insertedId uint64
	err = handler.hub.Database.GetDB().QueryRow("INSERT INTO public.roles (name, room_id) VALUES ($1, $2) RETURNING id", req.Name, req.RoomId).Scan(&insertedId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": insertedId})
}

type DeleteRoleReq struct {
	RoomId uint `json:"roomId"`
	RoleId uint `json:"roleId"`
}

func (handler *Handler) DeleteRole(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req DeleteRoleReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkPermission(handler, uint64(req.RoomId), user.Id, CanModifyRoles) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to delete roles"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.roles WHERE id = $1", req.RoleId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "role deleted"})
}

type RoleReq struct {
	RoomId uint64 `json:"roomId"`
	RoleId uint64 `json:"roleId"`
	UserId uint64 `json:"userId"`
}

func (handler *Handler) AddUserToRole(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req RoleReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkPermission(handler, req.RoomId, user.Id, CanModifyRoles) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to add users to roles"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("INSERT INTO public.role_members (user_id, role_id) VALUES ($1, $2)", req.UserId, req.RoleId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user added to role"})

}

func (handler *Handler) RemoveUserFromRole(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req RoleReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkPermission(handler, req.RoomId, user.Id, CanModifyRoles) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to remove users from roles"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.role_members WHERE user_id = $1 AND role_id = $2", req.UserId, req.RoleId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user removed from role"})
}

type EditRoleReq struct {
	RoomId          uint64 `json:"roomId"`
	RoleId          uint64 `json:"roleId"`
	BanUsersEdit    bool   `json:"banUsersEdit"`
	BanUsers        bool   `json:"banUsers"`
	KickUsersEdit   bool   `json:"kickUsersEdit"`
	KickUsers       bool   `json:"kickUsers"`
	ModifyRolesEdit bool   `json:"modifyRolesEdit"`
	ModifyRoles     bool   `json:"modifyRoles"`
}

func (handler *Handler) ModifyRolePermissions(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	var req EditRoleReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !checkPermission(handler, req.RoomId, user.Id, CanModifyRoles) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "you do not have permission to modify roles"})
		return
	}

	if req.BanUsersEdit {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.roles SET can_ban = $1 WHERE id = $2", req.BanUsers, req.RoleId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if req.KickUsersEdit {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.roles SET can_kick = $1 WHERE id = $2", req.KickUsers, req.RoleId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if req.ModifyRolesEdit {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.roles SET can_modify_roles = $1 WHERE id = $2", req.ModifyRoles, req.RoleId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{"message": "role permissions updated"})
}

func (handler *Handler) GetUser(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error ": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"error": "invalid token"})
		return
	}

	userId, err := strconv.Atoi(context.Param("userId"))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	var username string
	var displayName string
	var bio string

	err = handler.hub.Database.GetDB().QueryRow("SELECT username, display_name, bio FROM public.users WHERE id = $1", userId).Scan(&username, &displayName, &bio)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"username": username, "displayName": displayName, "bio": bio})
}

type UpdateUserReq struct {
	UpdateDisplayName bool   `json:"updateDisplayName"`
	DisplayName       string `json:"displayName"`
	UpdateBio         bool   `json:"updateBio"`
	Bio               string `json:"bio"`
	UpdateThemeId     bool   `json:"updateThemeId"`
	ThemeId           int    `json:"themeId"`
}

func (handler *Handler) UpdateSelfUser(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error ": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"error": "invalid token"})
		return
	}

	var req UpdateUserReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	if req.UpdateDisplayName {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.users SET display_name = $1 WHERE id = $2", req.DisplayName, user.Id)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	if req.UpdateBio {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.users SET bio = $1 WHERE id = $2", req.Bio, user.Id)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	if req.UpdateThemeId {
		_, err = handler.hub.Database.GetDB().Exec("UPDATE public.users SET theme_id = $1 WHERE id = $2", req.ThemeId, user.Id)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}
	context.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

type CreateChannelReq struct {
	RoomId  uint   `json:"roomId"`
	Name    string `json:"name"`
	IsVoice bool   `json:"isVoice"`
}

func (handler *Handler) CreateChannel(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"error": "invalid token"})
		return
	}

	var req CreateChannelReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	var friendRoom bool
	err = handler.hub.Database.GetDB().QueryRow("SELECT friendship_room FROM public.rooms WHERE id = $1", req.RoomId).Scan(&friendRoom)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	if friendRoom {
		context.JSON(http.StatusOK, gin.H{"error": "cannot kick from friendship room"})
		return
	}

	if !checkPermission(handler, uint64(req.RoomId), user.Id, CanModifyRoles) {
		context.JSON(http.StatusOK, gin.H{"error": "you do not have permission to create channels"})
		return
	}

	var channelName = req.Name
	// err = context.ShouldBindJSON(&channelName)
	// if err != nil {
	// 	context.JSON(http.StatusOK, gin.H{"error": err.Error()})
	// 	return
	// }

	var insertedId uint64
	err = handler.hub.Database.GetDB().QueryRow("INSERT INTO public.channels (name, room_id, is_voice) VALUES ($1, $2, $3) RETURNING id", channelName, req.RoomId, req.IsVoice).Scan(&insertedId)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": insertedId})
}

type DeleteChannelReq struct {
	RoomId    uint64 `json:"roomId"`
	ChannelId uint64 `json:"channelId"`
}

func (handler *Handler) DeleteChannel(context *gin.Context) {
	jwtToken, err := context.Cookie("jwt")
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": "no jwt token"})
		return
	}

	user, err := handler.UserHandler.GetUserByJWT(jwtToken)
	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"error": "invalid token"})
		return
	}

	var req DeleteChannelReq
	err = context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	if !checkPermission(handler, req.RoomId, user.Id, CanModifyRoles) {
		context.JSON(http.StatusOK, gin.H{"error": "you do not have permission to delete channels"})
		return
	}

	_, err = handler.hub.Database.GetDB().Exec("DELETE FROM public.channels WHERE id = $1", req.ChannelId)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "channel deleted"})
}

type Channel struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	IsVoice bool   `json:"isVoice"`
}

func (handler *Handler) GetChannels(context *gin.Context) {

	roomId, err := strconv.Atoi(context.Param("roomId"))
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	channels := make([]Channel, 1)
	rows, err := handler.hub.Database.GetDB().Query("SELECT id, name, is_voice FROM public.channels WHERE room_id = $1", roomId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var channel Channel
		err = rows.Scan(&channel.Id, &channel.Name, &channel.IsVoice)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		channels = append(channels, channel)
	}

	context.JSON(http.StatusOK, channels)

}

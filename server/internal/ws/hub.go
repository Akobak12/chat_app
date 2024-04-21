package ws

import (
	"fmt"
	"log"
	"server/db"

	"github.com/gorilla/websocket"
)

type Room struct {
	Id           uint64             `json:"id"`
	Clients      map[uint64]*Client `json:"clients"`
	AllowedUsers []uint64           `json:"allowed_users"`
}

type Hub struct {
	Rooms      map[uint64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Database   *db.Database
}

func NewHub(dbConnection *db.Database) *Hub {
	return &Hub{
		Rooms:      make(map[uint64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
		Database:   dbConnection,
	}
}

func (hub *Hub) JoinVoiceRoom(client *Client) bool {
	room, hasRoom := hub.Rooms[client.RoomID]
	if !hasRoom {
		return false
	}

	_, hasClient := room.Clients[client.Id]
	if hasClient {
		return false
	}

	room.Clients[client.Id] = client

	waiter := make(chan bool)

	go func() {
		for {
			messageType, bytes, err := client.Conn.ReadMessage()
			if err != nil {
				log.Printf("error: %v", err)
				break
			}
			if messageType != websocket.BinaryMessage {
				continue
			}

			for _, c := range room.Clients {
				if c.Id != client.Id {
					err := c.Conn.WriteMessage(websocket.BinaryMessage, bytes)
					if err != nil {
						log.Printf("error: %v", err)
						break
					}
				}
			}

		}
	}()

	<-waiter

	return true
}

func (hub *Hub) registerClient(client *Client) {
	room, hasRoom := hub.Rooms[client.RoomID]
	if !hasRoom {
		return
	}

	_, hasClient := room.Clients[client.Id]
	if !hasClient {
		room.Clients[client.Id] = client
	}

	messages, err := hub.Database.GetDB().Query("SELECT content, sender_id FROM public.messages WHERE channel_id = $1", client.RoomID)
	if err != nil {
		log.Fatalf("could not get messages: %s", err)
	}

	for messages.Next() {
		var content string
		var userID uint64
		err = messages.Scan(&content, &userID)
		if err != nil {
			log.Fatalf("could not scan messages: %s", err)
		}

		fmt.Printf("content: %s, userID: %d\n", content, userID)

		client.Conn.WriteJSON(&Message{
			Content: content,
			UserId:  userID,
			RoomId:  client.RoomID,
		})
	}

}

func (hub *Hub) unregisterClient(client *Client) {
	room, hasRoom := hub.Rooms[client.RoomID]
	if !hasRoom {
		return
	}
	_, hasClient := room.Clients[client.Id]
	if !hasClient {
		return
	}
	// todo: try remove if check, mostly unndeed??
	if len(hub.Rooms[client.RoomID].Clients) != 0 {
		hub.Broadcast <- &Message{
			Content: "user left the chat",
			RoomId:  client.RoomID,
			UserId:  client.Id,
		}
	}

	delete(hub.Rooms[client.RoomID].Clients, client.Id)
	close(client.Message)
}

func (hub *Hub) broadcast(message *Message) {
	_, hasRoom := hub.Rooms[message.RoomId]
	if !hasRoom {
		return
	}
	for _, client := range hub.Rooms[message.RoomId].Clients {
		client.Message <- message
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.registerClient(client)
		case client := <-hub.Unregister:
			hub.unregisterClient(client)
		case message := <-hub.Broadcast:
			hub.broadcast(message)
		}
	}
}

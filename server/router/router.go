package router

import (
	"net/http"
	"server/internal/user"
	"server/internal/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func RequireNoAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("jwt")
		if err != nil || cookie == "" {
			context.Next()
		} else {
			context.JSON(http.StatusForbidden, gin.H{"error": "You are already logged in"})
			context.Abort()
		}
	}
}

func RequireAuthMiddleware(svc user.Service) gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("jwt")
		if err != nil || cookie == "" {
			context.JSON(http.StatusForbidden, gin.H{cookie: "You must be logged in"})
			context.Abort()
			return
		}
		user, err := svc.GetUserByJWT(cookie)
		if err != nil || user == nil {
			context.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			context.Abort()
			return
		}
		context.Next()
	}
}

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	router = gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authMiddleware := RequireAuthMiddleware(userHandler.Service)

	wsHandler.UserHandler = userHandler

	// body: { email: string, username: string, password: string }
	router.POST("/api/register", RequireNoAuthMiddleware(), userHandler.CreateUser)
	// body: { email: string, password: string }
	router.POST("/api/login", RequireNoAuthMiddleware(), userHandler.Login)
	// body: none
	router.GET("/api/logout", authMiddleware, userHandler.Logout)

	// body: { name: string (room name) }
	router.POST("/api/ws/create-room", authMiddleware, wsHandler.CreateRoom)
	// body: uint - invite id
	router.POST("/api/ws/join-invite", authMiddleware, wsHandler.JoinInvite)
	// body: uint64 - user id
	router.POST("/api/ws/create-friend", authMiddleware, wsHandler.CreateFriend)

	router.GET("/api/ws/get-friend-requests", authMiddleware, wsHandler.ShowFriendReqs)
	router.GET("/api/ws/get-friends", authMiddleware, wsHandler.GetFriends)
	// body: { roomId: uint64 (room id) }
	router.POST("/api/ws/invite-to-room", authMiddleware, wsHandler.InviteToRoom)

	// body: none, connects to websocket
	router.GET("/api/ws/join-room/:roomId", authMiddleware, wsHandler.JoinRoom)

	// body: none
	router.GET("/api/ws/get-rooms", authMiddleware, wsHandler.GetRooms)
	router.GET("/api/ws/get-clients/:roomId", authMiddleware, wsHandler.GetClients)

	// multiform - "file"
	router.POST("/api/ws/upload-file", authMiddleware, wsHandler.UploadFile)
	// body: none
	router.GET("/api/upload/:hash", wsHandler.GetUpload)

	// body: string - username
	router.POST("/api/ws/update-username", authMiddleware, wsHandler.UpdateUsername)

	// body: uint - friend id
	router.POST("/api/ws/accept-friend", authMiddleware, wsHandler.AcceptFriend)

	router.POST("/api/ws/reject-friend", authMiddleware, wsHandler.RejectFriend)

	// body: { userId: uint, roomId: uint }
	router.DELETE("/api/ws/kick", authMiddleware, wsHandler.Kick)
	// body: uint - room id
	router.DELETE("/api/ws/leave-room", authMiddleware, wsHandler.LeaveRoom)
	// body: uint - room id
	router.DELETE("/api/ws/delete-room", authMiddleware, wsHandler.DeleteRoom)

	// body: { userId: uint, roomId: uint }
	router.POST("/api/ws/ban", authMiddleware, wsHandler.Ban)
	// body: { userId: uint, roomId: uint }
	router.POST("/api/ws/unban", authMiddleware, wsHandler.Unban)

	// body: { name: string, roomId: uint }
	router.POST("/api/ws/create-role", authMiddleware, wsHandler.CreateRole)
	// body: { roleId: uint, roomId: uint }
	router.POST("/api/ws/delete-role", authMiddleware, wsHandler.DeleteRole)

	// body: { userId: uint, roleId: uint, roomId: uint }
	router.POST("/api/ws/assign-role", authMiddleware, wsHandler.AddUserToRole)
	// body: { userId: uint, roleId: uint, roomId: uint }
	router.POST("/api/ws/remove-role", authMiddleware, wsHandler.RemoveUserFromRole)

	// edits determine whether to edit a permission or not
	// body: { roleId: uint, roomId: uint, banUsersEdit: bool, banUsers: bool, kickUsersEdit: bool, kickUsers: bool, modifyRolesEdit: bool, modifyRoles: bool }
	router.POST("/api/ws/modify-role", authMiddleware, wsHandler.ModifyRolePermissions)

	// body: none
	router.GET("/api/ws/get-user/:userId", authMiddleware, wsHandler.GetUser)
	// updates determine whether to update a field or not
	// body: { updateDisplayName: bool, displayName: string, updateBio: bool, bio: string, updateThemeId: bool, themeId: int }
	router.POST("/api/ws/update-self", authMiddleware, wsHandler.UpdateSelfUser)

	// body: { name: string, roomId: uint, isVoice: bool }
	router.POST("/api/ws/create-channel", authMiddleware, wsHandler.CreateChannel)
	// body: { channelId: uint, roomId: uint }
	router.POST("/api/ws/delete-channel", authMiddleware, wsHandler.DeleteChannel)

	// body: none
	router.GET("/api/ws/get-channels/:roomId", authMiddleware, wsHandler.GetChannels)
}

func Start(addr string) error {
	return router.Run(addr)
}

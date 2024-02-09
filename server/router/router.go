package router

import (
	"net/http"
	"server/internal/user"
	"server/internal/ws"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt")
		if err != nil || cookie == "" {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are already logged in"})
			c.Abort()
		}
	}
}
func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/api/register", userHandler.CreateUser)
	r.POST("/api/login", AuthMiddleware(), userHandler.Login)
	r.GET("/api/logout", userHandler.Logout)

	r.POST("/api/ws/create-room", wsHandler.CreateRoom)
	r.GET("/api/ws/join-room/:roomId", wsHandler.JoinRoom)
	r.GET("/api/ws/get-rooms", wsHandler.GetRooms)
	r.GET("/api/ws/get-clients/:roomId", wsHandler.GetClients)
}

func Start(addr string) error {
	return r.Run(addr)
}

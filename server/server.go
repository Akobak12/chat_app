package main

import (
	_ "chat-server/helpers"
	"chat-server/login"
	"chat-server/middleware"
	"chat-server/websocket"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

// insane code

func main() {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.POST("/login", login.LoginHandler)

	cleanup := websocket.Init(r)

	defer cleanup()
	log.Println("Server started on :3030")
	if err := r.Run("127.0.0.1:3030"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

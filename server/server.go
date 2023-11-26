package main

import (
	_ "chat-server/helpers"
	"chat-server/websocket"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// insane code

func main() {
	r := gin.Default()
	cleanup := websocket.Init(r)
	defer cleanup()
	log.Println("Server started on :3030")
	if err := r.Run("127.0.0.1:3030"); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

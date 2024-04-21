package ws

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	Id       uint64 `json:"id"`
	RoomID   uint64 `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content string `json:"content"`
	RoomId  uint64 `json:"roomId"`
	UserId  uint64 `json:"userId"`
}

func (client *Client) writeMessage() {
	defer client.Conn.Close()

	for {
		message, ok := <-client.Message
		if !ok {
			return
		}

		client.Conn.WriteJSON(message)
	}
}

func (client *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()
	first := true
	for {
		_, m, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		fmt.Printf("message: %s\n", m)

		msg := &Message{
			Content: string(m),
			RoomId:  client.RoomID,
			UserId:  client.Id,
		}

		// Get the first channel_id of the room
		var channelID uint64
		err = hub.Database.GetDB().QueryRow("SELECT id FROM public.channels WHERE room_id = $1", client.RoomID).Scan(&channelID)
		if err != nil {
			log.Printf("could not retrieve channel ID: %s", err)
			continue // Skip this iteration if we cannot find a channel ID
		}

		// Insert the message into public.messages table using the retrieved channel_id
		_, err = hub.Database.GetDB().Exec("INSERT INTO public.messages (content, channel_id, sender_id) VALUES ($1, $2, $3)", msg.Content, channelID, client.Id)
		if err != nil {
			log.Fatalf("could not insert message: %s", err)
		}

		if first {
			first = false
			var messages []Message
			rows, err := hub.Database.GetDB().Query("SELECT content, sender_id FROM public.messages WHERE channel_id = $1", channelID)
			if err != nil {
				log.Fatalf("could not retrieve messages: %s", err)
			}
			for rows.Next() {
				var content string
				var senderID uint64
				err = rows.Scan(&content, &senderID)
				if err != nil {
					log.Fatalf("could not scan messages: %s", err)
				}
				messages = append(messages, Message{Content: content, RoomId: client.RoomID, UserId: senderID})
			}
			rows.Close()
			for _, m := range messages {
				hub.Broadcast <- &m
			}
		}

		hub.Broadcast <- msg
	}

}

func (client *Client) ChangeUsername(username string) {
	client.Username = username
}

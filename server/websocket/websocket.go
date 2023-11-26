package websocket

import (
	"chat-server/helpers"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type Message struct {
	ID      int
	Content string
	Date    time.Time
	Sender  string
}

type Client struct {
	ID       int
	Nickname string
	Conn     *websocket.Conn
}

type Session struct {
	ID       int
	Clients  map[int]*Client
	Messages []Message
	mu       sync.Mutex
}

// only one session cause this thing doesnt support multiple sessions
type SessionManager struct {
	GlobalSession *Session
	db            *sql.DB
}

// constructor for session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		GlobalSession: &Session{
			ID:      1,
			Clients: make(map[int]*Client),
			Messages: []Message{
				{
					ID:      1,
					Content: "Session started.",
					Date:    time.Now(),
					Sender:  "SERVER",
				},
			},
		},
	}
}

func (sm *SessionManager) setup() {
	//set stuff up from the db
	dsn := "postgres://kryzanek_samuel_64d3f_z9sdt:FQ5lHUujGng5YyABdJ9lyH6tCBZTsHE2@hosting.ssps.cajthaml.eu:3337/kryzanek_samuel_64d3f_z9sdt_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	sm.db = db

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	//why am i doing this? i have no idea just to be safe ig
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--cleardb" {
		helpers.ClearTableData(sm.db, "public.clients", "public.messages", "public.sessions")
	}
	_, err = sm.db.Exec(`
		CREATE TABLE IF NOT EXISTS public.sessions
		(
		    id SERIAL PRIMARY KEY,
		    session_name TEXT NOT NULL,
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating sessions table: ", err)
	}

	_, err = sm.db.Exec(`
		CREATE TABLE IF NOT EXISTS public.clients
		(
		    id SERIAL PRIMARY KEY,
		    session_id INTEGER NOT NULL REFERENCES public.sessions(id) ON DELETE CASCADE,
		    username TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Error creating clients table: ", err)
	}
	_, err = sm.db.Exec(`
		INSERT INTO public.sessions (id, session_name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING; 
	`, sm.GlobalSession.ID, "Default Session Name")
	if err != nil {
		log.Fatal("Error inserting default session: ", err)
	}

	_, err = sm.db.Exec(`
		CREATE TABLE IF NOT EXISTS public.messages
		(
		    id SERIAL PRIMARY KEY,
		    session_id INTEGER NOT NULL REFERENCES public.sessions(id) ON DELETE CASCADE,
		    sender_id INTEGER NOT NULL REFERENCES public.clients(id) ON DELETE CASCADE,
		    content TEXT NOT NULL,
		    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating messages table: ", err)
	}

	rows, err := sm.db.Query(`SELECT id, content, timestamp, sender_id FROM public.messages WHERE session_id = $1`, sm.GlobalSession.ID)
	if err != nil {
		log.Fatal("Error loading messages: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		var senderID int
		err := rows.Scan(&msg.ID, &msg.Content, &msg.Date, &senderID)
		if err != nil {
			log.Fatal("Error scanning messages: ", err)
		}
		sm.GlobalSession.Messages = append(sm.GlobalSession.Messages, msg)
	}

	log.Println("All tables created successfully.")
}

func (sm *SessionManager) AddClient(conn *websocket.Conn) {
	client := &Client{
		Nickname: "User" + strconv.Itoa(len(sm.GlobalSession.Clients)+1),
		Conn:     conn,
	}

	err := sm.db.QueryRow(`
		INSERT INTO public.clients (session_id, username)
		VALUES ($1, $2) RETURNING id
	`, sm.GlobalSession.ID, client.Nickname).Scan(&client.ID)
	if err != nil {
		log.Printf("Error inserting new client: %v", err)
		return
	}

	sm.GlobalSession.mu.Lock()
	defer sm.GlobalSession.mu.Unlock()

	sm.GlobalSession.Clients[client.ID] = client
	for _, msg := range sm.GlobalSession.Messages {
		message, err2 := sm.parse(&msg)
		if err2 != nil {
			log.Printf("Error parsing message: %v", err2)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("write welcome message error: %v", err)
			return
		}
	}

	go sm.handleMessages(client)
}
func (sm *SessionManager) getMessageSender(messageID int) (string, error) {
	var senderName string
	err := sm.db.QueryRow(`
        SELECT c.username FROM public.clients c
        JOIN public.messages m ON c.id = m.sender_id
        WHERE m.id = $1
    `, messageID).Scan(&senderName)
	if err != nil {
		return "", err
	}
	return senderName, nil
}
func (sm *SessionManager) parse(message *Message) ([]byte, error) {
	// scuffed asl but it works
	if message.Sender == "" {
		senderName, err := sm.getMessageSender(message.ID)
		if err != nil {
			return nil, err
		}
		message.Sender = senderName
	}
	return []byte(strconv.Itoa(message.ID) + " " + message.Date.Format("2006-01-02 15:04:05") + " " + message.Sender + ": " + message.Content), nil
}

// removes client from session
func (sm *SessionManager) RemoveClient(clientID int) {
	sm.GlobalSession.mu.Lock()
	defer sm.GlobalSession.mu.Unlock()

	if client, ok := sm.GlobalSession.Clients[clientID]; ok {
		client.Conn.Close()
		delete(sm.GlobalSession.Clients, clientID)
	}
}

// reads msgs from clients and broadcasts to all clients
func (sm *SessionManager) handleMessages(client *Client) {
	defer func() {
		sm.RemoveClient(client.ID)
	}()
	for {
		_, messageBytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}

		sm.GlobalSession.mu.Lock()
		content := string(messageBytes)
		var msgID int
		err = sm.db.QueryRow(`
			INSERT INTO public.messages (session_id, sender_id, content)
			VALUES ($1, $2, $3) RETURNING id
		`, sm.GlobalSession.ID, client.ID, content).Scan(&msgID)
		if err != nil {
			log.Printf("Error inserting message: %v", err)
			sm.GlobalSession.mu.Unlock()
			continue
		}

		msg := Message{
			ID:      msgID,
			Content: content,
			Date:    time.Now(),
			Sender:  client.Nickname,
		}

		sm.GlobalSession.Messages = append(sm.GlobalSession.Messages, msg)
		sm.GlobalSession.mu.Unlock()

		sm.broadcastToSession(msg)
	}
}

// broadcasts message to all clients
func (sm *SessionManager) broadcastToSession(message Message) {
	sm.GlobalSession.mu.Lock()
	defer sm.GlobalSession.mu.Unlock()

	formattedMessage, err := sm.parse(&message)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	for _, client := range sm.GlobalSession.Clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(formattedMessage)); err != nil {
			log.Printf("write error: %v", err)
			sm.RemoveClient(client.ID)
		}
	}
}

func setUpWebSocketRoute(sm *SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		sm.AddClient(conn)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// quirky and weird but it works
func Init(r *gin.Engine) func() {
	sessionManager := NewSessionManager()
	sessionManager.setup()

	r.GET("/ws", setUpWebSocketRoute(sessionManager))
	return func() {
		sessionManager.db.Close()
	}
}

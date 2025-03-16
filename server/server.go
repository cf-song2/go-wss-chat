package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Upgrade Websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Struct: Client
type Client struct {
	conn *websocket.Conn
	room string
}

var rooms sync.Map
var broadcast = make(chan Message)

// Struct: Message
type Message struct {
	Room    string `json:"room"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Type    string `json:"type"` // "message" | "join" | "leave" | "ping"
}

func init() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Fail to create log file:", err)
	}
	log.SetOutput(file)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Fail to upgrade connection:", err)
		return
	}

	// Get inital message
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Fail to read initial message:", err)
		conn.Close()
		return
	}

	var initMessage Message
	if err := json.Unmarshal(msg, &initMessage); err != nil {
		log.Println("Fail to unmarshal initial message:", err)
		conn.Close()
		return
	}

	room := initMessage.Room
	sender := initMessage.Sender
	client := &Client{conn: conn, room: room}

	clients, _ := rooms.LoadOrStore(room, &sync.Map{})
	clients.(*sync.Map).Store(conn, client)

	log.Printf("[%s] %s is connected.", room, sender)
	broadcast <- Message{Room: room, Sender: sender, Content: "connected.", Type: "join"}

	go keepAlive(conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[%s] %s is disconnected: %v", room, sender, err)
			clients.(*sync.Map).Delete(conn)
			broadcast <- Message{Room: room, Sender: sender, Content: "disconnected.", Type: "leave"}
			conn.Close()
			break
		}

		broadcast <- Message{Room: room, Sender: sender, Content: string(msg), Type: "message"}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		if clients, ok := rooms.Load(msg.Room); ok {
			clients.(*sync.Map).Range(func(key, value interface{}) bool {
				client := value.(*Client)
				err := client.conn.WriteJSON(msg)
				if err != nil {
					log.Printf("[%s] %s failed to send message: %v", msg.Room, msg.Sender, err)
					client.conn.Close()
					clients.(*sync.Map).Delete(client.conn)
				}
				return true
			})
		}
	}
}

func keepAlive(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println("Failed to send ping message:", err)
			conn.Close()
			break
		}
	}
}

func main() {
	err := godotenv.Load(filepath.Join("server", ".env"))
	if err != nil {
		log.Fatal("Faled to load .env file:", err)
	}

	serverIP := os.Getenv("SERVER_IP")
	addrPort := os.Getenv("ADDR_PORT")
	serverAddr := fmt.Sprintf("%s:%s", serverIP, addrPort)

	tcpListener, err := net.Listen("tcp4", serverAddr)
	if err != nil {
		log.Fatalf("Failed to listen on: %v", err)
	}

	certPath := filepath.Join("cert", "fullchain.pem")
	keyPath := filepath.Join("cert", "server.key")

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.HandleFunc("/ws", handleConnections)

	go handleMessages()

	server := &http.Server{
		Addr:      serverAddr,
		Handler:   mux,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	fmt.Println("Start server at:", "wss://"+serverAddr+"/ws")
	err = server.ServeTLS(tcpListener, certPath, keyPath)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var mutex = sync.Mutex{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 업그레이드 실패:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("클라이언트 연결 종료:", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- string(msg)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("메시지 전송 실패:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func main() {
	err := godotenv.Load(filepath.Join("server", ".env"))
	if err != nil {
		log.Fatal("환경 변수 파일(.env) 로드 실패:", err)
	}

	serverIP := os.Getenv("SERVER_IP")
	addrPort := os.Getenv("ADDR_PORT")
	serverAddr := fmt.Sprintf("%s:%s", serverIP, addrPort)

	certPath := filepath.Join("cert", "server.crt")
	keyPath := filepath.Join("cert", "server.key")

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	server := &http.Server{
		Addr:      serverAddr,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	fmt.Println("WebSocket Secure 서버 시작:", "wss://"+serverAddr+"/ws")
	err = server.ListenAndServeTLS(certPath, keyPath)
	if err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}


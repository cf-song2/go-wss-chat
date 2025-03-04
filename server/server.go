package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
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

var clients sync.Map
var broadcast = make(chan string)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 업그레이드 실패:", err)
		return
	}
	defer conn.Close()

	clients.Store(conn, true)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("클라이언트 연결 종료:", err)
			clients.Delete(conn)
			break
		}
		broadcast <- string(msg)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		clients.Range(func(key, value interface{}) bool {
			client := key.(*websocket.Conn)
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("메시지 전송 실패:", err)
				client.Close()
				clients.Delete(client)
			}
			return true
		})
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

	tcpListener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("IPv4/IPv6 리슨 실패: %v", err)
	}

	certPath := filepath.Join("cert", "fullchain.pem")
	keyPath := filepath.Join("cert", "server.key")

	// ✅ 정적 파일 서빙 추가
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	server := &http.Server{
		Addr:      serverAddr,
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	fmt.Println("WebSocket Secure 서버 시작:", "wss://"+serverAddr+"/ws")
	err = server.ServeTLS(tcpListener, certPath, keyPath)
	if err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}


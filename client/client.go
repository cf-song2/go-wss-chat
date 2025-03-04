package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "wss://spectrum.cecil-personal.site/ws"

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	conn, _, err := dialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("서버 연결 실패:", err)
	}
	defer conn.Close()

	fmt.Println("서버에 연결되었습니다. 메시지를 입력하세요:")

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("서버에서 메시지 수신 실패:", err)
				return
			}
			fmt.Println("서버:", string(msg))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		if err != nil {
			log.Println("메시지 전송 실패:", err)
			return
		}
	}
}


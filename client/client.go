package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "wss://spectrum.cecil-personal.site/ws"

	// 재연결 로직 추가
	for attempt := 1; attempt <= 5; attempt++ {
		fmt.Printf("서버에 연결 시도 중... (Attempt %d/5)\n", attempt)

		dialer := websocket.Dialer{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}

		conn, _, err := dialer.Dial(serverAddr, nil)
		if err != nil {
			log.Println("서버 연결 실패:", err)
			time.Sleep(3 * time.Second) // 재시도 전 3초 대기
			continue
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
		break
	}
}


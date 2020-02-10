package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:5000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/helloworld"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	done := make(chan struct{})
	// open后就发送
	err = ws.WriteMessage(websocket.TextMessage, []byte("helloworld"))
	if err != nil {
		log.Println("write:", err)
	}

	go func() {
		defer close(done)
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			} else {
				switch mt {
				case websocket.CloseMessage:
					{
						log.Println("disconnected")
						return
					}
				case websocket.PingMessage, websocket.PongMessage:
					{
						log.Println("get ping pong")
					}
				case websocket.TextMessage:
					{
						msg := string(message)
						log.Printf("recv: %s", msg)
						return
					}
				case websocket.BinaryMessage:
					{
						log.Println("not support Binary now")
						return
					}
				default:
					{
						log.Println("not support now")
						return
					}
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

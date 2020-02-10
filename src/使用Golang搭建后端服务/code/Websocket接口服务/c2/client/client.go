package main

import (
	listener "c2/clienthanddler"
	config "c2/config"
	event "c2/event"
	logger "c2/logger"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func client(address string) {
	logger.Info("connecting to ", address)
	ws, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		logger.Fatal("dial:", err)
	}
	defer ws.Close()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})
	// open后就发送
	message := event.Event{EventType: "message", Message: "I'm HSZ"}
	err = ws.WriteJSON(message)
	if err != nil {
		logger.Error("write:", err)
	}
	go listener.Listerner(ws, done)

	message = event.Event{EventType: "publish", Message: "hello all, I'm HSZ"}
	err = ws.WriteJSON(message)

	message = event.Event{EventType: "publish_no_self", Message: "hello all except me, I'm HSZ"}
	err = ws.WriteJSON(message)
	for {
		select {
		case <-done:
			return

		case <-interrupt:
			logger.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error("write close:", err)
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

func main() {
	conf := config.LoadConfig()
	if conf.Debug == false {
		logger.Logger.SetLevel(log.InfoLevel)
	} else {
		logger.Logger.SetLevel(log.DebugLevel)
	}
	u1 := url.URL{Scheme: "ws", Host: conf.Address, Path: "/room/no1"}
	go client(u1.String())
	u2 := url.URL{Scheme: "ws", Host: conf.Address, Path: "/room/no2"}
	client(u2.String())
}

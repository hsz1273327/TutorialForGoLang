package serverhanddler

import (
	errors "c2/errors"
	event "c2/event"
	logger "c2/logger"

	"github.com/gorilla/websocket"
)

func Listerner(ws *websocket.Conn, done chan struct{}) error {
	defer close(done)
	for {
		e := event.Event{}
		err := ws.ReadJSON(&e)
		if err != nil {
			logger.Error("read:", err)
			return errors.GetMessageError
		}
		switch e.EventType {
		case "close":
			{
				logger.Info("disconnected")
				return nil
			}
		case "message":
			{
				logger.Info("get message:", e.Message)
			}
		default:
			{
				message := event.Event{EventType: "message", Message: "unkonwn command:" + e.EventType}
				err = ws.WriteJSON(message)
				if err != nil {
					logger.Error("write:", err)
					return errors.WriteMessageError
				}
			}
		}
	}
}

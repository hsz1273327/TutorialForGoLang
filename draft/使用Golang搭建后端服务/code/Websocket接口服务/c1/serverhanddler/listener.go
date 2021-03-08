package serverhanddler

import (
	errors "c1/errors"
	event "c1/event"
	logger "c1/logger"

	"github.com/gorilla/websocket"
)

func ServerHanddler(ws *websocket.Conn) error {
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
				return nil
			}
		case "helloworld":
			{
				logger.Debug("get helloworld")
				msg := "welcome " + e.Message
				message := event.Event{EventType: "message", Message: msg}
				err = ws.WriteJSON(message)
				if err != nil {
					logger.Error("write:", err)
					return errors.WriteMessageError
				}
				logger.Debug("send message ", msg)
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

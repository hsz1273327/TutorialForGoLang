package serverhanddler

import (
	errors "c3/errors"
	event "c3/event"
	logger "c3/logger"
	room "c3/room"
	exchange "c3/wsexchange"

	"github.com/gorilla/websocket"
)

func ServerHanddler(ws *websocket.Conn, room_name string, exchange *exchange.Exchange) error {
	for {
		e := event.Event{}
		err := ws.ReadJSON(&e)
		if err != nil {
			logger.Error("read:", err)
			return errors.GetMessageError
		}
		switch e.EventType {
		case "leave":
			{
				exchange.PubNoSelf("some one leave the room", ws)
				return nil
			}
		case "close":
			{
				exchange.Pub("this room will close")
				room.Close(room_name)
				return nil
			}
		case "publish":
			{
				logger.Debug("get publish")
				exchange.Pub(e.Message)
			}
		case "publish_no_self":
			{
				logger.Debug("get publish_no_self")
				exchange.PubNoSelf(e.Message, ws)
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

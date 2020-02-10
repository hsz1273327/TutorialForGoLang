package wsexchange

import (
	event "c2/event"
	logger "c2/logger"

	"github.com/gorilla/websocket"
	"github.com/rfyiamcool/syncmap"
)

type ClientLike interface {
	WriteJSON(interface{}) error
	Close() error
}

type Exchange struct {
	clients syncmap.Map
}

func New() *Exchange {
	var sm syncmap.Map
	exchange := &Exchange{
		clients: sm}
	return exchange
}

func (exchange *Exchange) Len() int64 {
	m := exchange.clients.Length()
	if m != nil {
		logger.Debug("length ", *m)
		return *m
	} else {
		logger.Error("length is nil")
		return 0
	}

}

func (exchange *Exchange) Sub(ws *websocket.Conn) {
	exchange.clients.Store(ws, true)
}

func (exchange *Exchange) DisSub(ws *websocket.Conn) {
	ok := exchange.clients.Delete(ws)
	if ok {
		logger.Debug("ws dissub the exchange")
	} else {
		logger.Debug("ws not in exchange")
	}

}

func (exchange *Exchange) Pub(msg string) {
	message := event.Event{EventType: "message", Message: msg}
	exchange.clients.Range(func(key, value interface{}) bool {
		client := key.(*websocket.Conn)
		err := client.WriteJSON(message)
		if err != nil {
			logger.Error("send to %v error: %v", client, err)
			return false
		}
		return true
	})
}

func (exchange *Exchange) PubNoSelf(msg string, ws *websocket.Conn) {
	message := event.Event{EventType: "message", Message: msg}
	exchange.clients.Range(func(key, value interface{}) bool {
		client := key.(*websocket.Conn)
		if client != ws {
			err := client.WriteJSON(message)
			if err != nil {
				logger.Error("send to %v error: %v", client, err)
				return false
			} else {
				return true
			}
		} else {
			return true
		}
	})
}
func (exchange *Exchange) Close() {
	exchange.clients.Range(func(key, value interface{}) bool {
		client := key.(*websocket.Conn)
		client.Close()
		return true
	})
}

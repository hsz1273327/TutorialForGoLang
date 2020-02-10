package wsexchange

import (
	errors "c3/errors"
	event "c3/event"
	logger "c3/logger"

	set "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
)

type ClientLike interface {
	WriteJSON(interface{}) error
	Close() error
}

type Exchange struct {
	clients set.Set
}

func New() *Exchange {
	exchange := &Exchange{
		clients: set.NewSet(),
	}
	return exchange
}

func (exchange *Exchange) Len() int {
	return exchange.clients.Cardinality()
}

func (exchange *Exchange) Sub(ws *websocket.Conn) error {
	ok := exchange.clients.Add(ws)
	if ok {
		return nil
	} else {
		return errors.WsExistError
	}
}

func (exchange *Exchange) DisSub(ws *websocket.Conn) {
	ok := exchange.clients.Contains(ws)
	if ok {
		exchange.clients.Remove(ws)
		logger.Debug("ws dissub the exchange")
	} else {
		logger.Debug("ws not in exchange")
	}
}

func (exchange *Exchange) Pub(msg string) {
	message := event.Event{EventType: "message", Message: msg}
	exchange.clients.Each(func(element interface{}) bool {
		client := element.(*websocket.Conn)
		err := client.WriteJSON(message)
		if err != nil {
			logger.Error("send to %v error: %v", client, err)
			return false
		} else {
			return true
		}
	})
}

func (exchange *Exchange) PubNoSelf(msg string, ws *websocket.Conn) {
	message := event.Event{EventType: "message", Message: msg}
	exchange.clients.Each(func(element interface{}) bool {
		client := element.(*websocket.Conn)
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
	exchange.clients.Each(func(element interface{}) bool {
		client := element.(*websocket.Conn)
		client.Close()
		return true
	})
}

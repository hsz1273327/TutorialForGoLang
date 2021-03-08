package room

import (
	errors "c2/errors"
	logger "c2/logger"
	exchange "c2/wsexchange"
	"time"

	"github.com/rfyiamcool/syncmap"
)

type RoomManager struct {
	rooms syncmap.Map
}

func New() *RoomManager {
	var sm syncmap.Map
	room_manager := &RoomManager{rooms: sm}
	return room_manager
}
func (rm *RoomManager) Len() int64 {
	return *rm.rooms.Length()
}

func (rm *RoomManager) Add(room_name string, exch *exchange.Exchange) {
	rm.rooms.Store(room_name, exch)
}
func (rm *RoomManager) GetRoom(room_name string) (*exchange.Exchange, error) {
	value, ok := rm.rooms.Load(room_name)
	if ok {
		room_exchange := value.(*exchange.Exchange)
		return room_exchange, nil
	} else {
		return nil, errors.RoomNotExistError
	}
}

func (rm *RoomManager) Close(room_name string) error {
	value, ok := rm.rooms.Load(room_name)
	if ok {
		room_exchange := value.(*exchange.Exchange)
		room_exchange.Close()
		rm.rooms.Delete(room_name)
		logger.Info("room close", room_name)
		return nil
	} else {
		logger.Info("room not exist", room_name)
		return errors.RoomNotExistError
	}
}

func (rm *RoomManager) AutoClose() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for _ = range ticker.C {
			rm.rooms.Range(func(key, value interface{}) bool {
				name := key.(string)
				exch := value.(*exchange.Exchange)
				if exch.Len() == 0 {
					rm.Close(name)
					return true
				} else {
					return false
				}
			})
		}
	}()
}

var DefaultRoomManager *RoomManager = New()

func Len() int64 {
	return DefaultRoomManager.Len()
}

func Add(room_name string, exch *exchange.Exchange) {
	DefaultRoomManager.Add(room_name, exch)
}

func GetRoom(room_name string) (*exchange.Exchange, error) {
	return DefaultRoomManager.GetRoom(room_name)
}

func Close(room_name string) error {
	return DefaultRoomManager.Close(room_name)
}

func AutoClose() {
	DefaultRoomManager.AutoClose()
}

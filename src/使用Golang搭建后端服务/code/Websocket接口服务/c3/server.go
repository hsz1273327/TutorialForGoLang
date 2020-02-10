package main

import (
	config "c3/config"
	logger "c3/logger"
	room "c3/room"
	handdler "c3/serverhanddler"
	exchange "c3/wsexchange"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

func time_pusher() {
	room.ForEach(func(key, value interface{}) bool {
		room_exchange := value.(*exchange.Exchange)
		if room_exchange.Len() != 0 {
			room_exchange.Pub(time.Now().String())
			return true
		} else {
			return false
		}
	})
}

func AutoPush() {
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for _ = range ticker.C {
			time_pusher()
		}
	}()
}

func RoomHttpHanddler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer ws.Close()
	vars := mux.Vars(r)
	room_name := vars["room_name"]
	var room_exchange *exchange.Exchange
	re, err := room.GetRoom(room_name)
	if err != nil {
		room_exchange = exchange.New()
		room.Add(room_name, room_exchange)
	} else {
		room_exchange = re
	}
	room_exchange.Sub(ws)
	defer room_exchange.DisSub(ws)
	logger.Info("conncet to ", room_name)
	handdler.ServerHanddler(ws, room_name, room_exchange)
}

func server(addr string) {
	r := mux.NewRouter()
	r.HandleFunc("/room/{room_name}", RoomHttpHanddler)
	http.Handle("/", r)
	logger.Info(http.ListenAndServe(addr, nil))
}

func main() {
	conf := config.LoadConfig()
	if conf.Debug == false {
		logger.Logger.SetLevel(log.WarnLevel)
	} else {
		logger.Logger.SetLevel(log.DebugLevel)
	}
	logger.Logger.Info("start @ ", conf.Address)
	room.AutoClose()
	AutoPush()
	server(conf.Address)

}

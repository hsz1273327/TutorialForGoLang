package main

import (
	config "c1/config"
	logger "c1/logger"
	wshanddlers "c1/serverhanddler"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

func helloworldHttpHanddler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer ws.Close()
	wshanddlers.ServerHanddler(ws)
}

func server(addr string) {
	http.HandleFunc("/helloworld", helloworldHttpHanddler)
	logger.Info(http.ListenAndServe(addr, nil))
}

func main() {
	conf := config.LoadConfig()
	if conf.Debug == false {
		logger.Logger.SetLevel(log.WarnLevel)
	} else {
		logger.Logger.SetLevel(log.DebugLevel)
	}
	logger.Info("start @", conf.Address)
	server(conf.Address)
}

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:5000", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func helloworldWsHanddler(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		} else {
			switch mt {
			case websocket.CloseMessage:
				{
					log.Println("get close signal")
					break
				}
			case websocket.PingMessage, websocket.PongMessage:
				{
					log.Println("get ping pong")
				}
			case websocket.TextMessage:
				{
					log.Printf("recv: %s", message)
					msg := string(message)
					switch msg {
					case "close":
						{
							break
						}
					case "helloworld":
						{
							err = ws.WriteMessage(websocket.TextMessage, []byte("Hello World"))
							if err != nil {
								log.Println("write:", err)
								break
							}
						}
					default:
						{
							err = ws.WriteMessage(websocket.TextMessage, []byte("unkonwn command"))
							if err != nil {
								log.Println("write:", err)
								break
							}
						}
					}
				}
			case websocket.BinaryMessage:
				{
					log.Println("not support Binary now")
				}
			default:
				{
					log.Println("not support now")
				}
			}
		}
	}
}

func helloworldHttpHanddler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()
	helloworldWsHanddler(ws)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/helloworld", helloworldHttpHanddler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

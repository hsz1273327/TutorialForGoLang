# 使用Golang构建Websocket接口服务

Golang可以使用[gorilla/websocket](https://github.com/gorilla/websocket)这个框架来实现websocket接口的构造.这个框架可以用于写客户端和服务器.

我们依然从一个[helloworld](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c0)开始.这个例子我们在客户端连同服务端后立即发送一个helloworld消息给后端服务器,服务器接到后则返回一个helloworld消息给客户端.
客户端在接收到服务器消息后发送一个close消息给服务器,服务器就断开和客户端的连接.

+ 服务端

```golang
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
```

服务端我们需要用一个http服务器监听一个url.在有客户端访问后,使用`websocket.Upgrader{}`来将http访问提升为websocket连接.使用
`mt, message, err := ws.ReadMessage()`来获取消息.返回值的第一项为消息类型,第二项为message体.消息类型包含如下几种:

+ `websocket.BinaryMessage`字节流数据
+ `websocket.TextMessage`文本数据
+ `websocket.PingMessage`保持连接用的消息
+ `websocket.PongMessage`保持连接用的消息
+ `websocket.CloseMessage`关闭信号

而发送回去一般使用的`ws.WriteMessage(int , []byte) error`来发送消息.

因为mt有多种情况,我们一般使用switch来区分.一般为了可读性我们使用`websocket.TextMessage`结合json来传递信息;
当然go语言天生亲和protobuf,所以也可以使用`websocket.BinaryMessage`结合protobuf来传递消息.


+ 客户端

```golang
package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:5000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/helloworld"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	done := make(chan struct{})
	// open后就发送
	err = ws.WriteMessage(websocket.TextMessage, []byte("helloworld"))
	if err != nil {
		log.Println("write:", err)
	}

	go func() {
		defer close(done)
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			} else {
				switch mt {
				case websocket.CloseMessage:
					{
						log.Println("disconnected")
						return
					}
				case websocket.PingMessage, websocket.PongMessage:
					{
						log.Println("get ping pong")
					}
				case websocket.TextMessage:
					{
						msg := string(message)
						log.Printf("recv: %s", msg)
						return
					}
				case websocket.BinaryMessage:
					{
						log.Println("not support Binary now")
						return
					}
				default:
					{
						log.Println("not support now")
						return
					}
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
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
```
客户端部分与上面也类似.区分好消息类型就行

比较坑的是`ws.WriteMessage`输入和`ws.ReadMessage()`都是一个`[]byte`类型的数据.因此即便是`websocket.TextMessage`类型的数据也要做好类型转换.

## Json数据传输

一种更常见的形式是使用结构化数据格式表示一个事件,客户端服务器按事件的不同来进行不同的处理.例子[c1](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c1)就是这个helloworld的改进.顺便我们也把这个代码拆分一下让它更加便于管理.

这个项目被拆分成了如下几个子模块

+ `logger`使用`github.com/sirupsen/logrus`的log模块
+ `config`使用`github.com/paked/configure`加载命令行和config文件的子模块

上面的几个是子模块具有一定的通用性,可以用于设定各种服务的配置.

+ `errors`用于维护这个项目的所有错误
+ `event`事件类型,事件使用json格式
+ `clienthanddler`客户端的事件控制逻辑
+ `serverhanddler`服务端的事件控制逻辑

在客户端和服务端中都使用定义好的事件结构

```golang
type Event struct {
	EventType string `json:"name"`
	Message   string `json:"message"`
}
```

来将消息转成事件处理.读取使用

```golang
e := event.Event{}
ws.ReadJSON(&e)
```

写使用

```golang
message := event.Event{EventType: "message", Message: "unkonwn command"}
err = ws.WriteJSON(message)
```

## 广播消息

广播就是向符合条件的连接同时发送相同的消息,也就是发布订阅模式.

下面的例子[C2](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c2)我们就利用http的url构造room的概念来限定广播的范围,进入一个room就相当于订阅了一个topic,离开room也就相当于取消订阅.为此我们构造了两个结构:

+ Exchange 用于管理一个room的订阅,取消订阅,广播消息,其使用`github.com/rfyiamcool/syncmap`的`syncmap.Map`结构代替set结构维护客户端的连接.这是标准库`sync.Map`的一个扩展,是一个线程安全的map结构,我们用它代替set,之所以不用默认的map是因为它不是线程安全的,我们有要定时删除不用room的逻辑,这可能引起资源冲突.不过要注意这个库有些小bug,用的时候要小心.当然我们也可以使用[github.com/deckarep/golang-set](https://github.com/deckarep/golang-set)中的线程安全set,我们在下一个例子中使用.

	```golang
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
	```


+ RoomManager 用于管理所有的room与Exchange的对应关系,使用`map[string]*exchange.Exchange`结构维护对应关系.

	```golang
	type RoomManager struct {
		rooms syncmap.Map
	}

	func New() *RoomManager {
		var sm syncmap.Map
		room_manager := &RoomManager{
			rooms: sm}
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
	```

	我们定义一个自动关闭的方法,每隔10s检查一次是否有room已经是空的了,如果有那么就将其关闭

这个例子中客户端通过发送类型为`message`,`publish`和`publish_no_self`的3个消息来验证广播功能.我们可以多开几个客户端来检查广播是否可用.


## 主动推送广播

在上面的例子中广播是由客户端触发发起的,这依然是请求响应模式,我们来构造一个由服务端主动推送的例子[C3](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c3),它会每隔3s向房间中的客户端推送当前时间

```golang
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
....

func main() {
	...
	room.AutoClose()
	AutoPush()
	server(conf.Address)
}
```

## 连接管理

我们可以看到这个websocket实现的api是比较底层的,和js原版的差不太多,只有点对点的消息传输,没有连接管理,为了可以做好连接管理我们通常需要一套用户系统,这将是下一节内容的一部分.

# Grpc接口服务

[GRpc](https://grpc.io/)正如其名,是一种RPC.它实际上和RESTful接口在功能上是相近的,本质都是一种请求响应模式的服务.只是作为一个RPC,GRpc一般描述动作而非资源,并且它可以返回的不光是一个数据,而是一组流数据.

GRpc是一种跨语言的Rpc,它建立在http2上使用[protobuf](https://developers.google.com/protocol-buffers/)作为结构化数据的序列化工具,

它有4种形式:

+ 请求-响应
+ 请求-流响应
+ 流请求-响应
+ 流请求-流响应
其基本使用方式是:

+ 服务端与客户端开发者协商创建一个protobuf文件用于定义rpc的形式和方法名以及不同方法传输数据的schema
+ 服务端实现protobuf文件中定义的方法
+ 客户端调用protobuf文件中定义的方法


在go中我们需要使用包[google.golang.org/grpc](https://github.com/grpc/grpc-go),[google.golang.org/genproto](https://github.com/google/go-genproto.git)和[github.com/golang/protobuf/protoc-gen-go](https://github.com/golang/protobuf)来实现上面的三个步骤.需要注意,这两个包都是被墙着的,而且实测用goproxy都没救,只能clone下来后改路劲install,注意安装顺序,先安装`github.com/golang/protobuf/protoc-gen-go`,`google.golang.org/genproto`再装`google.golang.org/grpc`,另外就是`google.golang.org/grpc`和`google.golang.org/genproto`clone下来后记得改名字

## 请求-响应

这个例子[C0](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/GRpc%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c0)我们来实现一个简单的服务--输入一个数,输出这个数的平方

### 最终整个项目的结构

go语言毕竟静态语言,不似动态语言这么灵活,尤其它的包引入机制并不先进,也就比c语言好那么一点.因此必须使用瀑布式的开发流程也就是必须先设计后实现,无法边写边设计.


我们使用go 1.11 新加入的特性`module`来构建项目`c0`,这个项目中会有一个子模块`c0/squarerpc_service`用于存放由protobuf编译来的go模块文件.因为有子模块在其中,我们就需要关注下项目的结构了,否则会无法编译.

```shell
|--bin--|
|          |--server //编译好的服务软件
|          |--client //编译好的客户端软件
|
|--client--|
|          |--client.go //客户端源码
|
|--schema--|
|          |--square_service.proto // 定义的protobuf
|
|--squarerpc_service--|
|                  |--square_service.pb.go //由square_service.proto编译过来的go模块
|                  |--go.mod //定义子模块squarerpc_service的依赖关系
|
|--go.mod //定义项目的依赖关系
|--makefile //定义编译过程
|--server.go //服务端源码
```

需要注意的是为了可以引用本地的子模块,父模块的`go.mod`需要使用`replace`声明将子模块替换为本地.

+ 父级go.mod

```
module c0

require (
	c0/squarerpc_service v0.0.0
	github.com/golang/protobuf v1.3.1
	google.golang.org/grpc v1.19.0
)

replace c0/squarerpc_service v0.0.0 => ./squarerpc_service

go 1.12
```

### 创建一个protobuf文件

go语言已经强制使用大写字母作为public的标志,那我们定义protobuf的时候也最好按这个来定义,否则在go中也会强制转成首字母大写,造成不一致.

```proto
syntax = "proto3";
package squarerpc_service;

service SquareService {
    rpc Square (Message) returns (Message){}
}

message Message {
    double Message = 1;
}
```

在这个项目中我们使用命令

```shell
protoc -I schema schema/square_service.proto --go_out=plugins=grpc:squarerpc_service
```

+ `-I`指定了protobuf文件所在的文件夹
+ `--go_out=plugins=grpc`指定了是grpc使用,`:squarerpc_service`则是指定了编译后的目标文件夹


### 服务端实现定义的方法

go的服务端写起来还是比较简单的

```go
package main

import (
	"context"
	"log"
	"net"

	pb "c0/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

// server 结构体用于实现pb中的SquareService.
type server struct{}

// 为结构体绑定要实现的方法
func (s *server) Square(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Received: %v", in.Message)
	return &pb.Message{Message: in.Message * in.Message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("server started @", ADDRESS)
	s := grpc.NewServer()
	pb.RegisterSquareServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

+ 首先是定义一个空的结构`server`,并为其绑定上我们要实现的pb文件中定义的方法实现.注意返回的是结构的实例地址而非实例

+ 然后使用`net`模块的`Listen`方法创建一个网络监听器`lis`

+ 之后使用`grpc.NewServer()`创建一个服务器`s`

+ 使用pb中的`pb.RegisterSquareServiceServer(s, &server{})`方法把`s`和一个server结构体的实例地址绑定

+ 最后服务器`s`使用`s.Serve(lis)`方法绑定监听器

### 客户端实现方式

客户端一块主要是要关注下连接的退出和过期后的关闭上下文的操作.这边都是用defer关键字实现的.

```go
package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "c0/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	//server address
	ADDRESS = "localhost:5000"
)

func main() {
	// 连接到服务器.
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 创建client对象
	c := pb.NewSquareServiceClient(conn)

	// 构造要发送的请求
	query := 2.0
	if len(os.Args) > 1 {
		query, err = strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			log.Fatalf("query can not parse from string: %v", err)
			panic(err)
		}
	}
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用rpc的函数
	r, err := c.Square(ctx, &pb.Message{Message: query})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Printf("Result: %f", r.Message)
}
```

### 使用makefile流程化编译

都写好了后来定义一个makefile用于流程化编译

```makefile
ASSETS=assets
default: server client

squarerpc_service/square_service.pb.go: schema/square_service.proto squarerpc_service/go.mod
	protoc -I schema schema/square_service.proto --go_out=plugins=grpc:squarerpc_service
server: squarerpc_service/square_service.pb.go server.go
	go build -o $(ASSETS)/server

client: squarerpc_service/square_service.pb.go client/client.go
	go build -o $(ASSETS)/client client/client.go
```
这个makefile肯定不是生产环境可以用的,它都没有做针对部署目标的交叉编译.这边写出来只是给个样例.


## 请求-流响应

这种需求比较常见,有点类似,python中的range函数,它生成的是一个流而非一个数组,它会一次一条的按顺序将数据发送回请求的客户端.

这个例子[C1](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/GRpc%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c1)实现了给出一个正整数,它会返回从0开始到它为止的每个整数的平方.

### 修改protobuf文件

```proto
...
service SquareService {
    rpc RangeSquare (Message) returns (stream Message){}
}
...
```

我们只要在返回值前面声明stream就可以定义一个流响应

### 修改服务端

```go
// 第一位为请求,第二位为流对象,用于发送回消息
func (s *server) RangeSquare(in *pb.Message, stream pb.SquareService_RangeSquareServer) error {
	limit := int(in.Message)
	for i := 0; i <= limit; i++ {
		err := stream.Send(&pb.Message{Message: float64(i * i)})
		if err != nil {
			return err
		}
	}
	return nil
}
```

go的流响应主要是靠调用参数中的stream的`Send`方法实现的.

### 修改客户端

```go
stream, err := client.RangeSquare(ctx, &pb.Message{Message: query})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	for {
		result, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("%v.RangeSquare(_) = _, %v", client, err)
			}
		}
		log.Printf("Result: %f", result.Message)
	}
```

客户端调用在前面没什么不同,只是返回的值是一个stream,我们需要使用for语句不断调用`stream.Recv()`直到返回的err为`io.EOF`为止.


## 流请求-响应

这种需求不是很多见,可能用的比较多的是收集一串数据后统一进行处理吧,流只是可以确保是同一个客户端发过来的而已.

这个例子[C2](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/GRpc%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c2)实现了传过来一串数,之后返回他们的平方和

### 修改protobuf文件

```protobuf
...
service SquareService {
    rpc SumSquare (stream Message) returns (Message){}
}
...
```

我们只要在请求前面声明stream就可以定义一个流响应

### 修改服务端

```go
func (s *server) SumSquare(stream pb.SquareService_SumSquareServer) error {
	var sum float64 = 0.0
	for {
		data, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&pb.Message{Message: sum})
			}
			if err != nil {
				return err
			}
		} else {
			sum += data.Message * data.Message
		}
	}
}
```

请求流的返回就是一个error,操作其实就在stream上,接收使用`Recv()`,返回使用`SendAndClose()`

### 修改客户端

```go
stream, err := c.SumSquare(ctx)
if err != nil {
    log.Fatalf("could not call: %v", err)
}
for i := 0; i < int(query); i++ {
    msg := pb.Message{Message: float64(i)}
    err := stream.Send(&msg)
    if err != nil {
        if err == io.EOF {
            break
        }
        log.Fatalf("%v.Send(%v) = %v", stream, msg, err)
    }
}
reply, err := stream.CloseAndRecv()
if err != nil {
    log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
}
log.Printf("Route summary: %v", reply)
```

请求也是类似,由于go的语法表现力不足,所以只能调用后返回一个流对象,再使用流的`CloseAndRecv()`手工结束流后等待结果返回.

## 流请求-流响应

将上面两种方式结合起来,就是我们的第四种方式,请求为一个流,响应也是流.代码在[C3](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/GRpc%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c3)可以这两个流可以是相互交叉的也可以是请求完后再返回一个流.他们在写pb文件时是相同的写法

```protobuf
service SquareService {
    rpc StreamrangeSquare (stream Message) returns (stream Message){}
}
```

### 请求流完成后返回流

+ 修改服务端

```go
func (s *server) StreamrangeSquare(stream pb.SquareService_StreamrangeSquareServer) error {
	l := list.New()
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			l.PushBack(in.Message * in.Message)
		}
	}
	for e := l.Front(); e != nil; e = e.Next() {
		stream.Send(&pb.Message{Message: e.Value.(float64)})
	}
	return nil
}
```


+ 修改客户端

```go
// 调用rpc的函数
stream, err := c.StreamrangeSquare(ctx)
if err != nil {
    log.Fatalf("could not call: %v", err)
}
waitc := make(chan struct{})
go func() {
    for {
        in, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                // read done.
                close(waitc)
                return
            } else {
                log.Fatalf("Failed to receive a note : %v", err)
            }
        } else {
            log.Printf("Got message %f", in.Message)
        }

    }
}()
limit := int(query)
for i := 0; i < limit; i++ {
    stream.Send(&pb.Message{Message: float64(i)})
}
stream.CloseSend()
<-waitc
```

客户端有点特殊,使用子协程处理读入,主协程处理请求流,通过阻塞信道waitc来等待读取请求完成,当读取完成后协程会关闭waitc,这样阻塞取消程序就正常结束了.

### 请求的进行中就返回响应

代码在[C4](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/GRpc%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c4)

+ 修改服务端

```go
func (s *server) StreamrangeSquare(stream pb.SquareService_StreamrangeSquareServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			err := stream.Send(&pb.Message{Message: in.Message * in.Message})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
```
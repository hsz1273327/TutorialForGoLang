# 使用Golang搭建后端服务

go语言现在的主战场就是网络编程.这篇其实主要也是为了讲最常用的go语言网络编程工具.

网络编程最常见的就两种形式:

1. http协议的服务,一般用于和前端业务交互以完成业务需求
2. rpc服务,一般用于封装算法供业务调取

冷门些的可能会有一些p2p协议的网络编程,这种更多的用在加密通信和一些特殊领域

本文主要讲:

+ http服务
    + RESTful接口服务
    + SSE服务用于推送消息

+ grpc服务

+ websocekt服务

+ p2p程序
    + webrtc程序(+ [webrtc](https://github.com/pions/webrtc)一种在浏览器端也有实现的p2p即时通信技术)

而后端相关的技术有:

+ 关系数据库技术,常用于保存业务数据.常见的有
    + [PostgreSQL](http://www.postgres.cn/docs/12/),一般用在服务端
    + [sqlite3](https://www.sqlite.org/doclist.html),一般用在客户端

+ orm技术,业务上一般使用orm来操作关系数据库.常用的orm有[xormplus](https://github.com/xormplus/xorm)

+ 共享内存技术,常见的是Redis,我们一般使用[redis](https://github.com/go-redis/redis)

+ 消息中间件技术,常见的有:
    + rabbitMQ,用于相对轻量的分发任务.我们使用[amqp](https://github.com/streadway/amqp)
    + redis,用于在追求实时性,不在意数据完整性的业务场景下使用,常见的场景比如聊天室.我们通常用[github.com/go-redis/redis/v8](https://github.com/go-redis/redis)
    + kafka,用于在严格追求数据完整性和顺序,同时对吞吐量有要求时使用,比如业务层向数据层同步数据,事件驱动任务等.我们通常使用[gopkg.in/confluentinc/confluent-kafka-go.v1/kafka](https://github.com/confluentinc/confluent-kafka-go)

+ 数据序列化反序列化技术,常见的有:
    + 标准库的json性能比较差,我们有时会用[github.com/json-iterator/go](https://github.com/json-iterator/go)替代
    + 有时我们也会考虑使用[msgpack](https://github.com/vmihailenco/msgpack)代替json作为表现层序列化协议
    + 另一种思路使用必须预先定义好schema的[Protobuffer](https://github.com/protocolbuffers/protobuf),一般用在rpc技术或者一些定义严格的服务中

+ 科学计算,虽然多数时候go主要处理的都是io密集型任务,但指不定也会用到要计算的部分,一般会用[gonum](https://github.com/gonum/gonum)

这些库就不一一介绍了,用的的时候去现查即可.
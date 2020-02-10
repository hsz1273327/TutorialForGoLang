# 使用Golang搭建后端服务

go语言现在的主战场就是网络编程.这篇其实主要也是为了讲最常用的go语言网络编程工具.

本文主要讲:

+ 使用gin构造RESTful接口服务
+ websocekt构造接口服务
+ grpc构造接口服务

而后端相关的技术有:

+ 关系数据库技术,常见的有[PostgreSQL](http://www.postgres.cn/docs/10/),业务上一般使用orm来操作数据库.常用的orm有[xorm](https://github.com/go-xorm/xorm)
+ 共享内存技术,常见的是Redis,我们使用[redis](https://github.com/go-redis/redis)
+ 消息队列技术,常见的有rabbitMQ,我们使用[amqp](https://github.com/streadway/amqp);redis;kafka我们用[gopkg.in/confluentinc/confluent-kafka-go.v1/kafka](https://github.com/confluentinc/confluent-kafka-go)
+ 消息的发布订阅工具,常见的有Redis,rabbitMQ,postgreSQL

+ [zmq](https://github.com/pebbe/zmq4)一种基于通信模式的消息组件框架

+ [webrtc](https://github.com/pions/webrtc)一种在浏览器端也有实现的p2p即时通信技术

+ 由于go语言标准库的log工具比较弱,我们有时用[logrus](https://github.com/sirupsen/logrus)来代替
+ 标准库的json性能比较差,我们有时会用[github.com/json-iterator/go](https://github.com/json-iterator/go)替代
+ 有时我们也会考虑使用[msgpack](https://github.com/vmihailenco/msgpack)代替json作为表现层协议
+ 标准库没有原生的协程池,我们有时用[ants](https://github.com/panjf2000/ants)
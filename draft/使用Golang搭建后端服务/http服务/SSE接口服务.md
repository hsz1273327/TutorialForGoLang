# SSE接口服务

sse接口一般分为两种形式.

+ 一种是服务启动后就会自动推送的,比如订阅时钟.
+ 一种是需要创建后才会推送的,我们的例子聊天室就是后一种的应用.

第一种我们只需要定义一个GET方法的接口,用户请求后一直给他推送数据即可

第二种我们一般提供3个接口:

1. 创建接口,用于创建一个流按逻辑推送数据到一个发布订阅器
2. 单独监听接口,用于给用户推送创建的流里的内容
3. [可选]公频监听接口,用于将所有创建出来的流推送给请求了这个接口的用户

无论是哪一种,我们的sse接口都需要一个发布订阅器用于广播消息给订阅的用户.

## 发布订阅器

我们发布订阅器的要求是可以广播消息,在考虑水平扩展的情况下我们一般需要借助外部的有状态服务实现,最常用的就是[消息中间件](https://blog.hszofficial.site/experiment/2019/04/09/%E5%B8%B8%E8%A7%81%E7%9A%84%E6%B6%88%E6%81%AF%E4%B8%AD%E9%97%B4%E4%BB%B6%E6%A8%A1%E5%BC%8F/).

一般是两种形式:

1. 直接消费消息中间件中的结果(对应服务启动后就会自动推送的的数据)
2. 借助消息中间件广播生产的消息(对应需要创建后才会推送的的服务)

无论哪种最关键的部分就是将结果用流的形式发送出去

## 发送SSE

`gin`中使用`c.Stream(func (w io.Writer) bool)`发送流数据,这个`Stream`方法的参数是一个函数,它会循环的一直被调用直到它的返回值为`false`为止.当函数的返回值为false时`c.Stream(func (w io.Writer) bool)`这个方法也就执行结束了.一般来说这个sse推送也就关闭了

而在这个函数中我们需要控制每次的发送数据,`gin.Context`自己提供了一个`c.SSEvent(name string, message interface{})`方法作为sse发送的快捷方式,但很明显虽然够用,但实际上是不完整的.相对完整的方法是使用包`"github.com/gin-contrib/sse"`(虽然还是缺少注释这种字段),用`c.Render(-1, sse.Event)`替代上面的方法发送消息.

总结下一个发送sse的模板大致可以套用如下代码:

```golang
c.Header("Connection", "keep-alive") // 设置流的头以控制keep-alive
c.Writer.Flush()//清空缓冲区数据
clientGone := c.Writer.CloseNotify()// 获取客户端主动断开的信号
c.Stream(func(w io.Writer) bool {
    select {
    case <-clientGone: //监听客户端是否断开
        {
            log.Debug("client close", log.Dict{
                "channelid": "timer::" + cq.ChannelID,
            })
            return false
        }
    case message, isopen := <-ssech: //监听sse数据的生成流,通常借助消息中间件获得
        {
            if isopen {
                log.Debug("channel open", log.Dict{
                    "channelid": "timer::" + cq.ChannelID,
                })
                msg := Event{}
                err := json.UnmarshalFromString(message.Payload, &msg)
                if err != nil {
                    log.Error("UnmarshalFromString error", log.Dict{"Payload": message.Payload})
                    return true
                }
                if msg.Event == "EOF" {
                    log.Debug("publisher close", log.Dict{
                        "channelid": "timer::" + cq.ChannelID,
                    })
                    return false
                }
                // 发送一段sse数据到客户端
                c.Render(-1, &sse.Event{
                    Id:    msg.Id,
                    Event: msg.Event,
                    Data:  msg.Data,
                    Retry: msg.Retry,
                })
                return true
            } else {
                log.Debug("channel close", log.Dict{
                    "channelid": "timer::" + cq.ChannelID,
                })
                return false
            }
        }
    }
})
```

## 测试sse接口

gin并没有提供测试sse的工具,我们只能先起一个服务然后使用http客户端和sse客户端来实现测试.

http客户端可以直接用标准库,sse的客户端我们可以使用[github.com/r3labs/sse/v2](https://github.com/r3labs/sse)这个包既包含客户端实现也包含服务器实现,但由于一般业务上很少有纯sse的服务端的需求,就不推荐用它来实现服务端了.

大致形式如下:

```golang
//Http请求激活channel并获得id
req, _ := http.NewRequest("POST", "http://localhost:5000/v1_0_0/event/timer", bytes.NewBuffer([]byte(`{"second": 10}`)))
req.Header.Set("Content-Type", "application/json")
httpclient := &http.Client{}
resp, err := httpclient.Do(req)
if err != nil {
    panic(err)
}
defer resp.Body.Close()
assert.Equal(t, 200, resp.StatusCode)
body, _ := ioutil.ReadAll(resp.Body)
msg := timernamespace.CounterDownResponse{}
json.Unmarshal(body, &msg)
log.Info("get res", log.Dict{"msg": msg})

//sse请求监听消息
sse_client := sse_test.NewClient("http://localhost:5000/v1_0_0/event/timer/" + msg.ChannelID)
sse_client.Subscribe("messages", func(msg *sse_test.Event) {
    // Got some data!
    fmt.Println(string(msg.Data))
})
```

## 响应流应用扩展

sse不过是http响应流的一个应用,实际上响应流的应用领域非常多,而与构造接口最相关的是提供表格下载接口.一般的表格可能直接写就是了,但非流响应的数据规模上限

## 实现聊天室
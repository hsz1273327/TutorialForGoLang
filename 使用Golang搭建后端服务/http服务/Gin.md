# Gin

[Gin](https://github.com/gin-gonic)是一个相当流行且简单的http服务框架,它的性能足够强大,关键是结构足够简单可以充分扩展以适应各种需求.本文算是这部分的一个预热,先简单介绍下Gin的使用,以便于快速进入正题

## helloworld

一个最基本的gin项目大致是这样一个结构:

```bash
|--dist--|
|        |--server //编译好的服务执行程序
|
|--go.mod //模块信息描述依赖
|--helloworld.go //源码
```

我们来根据这个结构构造一个最简单的helloworld程序

+ `helloworld.go`

    ```golang
    package main

    import (
        "github.com/gin-gonic/gin"
    )


    func main() {
        router := gin.Default()
        router.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "pong",
            })
        })
        router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
    }
    ```

可以看出一个最简单的gin程序,需要如下几个部分:

1. 一个`*gin.Engine`的路由实例,上例中就是`router`
2. 为这个路由实例绑定执行过程,
   1. `GET`方法表示handler处理的是Http方法`GET`的请求,同理还有`POST`,`PUT`,`DELETE`,`PATCH`,`HEAD`,`OPTIONS`;
   2. 其第一个参数为请求的路径,这个路径可以是类似例子中这样静态的路径,也可以是动态路径,动态路径有两种形式:
        1. 如`"/user/:name"`这样的路径参数.它会匹配`/user/john`这样的路径,但不会匹配`/user`或者`/user/`
        2. 如`"/user/:name/*action"`这样的路径参数.`/user/john/`和`/user/john/send`这样的路径,同时如果没有其他路由匹配`/user/john`,它将重定向到`/user/john/`
        动态路由中动态的部分可以在处理过程中通过`c.Param("name")`这样的方式获得,但注意**得到的都是字符串**
   3. 其除第一个参数以外的参数都是对符合上面定义的路径和Http方法的处理过程(`gin.HandlerFunc`),他们都有一个统一的入参`c *gin.Context`并且没有返回.每个处理过程可以分为两种:
        1. 内部有`c.Next()`的称为中间件(middleware)
        2. 内部没有`c.Next()`的我们暂且称为handler
        在处理请求时gin会从左到右的调用注册进去的处理过程,直到处理完handler后就不再向后处理了并再次从右往左的将前面的中间件都执行完.这样完成一次请求-响应的处理链路

3. 在绑定好路由的执行过程后我们需要执行`Run(address)`让程序跑起来等待请求

## 路由组

一个http服务当然不会只有一个可用路由,为了更好的组织和维护路由,gin提供了路由组(`gin.RouterGroup`)这个工具.

```golang
package main

import (
    "github.com/gin-gonic/gin"
)

const (
    ADDRESS = "0.0.0.0:5000"
)

func main() {
    router := gin.Default()
    group1 := router.Group("/group1")
    {
        group1.GET("/ping", func(ctx *gin.Context) {
            ctx.JSON(200, gin.H{
                "message": "pong",
            })
        })
        group1.GET("/pong", func(ctx *gin.Context) {
            ctx.JSON(200, gin.H{
                "message": "ping",
            })
        })
    }
    router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}
```

使用路由组我们就可以将不同含义的路由进行分组,同时用路由组还可以将业务逻辑和具体的某个`gin.Engine`实例解耦,从而达到复用的效果

注意我们可以定义`router.RedirectTrailingSlash = true`自动重新定向末尾为`/`的请求到去掉末尾`/`的路径

如果我们需要在启动时就获取整个app的路由注册表,我们可以通过调用`router.Routes()`或则路由信息列表.这个方法在需要做权限控制的场景非常有用,可以用于注册资源权限到数据库等

## `gin.Context`

我们的处理过程中全程有参数`c *gin.Context`,它几乎覆盖了所有响应流程中需要用到的功能,包括:

1. 传递请求数据和状态
2. 向整个请求-响应链上传递数据
3. 构造响应

### 传递请求数据和状态

下面是可以直接获得的数据和状态表格

| 含义                                                          | 获取方法                                                      | 类型                            |
| ------------------------------------------------------------- | ------------------------------------------------------------- | ------------------------------- |
| 路由参数                                                      | `c.Param("name")`                                             | `string`                        |
| url请求参数                                                   | `c.DefaultQuery("firstname", "Guest")`或`c.Query("lastname")` | `string`                        |
| `Content-Type: application/x-www-form-urlencoded`时的body参数 | `c.PostForm("name")` 或`c.PostFormMap("namemap")`             | `string`或`map[string][string]` |
| 获取header中字段                                              | `c.GetHeader(key string)`                                     | `string`                        |
| header中的cookie                                              | `c.Cookie("gin_cookie")`                                      | string                          |
| 请求header中的`Content-Type`字段                              | `c.ContentType()`                                             | `string`                        |
| 请求方ip信息                                                  | `c.ClientIP()`                                                | `string`                        |
| 请求的完整路径(不含host,port信息)                             | `c.FullPath()`                                                | string                          |

除了上面可以直接获取的数据我们也可以直接访问`c.Request`来获取原始的请求数据,它的类型就是[标准库`net/http`里的`Request`](http://doc.golang.ltd/pkg/net_http.htm#Request)要啥有啥可以直接获取

除了上面的方法,我们还可以通过接口将对应位置的数据绑定到事先定义好的结构体上,当然如果数据和结构体不匹配会报错.

+ 针对header,可以使用`c.ShouldBindHeader`

    ```golang
    type testHeader struct {
        Rate   int    `header:"Rate"`
        Domain string `header:"Domain"`
    }
    h := testHeader{}
    err := c.ShouldBindHeader(&h)
    ```

+ 针对路由参数,可以使用`c.ShouldBindUri`

    ```golang
    type Person struct {
        ID string `uri:"id" binding:"required,uuid"`
        Name string `uri:"name" binding:"required"`
    }

    person := Person{}
    err := c.ShouldBindUri(&person)
    ```

+ 针对url请求参数

    ```golang
    type Person struct {
        Name    string `form:"name"`
        Address string `form:"address"`
    }

    
    person := Person{}
    err := c.ShouldBindQuery(&person)
    ```

+ 针对body,可以使用`c.ShouldBindBodyWith(&objA,format interface{})`

    ```golang
    type Person struct {
        Name    string `json:"name"`
        Address string `json:"address"`
    }

    
    person := Person{}
    err := c.ShouldBindBodyWith(&person,binding.JSON)
    ```

    `format`支持的序列化格式有:

    | 序列化格式      | format值                |
    | --------------- | ----------------------- |
    | `JSON`          | `binding.JSON`          |
    | `MsgPack`       | `binding.MsgPack`       |
    | `YAML`          | `binding.YAML`          |
    | `ProtoBuf`      | `binding.ProtoBuf`      |
    | `XML`           | `binding.XML`           |
    | `Form`          | `binding.Form`          |
    | `FormPost`      | `binding.FormPost`      |
    | `FormMultipart` | `binding.FormMultipart` |

### 向整个请求-响应链上传递数据

`gin.Context`可以使用`Set(key string, value interface{})`和`Get(key string) (interface{},bool)`两个接口在其上设置键值对数据.从而在中间件和handler间传递数据.这样在数据上一条请求-响应链路上的数据就可以打通了.

### 构造响应

响应一般可以分为两部分

1. header

2. body

#### 构造响应头

响应头可以通过`c.Header(key string, value string)`来设置.

比较特殊的可以通过`c.Status(int)`来设置返回的状态码,`SetCookie`可以设置头部的Cookie

#### 构造响应body

`gin.Context`默认支持渲染如下类型的数据以构造响应body

| 格式       | 方法                                                                              |
| ---------- | --------------------------------------------------------------------------------- |
| `html`     | `c.HTML(state_code int, html_template_path string, args map[string]interface{})`  |
| `json`     | `c.JSON(state_code int, interface{})`或`c.AsciiJSON(state_code int, interface{})` |
| `xml`      | `c.XML(state_code int, interface{})`                                              |
| `yaml`     | `c.YAML(state_code int, interface{})`                                             |
| `protobuf` | `c.ProtoBuf(state_code int, interface{})`                                         |

很奇怪的是msgpack的支持是现成的,但gin就是没有实现它的render,我们可以强行使用`c.Render(code, render.MsgPack{Data: obj})`来渲染

注意,当调用了上面这些方法后(实际是调用了`c.Render`后),响应就设置完成了,因此再对其设置header什么的就无效了.

## 中间件

我们上面已经介绍了中间件的定义和基本原理了,但这里有一个问题--如果我们的所有路由都需要用到同一个中间件怎么办?显然每个都在绑定处理过程的函数里写上是相当丑陋的方法,gin提供了一个快捷方式--`gin.Engine.Use(middleware ...gin.HandlerFunc)`/`gin.RouterGroup.Use(middleware ...gin.HandlerFunc)`接口.

通常我们不太会直接构造一个中间件,而是做一个工厂函数来创建中间件.

它的作用是在其下所属所有的路由上绑定中间件,至于执行顺序,我们看下面的例子:

```golang
package main

import (
    "log"

    "github.com/gin-gonic/gin"
)

const (
    ADDRESS = "0.0.0.0:5000"
)

//UselessMiddlewareFactory 没什么用的中间件
func UselessMiddlewareFactory(name string) gin.HandlerFunc {

    return func(ctx *gin.Context) {

        log.Println("Hello ", name)
        ctx.Writer.Header().Set("Author", name)
        ctx.Next()
        log.Println("bye ", name)
    }
}

func main() {
    router := gin.Default()
    router.Use(UselessMiddlewareFactory("router step1"))
    router.Use(UselessMiddlewareFactory("router step2"))
    group1 := router.Group("/group1")
    group1.Use(UselessMiddlewareFactory("group1 step1"))
    group1.Use(UselessMiddlewareFactory("group1 step2"))
    {
        group1.GET("/ping",
            UselessMiddlewareFactory("ping step1"),
            UselessMiddlewareFactory("ping step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "pong",
                })
            })
        group1.GET("/pong",
            UselessMiddlewareFactory("pong step1"),
            UselessMiddlewareFactory("pong step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "ping",
                })
            })
    }
    router.Run(ADDRESS) // listen and serve on 0.0.0.0:5000
}
```

启动后请求`http://localhost:5000/group1/ping`得到结果

```txt
Hello  router step1
Hello  router step2
Hello  group1 step1
Hello  group1 step2
Hello  ping step1
Hello  ping step2
bye  ping step2
bye  ping step1
bye  group1 step2
bye  group1 step1
bye  router step2
bye  router step1
```

可以看出:

1. 整个请求-响应的流程就是从外层向内存,然后再从内存向外层执行.而`Use`接口则是先绑定先执行.
2. 构造完成响应后请求-响应的流程并不会终止,必须每个过程都执行完毕才会最终返回响应

第二点特性其实很多时候并不符合我们的要求,如果我们希望直接中断后面的执行流程怎么办呢?--调用如下接口可以实现中断执行流程:

+ `c.Abort()`中断继续传递到后面的执行过程,注意当前执行过程还是会被执行完

+ `c.AbortWithStatus(code int)`写好http状态码并中断继续执行后面的过程,注意当前执行过程还是会被执行完

+ `c.AbortWithStatusJSON(code int, jsonObj interface{})`写好http状态码和返回的json格式的body并中断继续执行后面的过程,注意当前执行过程还是会被执行完

+ `c.AbortWithError(code int, err error) *Error`返回error给`c.Errors`并中断继续执行后面的过程,注意当前执行过程还是会被执行完

由于这几个方法都会中断向后执行,所以一般也就中间件用.

### 常用中间件

1. [github.com/gin-contrib/cors](https://github.com/gin-contrib/cors)是一个管理跨域问题的中间件,一般挂在第一个加载
2. [github.com/toorop/gin-logrus](https://github.com/toorop/gin-logrus)是一个基于logrus的gin中间件,用于打印log,一般放在第二个加载
3. [github.com/gin-contrib/static"](https://github.com/gin-contrib/static)是一个管理静态http资源的中间件,一般挂在定义其他路由之前,防止被覆盖

## 优雅的关闭服务程序

go以并行见长,很有可能我们中断服务时还有任务还在执行,因此一般我们不会直接用默认提供的服务方法,而是借助标准库`net/http`来启动服务,同时转发系统信号并设置一个等待时间,等待各种`goroutine`都退出了再正真关闭服务.

```golang
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gin-gonic/gin"
)

const (
    ADDRESS = "0.0.0.0:5000"
)

//UselessMiddlewareFactory 没什么用的中间件
func UselessMiddlewareFactory(name string) gin.HandlerFunc {

    return func(ctx *gin.Context) {

        log.Println("Hello ", name)
        ctx.Writer.Header().Set("Author", name)
        ctx.Next()
        log.Println("bye ", name)
    }
}

func runserv(router *gin.Engine) {
    srv := &http.Server{
        Addr:    ADDRESS,
        Handler: router,
    }

    go func() {
        // 服务连接
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
    log.Println("Server exiting")
}

func main() {
    router := gin.Default()
    router.Use(UselessMiddlewareFactory("router step1"))
    router.Use(UselessMiddlewareFactory("router step2"))
    group1 := router.Group("/group1")
    group1.Use(UselessMiddlewareFactory("group1 step1"))
    group1.Use(UselessMiddlewareFactory("group1 step2"))
    {
        group1.GET("/ping",
            UselessMiddlewareFactory("ping step1"),
            UselessMiddlewareFactory("ping step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "pong",
                })
            })
        group1.GET("/pong",
            UselessMiddlewareFactory("pong step1"),
            UselessMiddlewareFactory("pong step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "ping",
                })
            })
    }
    runserv(router) // listen and serve on 0.0.0.0:5000
}
```

## https支持

使用标准库的另一好处是可以支持https协议,方法就是将`srv.ListenAndServe`改为`srv.ListenAndServeTLS`

```golang
package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/gin-gonic/gin"
)

const (
    ADDRESS         = "0.0.0.0:5000"
    Serv_Cert_Path  = ""
    Serv_Key_Path   = ""
    Ca_Cert_Path    = ""
    Client_Crl_Path = ""
)

//UselessMiddlewareFactory 没什么用的中间件
func UselessMiddlewareFactory(name string) gin.HandlerFunc {

    return func(ctx *gin.Context) {

        log.Println("Hello ", name)
        ctx.Writer.Header().Set("Author", name)
        ctx.Next()
        log.Println("bye ", name)
    }
}

func runserv(router *gin.Engine) {
    srv := &http.Server{
        Addr:    ADDRESS,
        Handler: router,
    }
    usetls := false
    if Serv_Cert_Path != "" && Serv_Key_Path != "" {
        usetls = true
        //双向认证
        if Ca_Cert_Path != "" {
            capool := x509.NewCertPool()
            caCrt, err := ioutil.ReadFile(Ca_Cert_Path)
            if err != nil {
                log.Println("read pem file error")
                os.Exit(2)
            }
            capool.AppendCertsFromPEM(caCrt)
            tlsconf := &tls.Config{
                RootCAs:    capool,
                ClientAuth: tls.RequireAndVerifyClientCert, // 检验客户端证书
            }
            //指定client名单
            if Client_Crl_Path != "" {
                clipool := x509.NewCertPool()
                cliCrt, err := ioutil.ReadFile(Client_Crl_Path)
                if err != nil {
                    log.Println("read pem file error")
                    os.Exit(2)
                }
                clipool.AppendCertsFromPEM(cliCrt)
                tlsconf.ClientCAs = clipool
            }
            srv.TLSConfig = tlsconf
        }
    }

    //启动服务器
    go func() {
        // 服务连接
        log.Println("servrt start")
        var err error
        if usetls {
            err = srv.ListenAndServeTLS(Serv_Cert_Path, Serv_Key_Path)

        } else {
            err = srv.ListenAndServe()
        }
        if err != nil && err != http.ErrServerClosed {
            log.Println("listen error")
            os.Exit(2)
        }
    }()

    // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
    log.Println("Server exiting")
}

func main() {
    router := gin.Default()
    router.Use(UselessMiddlewareFactory("router step1"))
    router.Use(UselessMiddlewareFactory("router step2"))
    group1 := router.Group("/group1")
    group1.Use(UselessMiddlewareFactory("group1 step1"))
    group1.Use(UselessMiddlewareFactory("group1 step2"))
    {
        group1.GET("/ping",
            UselessMiddlewareFactory("ping step1"),
            UselessMiddlewareFactory("ping step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "pong",
                })
            })
        group1.GET("/pong",
            UselessMiddlewareFactory("pong step1"),
            UselessMiddlewareFactory("pong step2"),
            func(ctx *gin.Context) {
                ctx.JSON(200, gin.H{
                    "message": "ping",
                })
            })
    }
    runserv(router) // listen and serve on 0.0.0.0:5000
}
```

## 提高性能

通常情况下gin默认执行性能就已经可以满足大部分需求了,但如果想进一步提高性能,我们还有如下手段:

1. 项目编译时加上`-tags=jsoniter`,gin默认使用的是标准库的json,由于标准库的json性能很一般,我们可以用这种方式使用`jsoniter`替代的
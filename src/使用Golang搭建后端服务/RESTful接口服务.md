# 使用Go构建RESTful接口服务

Go有很多优秀的http服务框架,目前用的最广的是[Gin](https://github.com/gin-gonic),本文也将基于Gin来讲如何使用go语言构建RESTful接口服务.



先来一个[helloworld](https://github.com/hsz1273327/TutorialForGoLang/tree/master/src/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c0)简单了解下通常go语言项目的结构和gin的基本用法

## helloworld

一个最基本的项目大致是这样一个结构:

```bash
|--assets--|
|          |--server //编译好的服务执行程序
|
|--go.mod //模块信息描述依赖
|--helloworld.go //源码
|--makefile //编译流程控制
```

+ helloworld.go

```go
package main

import (
	"github.com/gin-gonic/gin"
)

const (
	ADDRESS = "0.0.0.0:5000"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}
```

+ go.mod

```
module c0

require (
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

go 1.12
```

+ makefile

```makefile
ASSETS=assets

server: helloworld.go
	go build -o $(ASSETS)/server

```

如果我们要为其添加组件,应该以子模块的形式添加.

## 为gin写一个插件

Gin的结构类似[koa](https://tutorialforjavascript.github.io/%E4%BD%BF%E7%94%A8Javascript%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1.html#koa%E7%9A%84%E5%8E%9F%E7%90%86)也是洋葱皮结构,它使用`ctx.Next()`来区分请求和响应,一个典型的的插件如下[c1](https://github.com/hsz1273327/TutorialForGoLang/tree/master/src/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c1)

这个插件的作用是给所有的response的headers体内加上`author:hsz`

我们在helloworld的基础上个增加一个插件子模块:`middleware`

+ middleware/useless.go

```go
package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

//UselessMiddleware 没什么用的中间件
func UselessMiddleware(name string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		log.Println("Hello ", name)
		ctx.Writer.Header().Set("Author", name)
		ctx.Next()
		log.Println("bye ", name)
	}
}
```

使用这个插件只需要使用`Use`方法即可.

```go
package main

import (
	"c1/middleware"

	"github.com/gin-gonic/gin"
)

const (
	ADDRESS = "0.0.0.0:5000"
)

func main() {
	router := gin.Default()
	somemid := middleware.UselessMiddleware("hsz")
	router.Use(somemid)
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}
```

一般来说用gin不太需要用第三方插件,要用插件也是自己写为主.我们之前只用了一个第三方插件[github.com/toorop/gin-logrus](https://github.com/toorop/gin-logrus),因为我们的业务是将log打在stdout,通过docker管理服务,并且使用`ELK`工具组收集log数据.而这套需要输出的格式为json,所以我们才使用了这个插件配合[logrus](github.com/sirupsen/logrus)实现json格式输出.


## 路由组

gin对路由分组的支持是原生的并不需要插件,这就方便我们基于资源划分模块构建RESTful接口了,下面[C2](https://github.com/hsz1273327/TutorialForGoLang/tree/master/src/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/c2)是一个典型的RESTful接口,用于描述用户.

+ 路由

```go

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginlogrus.Logger(logger.Logger), gin.Recovery())
	userlistsource := sources.UserListSource{}
	usersource := sources.UserSource{}
	user := router.Group("/user")
	{
		user.POST("/", userlistsource.Post)
		user.GET("/:id", usersource.Get)
		user.PUT("/:id", usersource.Put)
		user.DELETE("/:id", usersource.Delete)
	}

	userpasswordsource := sources.UserPassWordSource{}
	userpassword := router.Group("/user-password")
	{
		userpassword.POST("/", userpasswordsource.Post)
		userpassword.PUT("/:id", userpasswordsource.Put)
	}
	userpermissionsource := sources.UserPermissionSource{}
	userpermission := router.Group("/user-permission")
	{
		userpermission.POST("/", userpermissionsource.Post)
		userpermission.DELETE("/:id", userpermissionsource.Delete)
	}

	authsource := sources.AuthSource{}
	auth := router.Group("/auth")
	{
		auth.POST("/", authsource.Post)
		auth.GET("/", authsource.Get)
	}
	permissionlistsource := sources.PermissionListSource{}
	permission := router.Group("/permission")
	{
		permission.GET("/", permissionlistsource.Get)
	}

	logger.Logger.Info("start @", ADDRESS)
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}
```

上面我们使用`router.Group`划分资源,并将对应资源对象的对应方法绑定到不同的子路径和对应HTTP方法.

+ sources

一个source可以这样写

```go
//UserSource 用户资源
type UserSource struct {
}

//Get 根据id获取用户信息
func (u UserSource) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "GET"})
}

//Put 根据id和body修改用户信息
func (u UserSource) Put(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "PUT"})
}

//Delete 根据id删除用户信息
func (u UserSource) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "DELETE"})
}
```

上面我们我们使用值绑定绑定方法是因为我们的`UserSource`是一个空的结构.



<!-- ## 一个完整的例子

一个基于JWT的用户系统,这个系统用于实现用户的注册,登录注销等操作.我们用pg作为数据库,使用gorm作为连接数据库的orm.同时为了方便其他需要登录消息的服务知道,我们会将成功的注册和登录信息通过redis广播出去.

这个例子在[]()我们通过这个项目来介绍如何使用go语言搭建一个安全可扩展的用户系统微服务组件. -->
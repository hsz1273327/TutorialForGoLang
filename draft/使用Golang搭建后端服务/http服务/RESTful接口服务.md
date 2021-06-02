# RESTful接口服务

http接口服务中最常见的就是RESTful接口服务,不了解RESTful接口的可以先去我的[这篇博文](https://blog.hszofficial.site/recommend/2019/03/14/RESTful%E9%A3%8E%E6%A0%BC%E7%9A%84%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1/)补补课.

一般来说一个RESTful接口服务需要有如下要素:

1. 能够提供满足RESTful接口规范的服务
2. 能够提供OpenAPI规范规定的接口文档
3. 能够提供swagger UI来展示OpenAPI规范的接口文档,并提供调试功能

gin用于构造RESTful接口服务的资本在于

1. 原生的`gin.RouterGroup`利于组织接口
2. 有工具可以生成OpenAPI 2规范的接口文档[https://github.com/swaggo/swag/blob/master/README_zh-CN.md]
3. 有库`github.com/swaggo/files`和`github.com/swaggo/gin-swagger`可以`swaggo/swag`生成的文件直接构造swagger UI

## 使用`gin.RouterGroup`组织资源

RESTful语境下通常我们管一种对象叫做资源,,一类相关资源叫命名空间(namespace),gin对资源支持的并不好,所以我们只能按命名空间组织,通常一个命名空间包括两个资源

+ `xxSource`,单个xx资源
+ `xxSourceList`,管理xx资源整体的资源(资源容器)

我们以一个简单的namespace--`user`为例,用gin来实现一个完整的RESTful服务

### 代码组织结构

在开始编码前我们来规定下代码的组织基本形式

```bash
项目名\
      |-apis\
      |     |-xxxnamespace  # 子模块用于实现特定命名空间,这个子模块需要实现函数`func Init(group *gin.RouterGroup)`用于将这个命名空间上的所有路由注册到gin的路由集合上
      |     |-apis.go # 实现函数`func Init(app *gin.Engine) *gin.Engine`用于通过路由集合将所有命名空间注册到gin的实例上
      |
      |-docs #用于构造swagger文档
      |-middlewares # 用于实现中间件和中间件工厂函数
      |-models # 用于构造和实现业务的数据模型
      |-serv # 用于构造服务的启动对象的模块
      |-servtest # 用于构造服务测试
      |-main.go # 入口文件
      |-...
```

这个项目的实现在[]()

使用gin构造RESTful服务最关键的地方在于组织命名空间和资源.这里给出一个个人的最佳实践:

1. 每个namespace都构造一个子模块,一个namespace使用一个`gin.RouterGroup`维护.
2. apis模块实现一个函数`func Init(app *gin.Engine) *gin.Engine`用于将所有的namespace对应的`RouterGroup`注册到`gin.Engine`

    ```bash
    // Init 初始化路由
    func Init(app *gin.Engine) *gin.Engine {
        // 注册api路由
        // 用户信息
        user := app.Group("/v1_0_0/api/user")
        usernamespace.Init(user)
        return app
    }
    ```

3. 每个namespace中都要实现一个函数`func Init(group *gin.RouterGroup)`用于将这个子模块下的资源的所有路由注册到路由组

    ```bash
    func Init(group *gin.RouterGroup) {
        us := UserSource{}
        uls := UserListSource{}
        group.GET("", uls.Get)
        group.POST("", uls.Post)
        group.GET("/:uid", us.Get)
        group.PUT("/:uid", us.Put)
        group.DELETE("/:uid", us.Delete)
    }
    ```

这样在我们只要在入口函数中创建好`gin.Engine`的实例后使用`apis.Init(app)`就可以将组织好的接口路由都注册到其中.

## 使用`swag`结合注释生成OpenAPI 2规范的接口文档

## 启动swagger UI
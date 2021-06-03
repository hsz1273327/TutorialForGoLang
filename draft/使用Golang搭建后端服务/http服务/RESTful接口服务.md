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

[github.com/swaggo/swag](https://github.com/swaggo/swag)可以用于生成OpenAPI v2规范的接口文档.它的使用方法是在需要申明的地方通过注释申明,

其基本形式为`// @<key> value1 value2 ...`

大致分为如下步骤

1. 在`main`函数上添加注释申明根字段的值和info字段的值
2. 在每个路由注册的回调函数上添加注释申明这个路由的操作字段
3. 在每个请求和返回接收的结构体上添加注释声明字段的约束
4. 在根目录下执行`swag init`,这样在`docs`文件夹下就会生成与你的注释一致的OpenAPI 2.0生命文件(swagger.json和swagger.yml),同时这个docs还是个go语言模块,我们后面可以用它直接启动swagger ui服务.后面如果注释更新,执行一样的命令就可以更新文档了

### 申明根字段的值和info字段的值

注意只有声明在`main`函数上的注释才可以申明根字段的值和info字段的值

比如这个例子:

```golang
// @title tp_go_gin_complex
// @version 1.0
// @description 测试

// @contact.name hsz
// @contact.email hsz1273327@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
func main(){...}
```

这部分常见的字段有

| OpenAPI中的字段      | 是否必须 | 说明            | 例子                                                              |
| -------------------- | -------- | --------------- | ----------------------------------------------------------------- |
| `info.title`         | true     | 应用名          | `// @title Swagger Example API`                                   |
| `info.version`       | true     | 应用版本        | `// @version 1.0.0`                                               |
| `info.description`   | false    | 应用简介        | `// @description This is a sample server celler server.`          |
| `tag.name`           | false    | 标签名          | `// @tag.name This is the name of the tag`                        |
| `tag.description`    | false    | 标签说明        | `// @tag.description Cool Description`                            |
| `info.contact.name`  | false    | 联系人信息      | `// @contact.name API Support`                                    |
| `info.contact.url`   | false    | 联系人网址      | `// @contact.url http://www.swagger.io/support`                   |
| `info.contact.email` | false    | 联系人邮箱      | `// @contact.email support@swagger.io`                            |
| `info.license.name`  | true     | 协议名          | `// @license.name Apache 2.0`                                     |
| `info.license.url`   | false    | 协议url         | `// @license.url http://www.apache.org/licenses/LICENSE-2.0.html` |
| `host`               | false    | 主机名或ip+port | `// @host localhost:8080`                                         |
| `BasePath`           | false    | 基url           | `// @BasePath /api/v1`                                            |
| `schemes`            | false    | 使用的协议      | `// @schemes http https`                                          |

另外还有验证方式的声明也在这个部分,它以`securitydefinitions`开头,并会解析接下来的几行用于设置验证信息

+ `type:basic`
  
    ```golang
    // @securityDefinitions.basic BasicAuth
    ```

+ `type:apiKey`

    ```golang
    // @securityDefinitions.apikey ApiKeyAuth
    // @in header
    // @name Authorization
    ```

+ `type:oauth2`

    + `flow:implicit`

        ```golang
        // @securitydefinitions.oauth2.implicit OAuth2Implicit
        // @authorizationUrl https://example.com/oauth/authorize
        // @scope.write Grants write access
        // @scope.admin Grants read and write access to administrative information
        ```

    + `flow:password`

        ```golang
        // @securitydefinitions.oauth2.password OAuth2Password
        // @tokenUrl https://example.com/oauth/token
        // @scope.read Grants read access
        // @scope.write Grants write access
        // @scope.admin Grants read and write access to administrative information
        ```

    + `flow:application`

        ```golang
        // @securitydefinitions.oauth2.application OAuth2Application
        // @tokenUrl https://example.com/oauth/token
        // @scope.write Grants write access
        // @scope.admin Grants read and write access to administrative information
        ```

    + `flow:accessCode`

        ```golang
        // @securitydefinitions.oauth2.accessCode OAuth2AccessCode
        // @tokenUrl https://example.com/oauth/token
        // @authorizationUrl https://example.com/oauth/authorize
        // @scope.admin Grants read and write access to administrative information
        ```

### 申明路由的操作字段

这个部分需要声明到路由绑定的回调函数上,比如:

```golang
// @Summary 创建新用户
// @Tags user
// @accept json
// @Produce json
// @Param name body UserCreateQuery true "用户名"
// @Success 200 {object} user.User "{"Name":"1234","ID":1}"
// @Router /v1/api/user [post]
func (s *UserListSource) Post(c *gin.Context) 
```

下面是常用的声明字段中的元数据部分:

| OpenAPI中的字段   | 是否必须 | 说明                       | 例子                             |
| ----------------- | -------- | -------------------------- | -------------------------------- |
| `path`前两层的key | true     | 指定路由路径               | `// @Router /bottles/{id} [get]` |
| `deprecated`      | false    | 指定接口是否过时,默认false | `// @Deprecated`                 |
| `description`     | false    | 接口介绍                   | `// @Description 创建新用户`     |
| `summary`         | false    | 接口简介                   | `// @Summary 创建新用户`         |
| `operationId`     | false    | 接口唯一id                 | `// @ID 123413`                  |
| `tags`            | false    | 接口标签                   | `// @Tags a b c`                 |

需要注意`@Router`中声明的`{id}`需要在`@Param`中有对应的path类型的参数,且参数名必须一致

#### 认证信息部分

| OpenAPI中的字段 | 是否必须 | 说明                                                           | 例子                      |
| --------------- | -------- | -------------------------------------------------------------- | ------------------------- |
| `security`      | false    | 申明接口的验证方式,需要用到前面`securitydefinitions`中定义的值 | `// @Security ApiKeyAuth` |

#### 请求信息部分

| OpenAPI中的字段 | 是否必须 | 说明                                                                                | 例子                                                |
| --------------- | -------- | ----------------------------------------------------------------------------------- | --------------------------------------------------- |
| `consumes`      | false    | 指定接口接收的Mime Types                                                            | `// @Accept application/json application/msgpack`   |
| `parameters`    | false    | 声明请求数据的模式,其模式为`参数名 参数位置 数据类型 是否必须 参数说明 其他约束...` | `// @Param name body UserCreateQuery true "用户名"` |

参数位置是如下枚举:

+ `query`
+ `path`
+ `header`
+ `body`
+ `formData`

数据类型可以是如下值

+ `string`
+ `integer`
+ `number`
+ `boolean`
+ 用户自定义类型
+ `[]类型`用于声明参数为列表类型

当我们的参数位置为`query`/`path`/`header`/`formData`时,我们可以为它定义约束,根据不同的数据类型我们可以定义不同的约束,这个约束满足jsonschema定义,比如:

```golang
// @Param strq query string false "string valid" minlength(5) maxlength(10) Enums("1.1", "1.2", "1.3")
```

当我们的参数位置为`body`时,我们的数据类型必定为用户自定义类型或者列表类型,这个时候就需要通过结构体的标签来定义约束了.非常遗憾目前标签并不能使用`jsonschema`标签,它的约束是单独的如下例:

```golang
type UserCreateQuery struct {
    Name string `json:"Name" minLength:"4" maxLength:"16"`
}
```

#### 响应信息部分

| OpenAPI中的字段          | 是否必须 | 说明                                                                                  | 例子                                                      |
| ------------------------ | -------- | ------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| `produces`               | true     | 指定接口返回的Mime Types                                                              | `// @Produce application/json application/msgpack`        |
| `responses.<xxx>.header` | false    | 指定http状态码返回值下添加Header部分,顺序为`状态码 {返回类型} 数据类型 注释`          | `// @Header 200 {string} Token "qwerty"`                  |
| `responses.[default]`    | true     | 自定请求成功时的响应状态,顺序为`状态码 {返回类型} 数据类型 注释`                      | `// @Success 200 {array} Token "pong"`                    |
| `responses.<xxx>`        | false    | 自定请求失败时的响应状态,顺序为`状态码 {返回类型} 数据类型 注释`                      | `// @Failure 404 {string} NotFoundResponse "未找到资源"`  |
| `responses.<xxx>`        | false    | 自定响应状态,也就是上面success和failure的整合,顺序为`状态码 {返回类型} 数据类型 注释` | `// @Response 404 {string} NotFoundResponse "未找到资源"` |

返回类型支持的有:

+ string,表示返回为纯文本字符串
+ object,表示返回指定结构体声明模式的数据
+ array,表示返回的是列表

数据类型支持的有:

+ `string`
+ `integer`
+ `number`
+ `boolean`
+ 用户自定义类型

### 启动swagger UI

前面生成docs模块我们可以借助包`github.com/swaggo/files`和`github.com/swaggo/gin-swagger`将其注册到gin.Engine实例上

```golang
import (
    _ "你的项目/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main(){
    ...
    url := ginSwagger.URL("http://localhost:5000/swagger/doc.json") // The url pointing to API definition
    app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
```

我们以5000端口启动服务器,执行后我们就可以访问<http://localhost:5000/swagger/index.html>来访问swagger ui了
需要注意一般swagger ui只会在开发阶段的测试服务器中执行,线上条件下不应将他拉起来.

## 服务测试

RESTful接口服务的服务测试相对还是比价好做的,我们可以参照下面来实现:

```golang
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "tp_go_gin_complex/apis"
    "tp_go_gin_complex/models"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

const DBURL = "sqlite://:memory:"

// 创建一个gin实例,初始化好测试数据库以供测试,一般用sqlite的内存模式以避免来回创建删除数据库
func setupRouter() *gin.Engine {
    r := gin.Default()
    r.RedirectTrailingSlash = true
    apis.Init(r)
    models.Init(DBURL)
    return r
}

// 测试接口
func TestPingRoute(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "pong", w.Body.String())
}
```
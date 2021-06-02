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

我们以



## 使用`swag`结合注释生成OpenAPI 2规范的接口文档

## 启动swagger UI
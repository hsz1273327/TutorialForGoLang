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
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}

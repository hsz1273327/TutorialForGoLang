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

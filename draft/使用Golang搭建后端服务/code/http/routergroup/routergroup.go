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

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

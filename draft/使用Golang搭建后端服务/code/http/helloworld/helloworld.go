package main

import (
	"github.com/gin-gonic/gin"
)

const (
	ADDRESS = "0.0.0.0:5000"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:5000
}

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

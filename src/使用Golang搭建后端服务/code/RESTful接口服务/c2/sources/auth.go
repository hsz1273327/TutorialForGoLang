package sources

import (
	"github.com/gin-gonic/gin"
)

//AuthSource 用户登录和验证服务
type AuthSource struct {
}

//Get 验证一个JWTtoken是否正确
func (u AuthSource) Get(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"method": "GET"})
}

//Post 生成一个新的JWTtoken.
func (u AuthSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

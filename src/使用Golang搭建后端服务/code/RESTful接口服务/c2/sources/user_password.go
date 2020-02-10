package sources

import (
	"c2/logger"

	"github.com/gin-gonic/gin"
)

//UserPassWordSource 用户密码资源
type UserPassWordSource struct {
}

//Put 根据id和body修改用户密码,这个指令触发发送邮件给用户,并附带一个会过期的token,
//用户需要使用这个token访问Post接口重置密码
func (u UserPassWordSource) Put(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "PUT"})
}

//Post 根据id和body重置用户密码
func (u UserPassWordSource) Post(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "POST"})
}

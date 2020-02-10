package sources

import (
	"github.com/gin-gonic/gin"
)

//UserListSource 用户资源列表
type UserListSource struct {
}

//Post 创建一个新用户
func (u UserListSource) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "POST"})
}

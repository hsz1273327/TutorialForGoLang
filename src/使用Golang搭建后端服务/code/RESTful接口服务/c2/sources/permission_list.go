package sources

import (
	"github.com/gin-gonic/gin"
)

//PermissionListSource 权限资源
type PermissionListSource struct {
}

//Get 获取全部权限列表
func (u PermissionListSource) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"method": "GET"})
}

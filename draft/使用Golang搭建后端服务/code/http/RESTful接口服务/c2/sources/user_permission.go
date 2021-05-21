package sources

import (
	"c2/logger"

	"github.com/gin-gonic/gin"
)

//UserPermissionSource 用户密码资源
type UserPermissionSource struct {
}

//Delete 根据用户id和Permission id删除用户的权限
func (u UserPermissionSource) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "Delete"})
}

//Post 根据用户id和Permission id 新建用户的权限
func (u UserPermissionSource) Post(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "POST"})
}

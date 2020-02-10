package sources

import (
	"c2/logger"

	"github.com/gin-gonic/gin"
)

//UserSource 用户资源
type UserSource struct {
}

//Get 根据id获取用户信息
func (u UserSource) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "GET"})
}

//Put 根据id和body修改用户信息
func (u UserSource) Put(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "PUT"})
}

//Delete 根据id删除用户信息
func (u UserSource) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logger.Logger.Info(id)
	ctx.JSON(200, gin.H{
		"id":     id,
		"method": "DELETE"})
}

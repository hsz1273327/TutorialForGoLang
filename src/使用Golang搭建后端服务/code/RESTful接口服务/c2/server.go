package main

import (
	"c2/logger"
	"c2/sources"

	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
)

const (
	ADDRESS = "0.0.0.0:5000"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginlogrus.Logger(logger.Logger), gin.Recovery())
	userlistsource := sources.UserListSource{}
	usersource := sources.UserSource{}
	user := router.Group("/user")
	{
		user.POST("/", userlistsource.Post)
		user.GET("/:id", usersource.Get)
		user.PUT("/:id", usersource.Put)
		user.DELETE("/:id", usersource.Delete)
	}

	userpasswordsource := sources.UserPassWordSource{}
	userpassword := router.Group("/user-password")
	{
		userpassword.POST("/", userpasswordsource.Post)
		userpassword.PUT("/:id", userpasswordsource.Put)
	}
	userpermissionsource := sources.UserPermissionSource{}
	userpermission := router.Group("/user-permission")
	{
		userpermission.POST("/", userpermissionsource.Post)
		userpermission.DELETE("/:id", userpermissionsource.Delete)
	}

	authsource := sources.AuthSource{}
	auth := router.Group("/auth")
	{
		auth.POST("/", authsource.Post)
		auth.GET("/", authsource.Get)
	}
	permissionlistsource := sources.PermissionListSource{}
	permission := router.Group("/permission")
	{
		permission.GET("/", permissionlistsource.Get)
	}

	logger.Logger.Info("start @", ADDRESS)
	router.Run(ADDRESS) // listen and serve on 0.0.0.0:8080
}

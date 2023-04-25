// Package handler
// Title       : routers.go
// Author      : Tuffy  2023/4/4 15:56
// Description :
package handler

import (
	"TrsFrame/src/api/internal/config"
	"TrsFrame/src/api/internal/handler/home"
	"github.com/gin-gonic/gin"
)

func InitRouter(apiConfig *config.ApiConfig) (rootGin *gin.Engine) {
	rootGin = gin.New()
	rootGin.Use(gin.Logger(), gin.Recovery())

	//rootGin.Use(middleware.ExceptionCaptureMiddleware)

	home.InitHomeGroup(rootGin.Group(apiConfig.RootURL), apiConfig)

	rootGin.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, apiConfig.RootURL) // 301永久重定向
		//ctx.Request.URL.Path = project_settings.RootURL // 路由重定向
		//rootGin.HandleContext(ctx)
	})
	return rootGin
}

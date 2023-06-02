// Package internal
// Title       : routers.go
// Author      : Tuffy  2023/4/4 15:56
// Description :
package internal

import (
	"TrsFrame/src/api/config"
	"TrsFrame/src/api/internal/index"
	"github.com/gin-gonic/gin"
)

func InitRouter(rootGin *gin.Engine, apiConfig *config.ApiConfig) {

	index.InitIndexGroup(rootGin.Group(apiConfig.RootURL), apiConfig)

	rootGin.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, apiConfig.RootURL) // 301永久重定向
		//ctx.Request.URL.Path = project_settings.RootURL // 路由重定向
		//rootGin.HandleContext(ctx)
	})
}

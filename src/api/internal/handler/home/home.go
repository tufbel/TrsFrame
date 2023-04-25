// Package home
// Title       : home.go
// Author      : Tuffy  2023/4/4 16:08
// Description :
package home

import (
	"TrsFrame/src/api/internal/config"
	"TrsFrame/src/api/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Home
// @Summary      home
// @Description  验活接口
// @Tags         index
// @Accept       json
// @Produce      json
// @Success      200  {object} types.HomeResp
// @Failure      400
// @Router       / [get]
func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &types.HomeResp{
		Message: "success",
	})
}

func InitHomeGroup(group *gin.RouterGroup, config *config.ApiConfig) {
	group.GET("", Home)
}

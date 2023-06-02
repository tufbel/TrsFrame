// Package index
// Title       : index.go
// Author      : Tuffy  2023/4/4 16:08
// Description :
package index

import (
	"TrsFrame/src/api/config"
	"TrsFrame/src/tools/mylog"
	"github.com/gin-gonic/gin"
	"net/http"
)

const PackageName = "index"

// Home
//
//	@Summary		index
//	@Description	验活接口
//	@Tags			index
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HomeResp
//	@Failure		400
//	@Router			/ [get]
func Home(ctx *gin.Context) {
	mylog.Logger.Debug("Visit Home success")
	ctx.JSON(http.StatusOK, &HomeResp{
		Message: "success",
	})
}

func InitIndexGroup(group *gin.RouterGroup, config *config.ApiConfig) {
	group.GET("", Home)
}

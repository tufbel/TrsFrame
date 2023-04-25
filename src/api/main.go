package main

import (
	_ "TrsFrame/src/api/docs"
	"TrsFrame/src/api/internal/config"
	"TrsFrame/src/api/internal/handler"
	"TrsFrame/src/tools/mylog"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	gSwag "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title           TrsFrame
// @version         1.0
// @description     TrsFrame RESTful API

// @host      localhost:22887
// @BasePath  /api/trsframe

// @schemes http https
func main() {
	apiConfig := &config.ApiConfig{
		RootURL:       "/api/trsframe",
		ListeningHost: "localhost",
		ListeningPort: 22887,
	}
	webRouter := handler.InitRouter(apiConfig)

	// 添加swagger文档
	{
		//docs.SwaggerInfo.Title = "TrsFrame"
		//docs.SwaggerInfo.Description = "TrsFrame RESTful API"
		//docs.SwaggerInfo.Version = "1.0"
		//docs.SwaggerInfo.Host = "localhost:22887"
		//docs.SwaggerInfo.BasePath = "/api/study_gin"
		//docs.SwaggerInfo.Schemes = []string{"http", "https"}
		webRouter.GET(apiConfig.RootURL+"/docs/*any", gSwag.WrapHandler(swagFiles.Handler))
		docsIndexFunc := func(ctx *gin.Context) { ctx.Redirect(301, apiConfig.RootURL+"/docs/index.html") }
		webRouter.GET(apiConfig.RootURL+"/docs", docsIndexFunc)
		mylog.Debug(fmt.Sprintf(
			"Docs: http://localhost:%d%s/docs/index.html",
			apiConfig.ListeningPort,
			apiConfig.RootURL,
		))
	}

	{
		srv := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", apiConfig.ListeningHost, apiConfig.ListeningPort),
			Handler: webRouter.Handler(),
		}

		// startup

		go func() {
			// Listen and Server
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				mylog.Error(fmt.Sprintf("Listen: %s", err))
			}
		}()

		{
			// wait signal
			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			sig := <-quit
			mylog.Info(fmt.Sprintf("Shutdown Server ... Reason: %s\n", sig))
		}

		// shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			mylog.Error(fmt.Sprintf("Server Shutdown: %s", err))
			os.Exit(1)
		}
		mylog.Info("Server exiting")
	}
}

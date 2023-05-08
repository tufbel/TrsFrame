package main

import (
	"TrsFrame/src/api/config"
	"TrsFrame/src/api/docs"
	"TrsFrame/src/api/internal"
	"TrsFrame/src/tools/mylog"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

// @title           TrsFrame
// @version         1.0.0
// @description     TrsFrame RESTful API

// @host      localhost:22887
// @BasePath  /api/trsframe

// @schemes http https
func main() {
	workdir, err := os.Getwd()
	if err != nil {
		mylog.Error(fmt.Sprintf("Failed to get working directory:", err))
		os.Exit(1)
	}

	apiConfig := &config.ApiConfig{
		WorkDir:       workdir,
		RootURL:       "/api/rcu_client",
		ListeningHost: "localhost",
		ListeningPort: 22887,
	}

	//gin.SetMode(gin.ReleaseMode)

	webRouter := gin.New()
	webRouter.Use(gin.Logger(), gin.Recovery())
	//rootGin.Use(middleware.ExceptionCaptureMiddleware)

	internal.InitRouter(webRouter, apiConfig)

	// 添加swagger文档
	{
		// gin-swagger文档
		//gDocs.SwaggerInfo.Version = "1.0.0"
		//webRouter.GET(apiConfig.RootURL+"/docs/*any", gSwag.WrapHandler(swagFiles.Handler))
		//docsIndexFunc := func(ctx *gin.Context) { ctx.Redirect(301, apiConfig.RootURL+"/docs/index.html") }
		//webRouter.GET(apiConfig.RootURL+"/docs", docsIndexFunc)
		//webRouter.GET("/docs", docsIndexFunc)
		//mylog.Debug(fmt.Sprintf(
		//	"Docs: http://localhost:%d%s/docs/index.html",
		//	apiConfig.ListeningPort,
		//	apiConfig.RootURL,
		//))

		// fast-swagger文档
		fastSwagger := &docs.FastSwagger{
			BaseURL:     apiConfig.RootURL,
			SwaggerPath: filepath.Join(apiConfig.WorkDir, "src/api/docs/swagger.json"),
			OpenapiPath: filepath.Join(apiConfig.WorkDir, "src/api/docs/openapi.json"),
		}
		fastSwagger.BuildOpenapi()
		fastSwagger.AddDocs(webRouter)
		mylog.Debug(fmt.Sprintf(
			"Docs: http://localhost:%d%s/docs",
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
				os.Exit(1)
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

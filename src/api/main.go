// Package api
// Title       : main.go
// Author      : Tuffy  2023/5/5 14:33
// Description :
package main

import (
	"TrsFrame/src/api/config"
	"TrsFrame/src/api/docs"
	"TrsFrame/src/api/internal"
	"TrsFrame/src/tools/mylog"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

//	@title			TrsFrame
//	@version		1.0.0
//	@description	TrsFrame RESTful API

//	@host		localhost:22887
//	@BasePath	/api/trsframe

// @schemes	http https
func main() {
	var workdir string
	{
		// 确定工作目录
		exePath, err := os.Executable()
		if err != nil {
			mylog.Error(fmt.Sprintf("Failed to get working directory: %v", err))
			panic(err)
		}
		workdir = filepath.Dir(exePath)
		if filepath.Base(workdir) == "build" {
			workdir = filepath.Dir(workdir)
		}
		if err := os.Chdir(workdir); err != nil {
			mylog.Error(fmt.Sprintf("Failed to switch the working to %s: %v", workdir, err))
			panic(err)
		}

		mylog.Info(fmt.Sprintf("Working directory: %s", workdir))
	}

	viperConfig := viper.New()
	{
		//	读取项目配置
		viperConfig.SetConfigName("projectConfig")
		viperConfig.SetConfigFile(filepath.Join(workdir, "pyproject.toml"))
		viperConfig.SetConfigType("toml")
		if err := viperConfig.ReadInConfig(); err != nil {
			mylog.Error(fmt.Sprintf("Failed to read pyproject.toml: %v", err))
			panic(err)
		}
	}

	apiConfig := &config.ApiConfig{
		WorkDir:       workdir,
		RootURL:       "/api/trsframe",
		ListeningHost: viperConfig.GetString("myproject.listening_host"),
		ListeningPort: viperConfig.GetInt("myproject.listening_port"),
	}
	if apiConfig.ListeningPort == 0 {
		panic("listening_port is not set in pyproject.toml")
	}

	//gin.SetMode(gin.ReleaseMode)

	webRouter := gin.New()
	webRouter.Use(gin.Logger(), gin.Recovery())
	//rootGin.Use(middleware.ExceptionCaptureMiddleware)

	internal.InitRouter(webRouter, apiConfig)

	// 添加OpenAPI文档文档
	{
		fastOpenAPI := &docs.FastOpenAPI{
			Title:           "TrsFrame",
			BaseURL:         apiConfig.RootURL,
			ApiDir:          filepath.Join(apiConfig.WorkDir, "src/api"),
			SwaggerFileName: "swagger.json",
			OpenapiFileName: "openapi.json",
		}
		go fastOpenAPI.BuildOpenapi()
		fastOpenAPI.AddDocs(webRouter)
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

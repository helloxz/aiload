package router

import (
	"aiload/api"
	"aiload/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Start() {
	// gin 运行模式
	RunMode := viper.GetString("server.mode")
	// 设置运行模式
	gin.SetMode(RunMode)
	// 运行 gin
	r := gin.Default()
	// 全局跨域中间件
	r.Use(middleware.CORSMiddleware())
	r.GET("/", api.Home)

	r.Any("/v1/chat/completions", middleware.Auth(), api.Relay)

	// 前台首页
	// 获取服务端配置
	port := ":" + viper.GetString("server.port")
	// 运行服务
	r.Run(port)
}

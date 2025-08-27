package main

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/server"
	"TFLanHttpDesktop/internal/ui"
	"github.com/gin-gonic/gin"
	"runtime"
)

func init() {
	logger.SetLogFile("./logs/", "TFLanHttpDesktop", 7)
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	//gin.DefaultWriter = io.Discard
	server.Router = gin.Default()
}

func main() {

	// 初始化http服务
	server.InitHttpServer()

	// 初始化ui
	ui.InitUI()

}

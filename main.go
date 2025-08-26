package main

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/ui"
)

func init() {
	logger.SetLogFile("./logs/", "TFLanHttpDesktop", 7)
}

func main() {

	// 初始化http服务

	// 初始化ui
	ui.InitUI()

}

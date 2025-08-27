package main

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"TFLanHttpDesktop/internal/server"
	"TFLanHttpDesktop/internal/ui"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	logger.SetLogFile("./logs/", "TFLanHttpDesktop", 7)
	initDB()
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	//gin.DefaultWriter = io.Discard
	server.Router = gin.Default()
}

func main() {

	// 初始化http服务
	listener, err := net.Listen("tcp", "0.0.0.0:0") // 关键：绑定0.0.0.0确保外部可访问
	if err != nil {
		logger.ErrorF("创建监听器失败: %s", err.Error())
		return
	}

	addr := listener.Addr().(*net.TCPAddr)
	actualPort := addr.Port

	lanIp, _ := utils.GetLocalIP()
	logger.Info("局域网ip ", lanIp)

	define.HttpPort = actualPort
	define.LanIP = lanIp
	define.DoMain = fmt.Sprintf("http://%s:%d", define.LanIP, define.HttpPort)
	logger.InfoF("http服务启动 %s/health", define.DoMain)
	server.InitHttpServer(listener)

	// 初始化ui
	ui.InitUI()
}

func initDB() {
	logger.Debug("初始化DB")
	dbPath := resolveDBPath()
	dir := filepath.Dir(dbPath)
	logger.Debug("应用数据文件: ", dir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	data.InitDB(dbPath)
}

func resolveDBPath() string {
	switch runtime.GOOS {
	case "windows":
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			userProfile := os.Getenv("USERPROFILE")
			appDataPath = filepath.Join(userProfile, "AppData", "Roaming")
		}
		return filepath.Join(appDataPath, define.DBFileDirName, define.DBFileFileName)
	case "linux":
		home := os.Getenv("HOME")
		return filepath.Join(home, ".local", "share", define.DBFileDirName, define.DBFileFileName)
	case "darwin":
		home := os.Getenv("HOME")
		return filepath.Join(home, "Library", "Application Support", define.DBFileDirName, define.DBFileFileName)
	default:
		panic("不支持的操作系统")
	}
}

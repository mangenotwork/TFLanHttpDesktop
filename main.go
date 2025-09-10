package main

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"TFLanHttpDesktop/internal/mq"
	"TFLanHttpDesktop/internal/server"
	"TFLanHttpDesktop/internal/ui"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

// 验证系统编码环境（非root用户可能缺失UTF-8配置）
func verifyEncoding() {
	imported := map[string]string{
		"LC_ALL":   "",
		"LANG":     "",
		"LC_CTYPE": "",
	}
	for k := range imported {
		imported[k] = os.Getenv(k)
	}
	fmt.Println("系统编码环境:")
	for k, v := range imported {
		fmt.Printf("  %s=%q\n", k, v)
		if v != "zh_CN.UTF-8" && v != "en_US.UTF-8" {
			fmt.Printf("警告: %s 不是UTF-8编码，可能导致乱码\n", k)
		}
	}
}

func init() {

	logger.Debug(theme.RootConfigDir())

	dir, _ := os.Getwd()
	tf := fmt.Sprintf("%s/fonts/NotoSans-Regular.ttf", dir)
	logger.Debug(tf)
	os.Setenv("FYNE_FONT", tf)

	verifyEncoding()

	logger.Debug("FYNE_FONT = ", os.Getenv("FYNE_FONT"))
	logger.Debug("FYNE_FONT_MONOSPACE = ", os.Getenv("FYNE_FONT_MONOSPACE"))
	logger.Debug("FYNE_FONT_SYMBOL = ", os.Getenv("FYNE_FONT_SYMBOL"))

	logger.SetLogFile("./logs/", "TFLanHttpDesktop", 7)
	initDB()
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)
	gin.DefaultWriter = io.Discard
	server.Router = gin.Default()
	mq.RunMq()
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

	// 初始化ui需要的数据
	ui.InitDB()

	ui.MainApp = app.NewWithID("TFLanHttpDesktop.2025.0826")

	icon, _ := fyne.LoadResourceFromPath("./icon.png")
	ui.MainApp.SetIcon(icon)
	// 应用自定义主题（使用嵌入的字体）
	//ui.MainApp.Settings().SetTheme(newEmbeddedFontTheme())

	ui.MainWindow = ui.MainApp.NewWindow(ui.ML(ui.MLTAppTitle))
	logger.Debug("初始化UI")

	ui.LogLifecycle(ui.MainApp)
	ui.MakeTray(ui.MainApp)
	ui.InitBus()

	ui.MainWindow.Resize(fyne.NewSize(1600, 900))
	ui.MainWindow.SetMainMenu(ui.MakeMenu())
	ui.MainWindow.SetMaster()
	ui.MainWindow.SetContent(ui.MainContent())

	//notice := widget.NewRichTextFromMarkdown(ui.MLGet(ui.MLWelcomeContent))
	notice := widget.NewRichTextFromMarkdown("你好这是中文")
	//notice.Segments[2].(*widget.HyperlinkSegment).Alignment = fyne.TextAlignCenter
	dialog.ShowCustom(ui.MLGet(ui.MLWelcome), "OK", notice, ui.MainWindow)

	ui.MainWindow.ShowAndRun()
}

func initDB() {
	logger.Debug("初始化DB")
	dbPath, fcPath, ciPath := resolveDBPath()
	dir := filepath.Dir(dbPath)
	logger.Debug("应用数据文件: ", dbPath, fcPath, ciPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	data.InitDB(dbPath, fcPath, ciPath)

}

func resolveDBPath() (string, string, string) {
	switch runtime.GOOS {
	case "windows":
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			userProfile := os.Getenv("USERPROFILE")
			appDataPath = filepath.Join(userProfile, "AppData", "Roaming")
		}
		return filepath.Join(appDataPath, define.DBFileDirName, define.DBFileFileName),
			filepath.Join(appDataPath, define.DBFileDirName, define.FcDBFileFileName),
			filepath.Join(appDataPath, define.DBFileDirName, define.CiDBFileFileName)
	case "linux":
		home := os.Getenv("HOME")
		return filepath.Join(home, ".local", "share", define.DBFileDirName, define.DBFileFileName),
			filepath.Join(home, ".local", "share", define.DBFileDirName, define.FcDBFileFileName),
			filepath.Join(home, ".local", "share", define.DBFileDirName, define.CiDBFileFileName)
	case "darwin":
		home := os.Getenv("HOME")
		return filepath.Join(home, "Library", "Application Support", define.DBFileDirName, define.DBFileFileName),
			filepath.Join(home, "Library", "Application Support", define.DBFileDirName, define.FcDBFileFileName),
			filepath.Join(home, "Library", "Application Support", define.DBFileDirName, define.CiDBFileFileName)
	default:
		panic("不支持的操作系统")
	}
}

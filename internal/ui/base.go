package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"net/url"
)

var MainApp fyne.App
var MainWindow fyne.Window

func LogLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logger.Debug("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logger.Debug("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		logger.Debug("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		logger.Debug("Lifecycle: Exited Foreground")
	})
}

// MakeTray 系统托盘
func MakeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		menuItem := fyne.NewMenuItem("TFLanHttpDesktop", func() {})
		menuItem.Icon = theme.HomeIcon()
		menuItem.Action = func() {
			logger.Debug("System tray menu tapped")

			// 关键步骤1：确保窗口可见（无论之前是隐藏还是最小化）
			MainWindow.Show()

			// 关键步骤2：将窗口置于前台并获取焦点
			MainWindow.RequestFocus()

			// 可选：更新菜单项文本
			menuItem.Label = "应用已打开"
			//// 刷新菜单显示
			//if menu := menuItem.Menu; menu != nil {
			//	menu.Refresh()
			//}
		}

		domain := fyne.NewMenuItem("项目地址", func() {
		})
		domain.Action = func() {
			u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
			_ = MainApp.OpenURL(u)
		}

		trayMenu := fyne.NewMenu(
			"TFLanHttpDesktop", // 菜单名称（桌面平台显示用）
			menuItem,           // 第一个菜单项
			domain,             // 第二个菜单项
		)
		desk.SetSystemTrayMenu(trayMenu)
	}
}

var RightContainer *container.Split
var LeftContainer *container.Split

func MainContent() *container.Split {
	MemoShow()
	DownloadContainerShow()
	UploadContainerShow()

	MemoEntry.OnChanged = func(val string) {
		//logger.Debug(val)
		_, err := data.SetMemoContent(NowMemoId, val)
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
		}
	}

	// 备忘录布局
	if MemoEntryContainer == nil {
		MemoEntryContainer = container.NewBorder(nil, nil, nil, nil, nil)
	}
	LeftContainer = container.NewHSplit(ListContainer, MemoEntryContainer)
	LeftContainer.SetOffset(0.3)

	// 下载上传布局
	RightContainer = container.NewVSplit(DownloadContainer, UploadContainer)
	RightContainer.SetOffset(0.5)

	mainContent := container.NewHSplit(LeftContainer, RightContainer)
	mainContent.SetOffset(0.60) // 左侧占比20%
	return mainContent
}

func InitDB() {
	downloadData, err := data.GetDownloadData()
	if err != nil {
		logger.Error(err)
	}
	if downloadData != nil {
		NowDownloadFilePath = downloadData.Path
	}

	uploadData, err := data.GetUploadData()
	if err != nil {
		logger.Error(err)
	}
	if uploadData != nil {
		NowUploadFilePath = uploadData.Path
	}
}

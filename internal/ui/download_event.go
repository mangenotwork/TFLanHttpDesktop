package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func DownloadEvent() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if reader == nil {
			logger.Debug("Cancelled")
			return
		}

		if reader == nil {
			logger.Debug("Cancelled")
			return
		}
		defer reader.Close()

		logger.Debug("选择的文件： ", reader.URI().Path())

		err = data.SetDownloadData(&data.DownloadNow{
			Path:       reader.URI().Path(),
			IsPassword: false,
			Password:   "",
		})
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}

		NowDownloadFilePath = reader.URI().Path()
		DownloadContainerShow()
	}, MainWindow)
	fd.Resize(fyne.NewSize(900, 600))
	fd.Show()
}

func DownloadCopyUrlEvent(url string) {
	if url == "" {
		dialog.ShowError(fmt.Errorf("复制失败，链接为空"), MainWindow)
		return
	}
	clipboard := MainApp.Clipboard()
	clipboard.SetContent(url)
	dialog.ShowInformation("复制成功", "链接已复制到剪贴板!", MainWindow)
}

func DownloadDelEvent() {
	_ = data.SetDownloadData(&data.DownloadNow{
		Path:       "",
		IsPassword: false,
		Password:   "",
	})
	NowDownloadFilePath = ""
	DownloadContainerShow()
	dialog.ShowInformation("删除成功", "已删除文件对外提供的下载链接!", MainWindow)
}

func DownloadPasswordEvent(value string) {
	password := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")
	password.SetText(value)
	items := []*widget.FormItem{
		widget.NewFormItem("Password", password),
	}
	passwordDialog := dialog.NewForm("设置密码", "保存", "取消", items, func(b bool) {
		logger.Info("Please Authenticate", password.Text)
		newDownloadData := &data.DownloadNow{
			Path:       NowDownloadFilePath,
			IsPassword: false,
			Password:   "",
		}
		if password.Text != "" {
			newDownloadData.IsPassword = true
			newDownloadData.Password = password.Text
		}

		err := data.SetDownloadData(newDownloadData)
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}

		DownloadContainerShow()

	}, MainWindow)
	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

func DownloadLogEvent() {
	logger.Debug("DownloadLogEvent")

	logList, _ := data.GetDownloadLog()
	logger.Debug("logList", logList)

	content := container.NewVBox()
	for _, v := range logList {
		content.Add(widget.NewLabel(fmt.Sprintf("%s | %s| %s| %s | %s", v.Time, v.Path, v.Size, v.IP, v.UserAgent)))
	}
	downloadDialog := dialog.NewCustom("下载日志", "关闭", container.NewScroll(content), MainWindow)
	downloadDialog.Resize(fyne.NewSize(500, 600))
	downloadDialog.Show()
}

package ui

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"time"
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
		defer func() {
			_ = reader.Close()
		}()

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

		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: fmt.Sprintf("选择了对外下载文件:%s", reader.URI().Path()),
		})

		NowDownloadFilePath = reader.URI().Path()
		DownloadContainerShow()
	}, MainWindow)
	fd.Resize(fyne.NewSize(960, 700))
	fd.Show()
}

func DownloadCopyUrlEvent(url string) {
	if url == "" {
		DialogCopyErr()
		return
	}

	i, ok := define.ShareHas[url]
	if !ok {
		define.ShareId++
		define.ShareHas[url] = define.ShareId
		define.ShareMap[define.ShareId] = url
		i = define.ShareId
	}
	url = fmt.Sprintf("%s/s/%d", define.DoMain, i)

	clipboard := MainApp.Clipboard()
	clipboard.SetContent(url)
	DialogCopySuccess(url)
	return
}

func DownloadDelEvent() {
	if NowDownloadFilePath != "" {
		_ = data.SetDownloadData(&data.DownloadNow{
			Path:       "",
			IsPassword: false,
			Password:   "",
		})
		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: "删除了对外下载文件",
		})
		NowDownloadFilePath = ""
		DownloadContainerShow()
	}
	DialogDelSuccess(MLTDownloadDelSuccess)
	return
}

func DownloadPasswordEvent(value string) {
	password := widget.NewPasswordEntry()
	password.SetText(value)
	items := []*widget.FormItem{
		widget.NewFormItem(MLGet(MLTSetPassword), password),
	}
	passwordDialog := dialog.NewForm(MLGet(MLTDownloadPasswordTitle), MLGet(MLTSave), MLGet(MLTCancel), items, func(b bool) {
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
		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: fmt.Sprintf("设置了对外下载文件的密码,文件:%s", NowDownloadFilePath),
		})
		DownloadContainerShow()

	}, MainWindow)
	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

func DownloadLogEvent() {
	logger.Debug("DownloadLogEvent")

	logList, _ := data.GetDownloadLog()
	content := container.NewVBox()
	for _, v := range logList {
		content.Add(widget.NewLabel(fmt.Sprintf("%s | %s| %s| %s | %s", v.Time, v.Path, v.Size, v.IP, v.UserAgent)))
	}
	downloadDialog := dialog.NewCustom(MLGet(MLTLog), MLGet(MLTClose), container.NewScroll(content), MainWindow)
	downloadDialog.Resize(fyne.NewSize(define.Level1Width, define.Level1Height))
	downloadDialog.Show()
}

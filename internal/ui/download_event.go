package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}

func DownloadCopyUrl(url string) {
	// 获取剪贴板并设置内容
	clipboard := MainApp.Clipboard()
	clipboard.SetContent(url)
	dialog.ShowInformation("复制成功", "链接已复制到剪贴板!", MainWindow)
}

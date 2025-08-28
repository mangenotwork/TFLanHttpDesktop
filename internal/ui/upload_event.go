package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func UploadEvent() {
	dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if list == nil {
			logger.Debug("Cancelled")
			return
		}

		logger.Debug(list.String())
		out := fmt.Sprintf("选择目录 %s :\n%s", list.Name(), list.String())
		// 改为确认对话框
		dialog.ShowConfirm(
			"确认目录", // 对话框标题
			out,    // 显示的消息内容
			func(confirmed bool) { // 用户选择后的回调函数
				if confirmed {
					// 用户点击了"确定"按钮，执行相应操作
					logger.Debug("用户确认了操作 : ", list.Path())
					// 可以在这里添加确认后的逻辑，比如实际打开目录

					err = data.SetUploadData(&data.UploadNow{
						Path:       list.Path(),
						IsPassword: false,
						Password:   "",
					})
					if err != nil {
						logger.Error(err)
						dialog.ShowError(err, MainWindow)
						return
					}
					NowUploadFilePath = list.Path()
					UploadContainerShow()

				} else {
					// 用户点击了"取消"按钮
					logger.Debug("用户取消了操作")
				}
			},
			MainWindow, // 父窗口
		)

	}, MainWindow)
}

func UploadCopyUrlEvent(url string) {
	if url == "" {
		dialog.ShowError(fmt.Errorf("复制失败，链接为空"), MainWindow)
		return
	}
	clipboard := MainApp.Clipboard()
	clipboard.SetContent(url)
	dialog.ShowInformation("复制成功", "链接已复制到剪贴板!", MainWindow)
}

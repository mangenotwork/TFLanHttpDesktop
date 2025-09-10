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

func UploadEvent() {
	fd := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
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
					_ = data.SetOperationLog(&data.OperationLog{
						Time:  time.Now().Format(utils.TimeTemplate),
						Event: fmt.Sprintf("设置了接收上传路径:%s", list.Path()),
					})
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
	fd.Resize(fyne.NewSize(960, 700))
	fd.Show()
}

func UploadCopyUrlEvent(url string) {
	if url == "" {
		dialog.ShowError(fmt.Errorf("复制失败，链接为空"), MainWindow)
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
}

func UploadDelEvent() {
	_ = data.SetUploadData(&data.UploadNow{
		Path:       "",
		IsPassword: false,
		Password:   "",
	})
	NowUploadFilePath = ""
	_ = data.SetOperationLog(&data.OperationLog{
		Time:  time.Now().Format(utils.TimeTemplate),
		Event: "删除了接收上传路径",
	})
	UploadContainerShow()
	dialog.ShowInformation("删除成功", "已删除接收上传文件链接!", MainWindow)
}

func UploadPasswordEvent(value string) {
	password := widget.NewPasswordEntry()
	password.SetText(value)
	items := []*widget.FormItem{
		widget.NewFormItem("Password", password),
	}
	passwordDialog := dialog.NewForm("设置密码", "保存", "取消", items, func(b bool) {
		logger.Info("Please Authenticate", password.Text)
		newUploadData := &data.UploadNow{
			Path:       NowUploadFilePath,
			IsPassword: false,
			Password:   "",
		}
		if password.Text != "" {
			newUploadData.IsPassword = true
			newUploadData.Password = password.Text
		}

		err := data.SetUploadData(newUploadData)
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		_ = data.SetOperationLog(&data.OperationLog{
			Time:  time.Now().Format(utils.TimeTemplate),
			Event: fmt.Sprintf("设置了接收上传路径密码:%s", NowUploadFilePath),
		})
		UploadContainerShow()

	}, MainWindow)
	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

func UploadLogEvent() {
	logList, _ := data.GetUploadLog()
	content := container.NewVBox()
	for _, v := range logList {
		content.Add(widget.NewLabel(fmt.Sprintf("%s | %s| %s| %s | %s", v.Time, v.Path, v.Files, v.IP, v.UserAgent)))
	}
	downloadDialog := dialog.NewCustom("上传日志", "关闭", container.NewScroll(content), MainWindow)
	downloadDialog.Resize(fyne.NewSize(define.Level1Width, define.Level1Height))
	downloadDialog.Show()
}

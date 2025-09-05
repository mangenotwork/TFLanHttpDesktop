package ui

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// UploadContainer 接收上传目录组件
var UploadContainer = container.New(layout.NewVBoxLayout())

var NowUploadFilePath = ""

func UploadContainerShow() {
	logger.Debug("渲染上传页面 上传目录: ", NowUploadFilePath)
	uploadData, _ := data.GetUploadData()
	UploadContainer.RemoveAll()
	UploadTitle := canvas.NewText(fmt.Sprintf("接收目录: %s", NowUploadFilePath), nil)
	UploadTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	UploadContainer.Add(layout.NewSpacer())
	UploadTitle.TextSize = 16
	uploadUrl := ""
	if NowUploadFilePath != "" {

		nowMd5 := utils.GetMD5Encode(NowUploadFilePath)
		define.UploadMem[nowMd5] = NowUploadFilePath
		uploadUrl = fmt.Sprintf("%s/upload/%s", define.DoMain, nowMd5)

		UploadTitleContainer := container.NewCenter(UploadTitle)

		UploadContainer.Add(UploadTitleContainer)
		UploadContainer.Add(layout.NewSpacer())
		qrImgUpload, _ := utils.GetQRCodeIO(uploadUrl)
		readerUpload := bytes.NewReader(qrImgUpload)
		UploadQr := canvas.NewImageFromReader(readerUpload, "移动设备在同一WiFi内扫码上传")
		UploadQr.FillMode = canvas.ImageFillOriginal
		UploadContainer.Add(UploadQr)
		UploadQrText := canvas.NewText("移动设备在同一WiFi内扫码上传", nil)
		UploadQrText.TextSize = 11
		UploadQrTextContainer := container.NewCenter(UploadQrText)
		UploadContainer.Add(UploadQrTextContainer)
	} else {
		logger.Debug("当前没有选择目录")
		UploadContainer.Add(container.NewCenter(widget.NewLabel("选择目录接收上传文件")))
	}

	uploadCopy := &widget.Button{
		Text: ML(MLTCopy),
		//Icon: theme.NavigateNextIcon(),
		OnTapped: func() {
			logger.Debug("复制上传链接")
			UploadCopyUrlEvent(uploadUrl)
		},
	}
	RegisterTranslatable(MLTCopy, uploadCopy)

	UploadTool := container.NewHBox(layout.NewSpacer(),
		&widget.Button{
			Text: "指定接收上传目录",
			Icon: theme.FolderIcon(),
			OnTapped: func() {
				logger.Debug("指定接收上传目录")
				UploadEvent()
			},
		},
		//&widget.Button{
		//	Text: "复制",
		//	//Icon: theme.NavigateNextIcon(),
		//	OnTapped: func() {
		//		logger.Debug("复制上传链接")
		//		UploadCopyUrlEvent(uploadUrl)
		//	},
		//},
		uploadCopy,
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除上传")
				UploadDelEvent()
			},
		},
		&widget.Button{
			Text: "密码",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("密码管理")
				UploadPasswordEvent(uploadData.Password)
			},
		},
		//&widget.Button{
		//	Text: "限制类型",
		//	//Icon: theme.NavigateNextIcon(),
		//	OnTapped: func() {
		//		logger.Debug("密码管理")
		//		// todo ...
		//	},
		//},
		&widget.Button{
			Text: "日志",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("接收日志")
				UploadLogEvent()
			},
		},
		layout.NewSpacer())
	UploadToolContainer := container.NewCenter(UploadTool)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Add(UploadToolContainer)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Refresh()
}

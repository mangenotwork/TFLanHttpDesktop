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
	UploadTitle := canvas.NewText(ML(MLUploadDirNow, NowUploadFilePath), nil)
	RegisterTranslatable(MLUploadDirNow, UploadTitle, NowUploadFilePath)
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
		UploadQr := canvas.NewImageFromReader(readerUpload, "")
		UploadQr.FillMode = canvas.ImageFillOriginal
		UploadContainer.Add(UploadQr)
		UploadQrText := canvas.NewText(ML(MLUploadQrTip), nil)
		RegisterTranslatable(MLUploadQrTip, UploadQrText)
		UploadQrText.TextSize = 11
		UploadQrTextContainer := container.NewCenter(UploadQrText)
		UploadContainer.Add(UploadQrTextContainer)
	} else {
		logger.Debug("当前没有选择目录")
		initUploadLabel := widget.NewLabel(ML(MLUploadNotDir))
		RegisterTranslatable(MLUploadNotDir, initUploadLabel)
		UploadContainer.Add(container.NewCenter(initUploadLabel))
	}

	uploadCopy := &widget.Button{
		Text: ML(MLTCopy),
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			logger.Debug("复制上传链接")
			UploadCopyUrlEvent(uploadUrl)
		},
	}
	RegisterTranslatable(MLTCopy, uploadCopy)

	dirBtn := &widget.Button{
		Text: ML(MLSpecifyUploadDir),
		Icon: theme.FolderIcon(),
		OnTapped: func() {
			logger.Debug("指定接收上传目录")
			UploadEvent()
		},
	}
	RegisterTranslatable(MLSpecifyUploadDir, dirBtn)

	delBtn := &widget.Button{
		Text: ML(MLTDel),
		Icon: theme.DeleteIcon(),
		OnTapped: func() {
			logger.Debug("删除上传")
			UploadDelEvent()
		},
	}
	RegisterTranslatable(MLTDel, delBtn)

	pwdBtn := &widget.Button{
		Text: ML(MLTSetPassword),
		Icon: theme.VisibilityOffIcon(),
		OnTapped: func() {
			logger.Debug("密码管理")
			UploadPasswordEvent(uploadData.Password)
		},
	}
	RegisterTranslatable(MLTSetPassword, pwdBtn)

	logBtn := &widget.Button{
		Text: ML(MLTLog),
		Icon: theme.ContentPasteIcon(),
		OnTapped: func() {
			logger.Debug("接收日志")
			UploadLogEvent()
		},
	}
	RegisterTranslatable(MLTLog, logBtn)

	UploadTool := container.NewHBox(layout.NewSpacer(),
		dirBtn,
		uploadCopy,
		delBtn,
		pwdBtn,
		//&widget.Button{
		//	Text: "限制类型",
		//	//Icon: theme.NavigateNextIcon(),
		//	OnTapped: func() {
		//		logger.Debug("密码管理")
		//		// todo ...
		//	},
		//},
		logBtn,
		layout.NewSpacer())
	UploadToolContainer := container.NewCenter(UploadTool)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Add(UploadToolContainer)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Refresh()
}

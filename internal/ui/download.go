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

// DownloadContainer 下载组件
var DownloadContainer = container.New(layout.NewVBoxLayout())

// NowDownloadFilePath 当前选择的下载文件路径
var NowDownloadFilePath = ""

func DownloadContainerShow() {
	logger.Debug("渲染下载页面 下载文件: ", NowDownloadFilePath)
	downloadData, _ := data.GetDownloadData()
	DownloadContainer.RemoveAll()
	DownloadTitle := canvas.NewText(fmt.Sprintf("下载文件: %s", NowDownloadFilePath), nil)
	DownloadTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	DownloadTitle.TextSize = 16
	DownloadTitleContainer := container.NewCenter(DownloadTitle)
	DownloadContainer.Add(layout.NewSpacer())
	downloadUrl := ""
	if NowDownloadFilePath != "" {
		nowMd5 := utils.GetMD5Encode(NowDownloadFilePath)
		define.DownloadMem[nowMd5] = NowDownloadFilePath
		downloadUrl = fmt.Sprintf("%s/download/%s", define.DoMain, nowMd5)
		DownloadContainer.Add(DownloadTitleContainer)
		qrImg, _ := utils.GetQRCodeIO(downloadUrl)
		reader := bytes.NewReader(qrImg)
		DownloadQr := canvas.NewImageFromReader(reader, "移动设备在同一WiFi内扫码下载")
		DownloadQr.FillMode = canvas.ImageFillOriginal
		DownloadContainer.Add(DownloadQr)
		DownloadQrText := canvas.NewText("移动设备在同一WiFi内扫码下载", nil)
		DownloadQrText.TextSize = 11
		DownloadQrTextContainer := container.NewCenter(DownloadQrText)
		DownloadContainer.Add(DownloadQrTextContainer)
	} else {
		DownloadContainer.Add(container.NewCenter(widget.NewLabel("选择提供下载的文件")))
	}
	DownloadContainer.Add(layout.NewSpacer())

	// 选择文件
	openFileBtn := &widget.Button{
		Text: ML(MLTSelectFile),
		Icon: theme.FileIcon(),
		OnTapped: func() {
			logger.Debug("选择文件")
			DownloadEvent()
		},
	}
	RegisterTranslatable(MLTSelectFile, openFileBtn)

	// 复制按钮
	downloadCopyBtn := &widget.Button{
		Text: ML(MLTCopy),
		//Icon: theme.NavigateNextIcon(),
		OnTapped: func() {
			logger.Debug("复制下载链接")
			DownloadCopyUrlEvent(downloadUrl)
		},
	}
	RegisterTranslatable(MLTCopy, downloadCopyBtn)

	DownloadTool := container.NewHBox(layout.NewSpacer(),
		openFileBtn,
		downloadCopyBtn,
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除下载")
				DownloadDelEvent()
			},
		},
		&widget.Button{
			Text: "密码",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("设置密码")
				DownloadPasswordEvent(downloadData.Password)
			},
		},
		&widget.Button{
			Text: "日志",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("下载日志")
				DownloadLogEvent()
			},
		},
		layout.NewSpacer())
	DownloadToolContainer := container.NewCenter(DownloadTool)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Add(DownloadToolContainer)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Refresh()
}

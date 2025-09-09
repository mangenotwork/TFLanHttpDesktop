package ui

import (
	"TFLanHttpDesktop/common/define"
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
	"path"
)

// DownloadContainer 下载组件
var DownloadContainer = container.New(layout.NewVBoxLayout())

// NowDownloadFilePath 当前选择的下载文件路径
var NowDownloadFilePath = ""

func DownloadContainerShow() {
	downloadData, _ := data.GetDownloadData()
	DownloadContainer.RemoveAll()
	_, fName := path.Split(NowDownloadFilePath)
	DownloadTitle := canvas.NewText(ML(MLTDownloadTitle, fName), nil)
	RegisterTranslatable(MLTDownloadTitle, DownloadTitle, fName)
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
		DownloadQr := canvas.NewImageFromReader(reader, "")
		DownloadQr.FillMode = canvas.ImageFillOriginal
		DownloadContainer.Add(DownloadQr)
		DownloadQrText := canvas.NewText(ML(MLTDownloadQrText), nil)
		RegisterTranslatable(MLTDownloadQrText, DownloadQrText)
		DownloadQrText.TextSize = 11
		DownloadQrTextContainer := container.NewCenter(DownloadQrText)
		DownloadContainer.Add(DownloadQrTextContainer)
	} else {
		choiceDownloadLabel := widget.NewLabel(ML(MLTChoiceDownloadLabel))
		RegisterTranslatable(MLTChoiceDownloadLabel, choiceDownloadLabel)
		DownloadContainer.Add(container.NewCenter(choiceDownloadLabel))
	}
	DownloadContainer.Add(layout.NewSpacer())

	// 选择文件
	openFileBtn := &widget.Button{
		Text: ML(MLTSelectFile),
		Icon: theme.FileIcon(),
		OnTapped: func() {
			DownloadEvent()
		},
	}
	RegisterTranslatable(MLTSelectFile, openFileBtn)

	// 复制按钮
	downloadCopyBtn := &widget.Button{
		Text: ML(MLTCopy),
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			DownloadCopyUrlEvent(downloadUrl)
		},
	}
	RegisterTranslatable(MLTCopy, downloadCopyBtn)

	// 删除按钮
	downloadDelBtn := &widget.Button{
		Text: ML(MLTDel),
		Icon: theme.DeleteIcon(),
		OnTapped: func() {
			DownloadDelEvent()
		},
	}
	RegisterTranslatable(MLTDel, downloadDelBtn)

	// 设置密码
	setPasswordBtn := &widget.Button{
		Text: ML(MLTSetPassword),
		Icon: theme.VisibilityOffIcon(),
		OnTapped: func() {
			DownloadPasswordEvent(downloadData.Password)
		},
	}
	RegisterTranslatable(MLTSetPassword, setPasswordBtn)

	logBtn := &widget.Button{
		Text: ML(MLTLog),
		Icon: theme.ContentPasteIcon(),
		OnTapped: func() {
			DownloadLogEvent()
		},
	}
	RegisterTranslatable(MLTLog, logBtn)

	DownloadTool := container.NewHBox(layout.NewSpacer(),
		openFileBtn,
		downloadCopyBtn,
		downloadDelBtn,
		setPasswordBtn,
		logBtn,
		layout.NewSpacer())
	DownloadToolContainer := container.NewCenter(DownloadTool)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Add(DownloadToolContainer)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Refresh()
}

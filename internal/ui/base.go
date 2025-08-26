package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var MainApp fyne.App
var MainWindow fyne.Window

func InitUI() {
	MainApp = app.NewWithID("TFLanHttpDesktop.2025.0826")
	MainWindow = MainApp.NewWindow("TFLanHttpDesktop")
	logger.Debug("初始化UI")
	MainWindow.Resize(fyne.NewSize(1600, 900))
	MainWindow.SetMainMenu(MakeMenu())
	MainWindow.SetMaster()
	MainWindow.SetContent(MainContent())
	MainWindow.ShowAndRun()
}

var RightContainer *container.Split
var LeftContainer *container.Split

func MainContent() *container.Split {
	DownloadContainer := container.New(layout.NewVBoxLayout())
	DownloadTitle := canvas.NewText("下载文件: aaaaaa.txt", nil)
	DownloadTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	DownloadTitle.TextSize = 19
	DownloadTitleContainer := container.NewCenter(DownloadTitle)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Add(DownloadTitleContainer)
	DownloadContainer.Add(layout.NewSpacer())
	//DownloadContainer.Add(widget.NewButtonWithIcon("选择文件", theme.FileIcon(), func() {
	//	logger.Debug("选择文件")
	//}))
	qrImg, _ := utils.GetQRCodeIO("https://www.baidu.com/")
	reader := bytes.NewReader(qrImg)
	DownloadQr := canvas.NewImageFromReader(reader, "移动设备扫码下载")
	DownloadQr.FillMode = canvas.ImageFillOriginal
	DownloadContainer.Add(DownloadQr)
	DownloadQrText := canvas.NewText("移动设备扫码下载", nil)
	DownloadQrText.TextSize = 11
	DownloadQrTextContainer := container.NewCenter(DownloadQrText)
	DownloadContainer.Add(DownloadQrTextContainer)

	DownloadTool := container.NewHBox(layout.NewSpacer(),
		&widget.Button{
			Text: "选择文件",
			Icon: theme.FileIcon(),
			OnTapped: func() {
				logger.Debug("选择文件")
				// todo ...
			},
		},
		&widget.Button{
			Text: "复制链接",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("复制下载链接")
				// todo ...
			},
		},
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除下载")
				// todo ...
			},
		},
		&widget.Button{
			Text: "设置密码",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("设置密码")
				// todo ...
			},
		},
		&widget.Button{
			Text: "下载日志",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("下载日志")
				// todo ...
			},
		},
		layout.NewSpacer())
	DownloadToolContainer := container.NewCenter(DownloadTool)
	DownloadContainer.Add(layout.NewSpacer())
	DownloadContainer.Add(DownloadToolContainer)
	DownloadContainer.Add(layout.NewSpacer())

	UploadContainer := container.New(layout.NewVBoxLayout())
	UploadTitle := canvas.NewText("接收上传文件目录: /home/aaa/", nil)
	UploadTitle.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	UploadTitle.TextSize = 18
	UploadTitleContainer := container.NewCenter(UploadTitle)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Add(UploadTitleContainer)
	DownloadContainer.Add(layout.NewSpacer())
	//UploadContainer.Add(widget.NewButtonWithIcon("指定接收上传文件目录", theme.FolderIcon(), func() {
	//	logger.Debug("选择文件目录")
	//}))
	qrImgUpload, _ := utils.GetQRCodeIO("https://www.baidu.com/")
	readerUpload := bytes.NewReader(qrImgUpload)
	UploadQr := canvas.NewImageFromReader(readerUpload, "移动设备扫码上传")
	UploadQr.FillMode = canvas.ImageFillOriginal
	UploadContainer.Add(UploadQr)
	UploadQrText := canvas.NewText("移动设备扫码上传", nil)
	UploadQrText.TextSize = 11
	UploadQrTextContainer := container.NewCenter(UploadQrText)
	UploadContainer.Add(UploadQrTextContainer)

	UploadTool := container.NewHBox(layout.NewSpacer(),
		&widget.Button{
			Text: "指定接收上传目录",
			Icon: theme.FolderIcon(),
			OnTapped: func() {
				logger.Debug("指定接收上传目录")
				// todo ...
			},
		},
		&widget.Button{
			Text: "复制链接",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("复制上传链接")
				// todo ...
			},
		},
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除上传")
				// todo ...
			},
		},
		&widget.Button{
			Text: "设置密码",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("设置密码")
				// todo ...
			},
		},
		&widget.Button{
			Text: "接收日志",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("接收日志")
				// todo ...
			},
		},
		layout.NewSpacer())
	UploadToolContainer := container.NewCenter(UploadTool)
	UploadContainer.Add(layout.NewSpacer())
	UploadContainer.Add(UploadToolContainer)
	UploadContainer.Add(layout.NewSpacer())

	RightContainer = container.NewVSplit(DownloadContainer, UploadContainer)
	RightContainer.SetOffset(0.5)

	// 备忘录

	data := make([]string, 1000)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	entryLoremIpsum := widget.NewMultiLineEntry()
	entryLoremIpsum.Wrapping = fyne.TextWrapWord

	entryLoremIpsumBtn := container.NewHBox(layout.NewSpacer(),
		&widget.Button{
			Text: "刷新",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("刷新")
				// todo ...
			},
		},
		&widget.Button{
			Text: "复制链接",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("复制链接")
				// todo ...
			},
		},
		&widget.Button{
			Text: "打开二维码",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("打开二维码")
				// todo ...
			},
		},
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除")
				// todo ...
			},
		},
		&widget.Button{
			Text: "另存为txt",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("另存为txt")
				// todo ...
			},
		},
		layout.NewSpacer())

	hbox := container.NewBorder(nil, entryLoremIpsumBtn, nil, nil, entryLoremIpsum)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		//label.SetText(data[id])
		entryLoremIpsum.SetText(data[id])
		//icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		//label.SetText("Select An Item From The List")
		entryLoremIpsum.SetText(data[id])
		//icon.SetResource(nil)
	}

	ListContainerTop := container.NewVBox(
		layout.NewSpacer(),
		//container.NewCenter(widget.NewLabel("共享备忘录")),
	)
	ListContainerTop.Add(container.NewHBox(
		&widget.Button{
			Text: "共享备忘录",
			Icon: theme.ContentAddIcon(),
			OnTapped: func() {
				logger.Debug("新建备忘录")
				// todo ...
			},
		},
		&widget.Button{
			Text: "导入本地txt",
			Icon: theme.FolderOpenIcon(),
			OnTapped: func() {
				logger.Debug("导入本地txt")
				// todo ...
			},
		},
		&widget.Button{
			//Text: "打开二维码",
			Icon: theme.ViewRefreshIcon(),
			OnTapped: func() {
				logger.Debug("刷新")
				// todo ...
			},
		},
		layout.NewSpacer(),
	))
	ListContainerTop.Add(NewSearchBox())
	ListContainerTop.Add(layout.NewSpacer())
	ListContainer := container.NewBorder(ListContainerTop, nil, nil, nil, list)
	LeftContainer = container.NewHSplit(ListContainer, hbox)
	LeftContainer.SetOffset(0.3)

	mainContent := container.NewHSplit(LeftContainer, RightContainer)
	mainContent.SetOffset(0.60) // 左侧占比20%
	return mainContent
}

var NowDownloadFilePath = ""
var UploadFilePath = ""

// 创建搜索框组件
func NewSearchBox() *fyne.Container {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("请输入搜索内容...")
	entry.OnChanged = func(s string) {
		logger.Debug("搜索 ", s)
	} // 支持回车搜索
	entryContainer := container.NewStack(entry)
	return entryContainer
}

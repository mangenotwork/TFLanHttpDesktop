package ui

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var MemoEntry = widget.NewMultiLineEntry()
var MemoEntryContainer *fyne.Container
var ListContainer *fyne.Container
var MemoListContainer *fyne.Container
var NowMemoId string = ""

func MemoListShow() {
	// 备忘录
	memoList, _ := data.GetMemoList()
	dataList := make(map[int]*data.Memo)
	for i, v := range memoList {
		dataList[i] = v
	}

	if MemoListContainer == nil {
		MemoListContainer = container.NewStack()
	}
	MemoListContainer.RemoveAll()
	MemoList := widget.NewList(
		func() int {
			return len(dataList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			//logger.Info("id = ", id)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(dataList[id].Name)
		},
	)
	MemoList.OnSelected = func(id widget.ListItemID) {
		logger.Debug("id = ", id)
		logger.Debug("data = ", dataList[id])
		//MemoEntry.SetText(dataList[id].Name)
		MemoEntryContainerShow(dataList[id].Id)
	}
	MemoList.OnUnselected = func(id widget.ListItemID) {
		//MemoEntry.SetText(dataList[id].Name)
		MemoEntryContainerShow(dataList[id].Id)
	}
	MemoListContainer.Add(MemoList)
	MemoListContainer.Refresh()
}

func MemoShow() {
	logger.Debug("显示备忘录")
	MemoListShow()
	ListContainerTop := container.NewVBox(
		layout.NewSpacer(),
	)
	ListContainerTop.Add(container.NewHBox(
		&widget.Button{
			Text: "共享备忘录",
			Icon: theme.ContentAddIcon(),
			OnTapped: func() {
				logger.Debug("新建备忘录")
				NewMemoEvent(false, "")
			},
		},
		&widget.Button{
			Text: "导入本地txt",
			Icon: theme.FolderOpenIcon(),
			OnTapped: func() {
				logger.Debug("导入本地txt")
				ImportTxtEvent()
			},
		},
		&widget.Button{
			Icon: theme.ViewRefreshIcon(),
			OnTapped: func() {
				logger.Debug("刷新")
				MemoListShow()
				dialog.ShowInformation("刷新成功", "刷新成功!", MainWindow)
			},
		},
		layout.NewSpacer(),
	))
	ListContainerTop.Add(NewSearchBox())
	ListContainerTop.Add(layout.NewSpacer())
	ListContainer = container.NewBorder(ListContainerTop, nil, nil, nil, MemoListContainer)
}

func MemoEntryContainerShow(id string) {
	logger.Debug("MemoEntryContainerShow... id=", id)

	NowMemoId = id

	content, err := data.GetMemoContent(id)
	if err != nil {
		logger.Error(err)
		dialog.ShowError(err, MainWindow)
	}

	MemoEntryContainer.RemoveAll()
	MemoEntry.Wrapping = fyne.TextWrapWord
	MemoEntry.SetText(content.String())
	MemoEntry.Refresh()

	memoUrl := fmt.Sprintf("%s/memo/%s", define.DoMain, id)

	entryLoremIpsumBtn := container.NewHBox(layout.NewSpacer(),
		&widget.Button{
			Text: "刷新",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("刷新")
				newContent, newContentErr := data.GetMemoContent(NowMemoId)
				if newContentErr != nil {
					logger.Error(newContentErr)
					dialog.ShowError(newContentErr, MainWindow)
				}
				MemoEntry.SetText(newContent.String())
				MemoEntry.Refresh()
				dialog.ShowInformation("刷新成功", "刷新成功!", MainWindow)
			},
		},
		&widget.Button{
			Text: "复制链接",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("复制链接")
				CopyMemoEvent(memoUrl)
			},
		},
		&widget.Button{
			Text: "打开二维码",
			OnTapped: func() {
				logger.Debug("打开二维码")
				OpenMemoEvent(memoUrl)
			},
		},
		&widget.Button{
			Text: "删除",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("删除")
				DelMemoEvent()
			},
		},
		&widget.Button{
			Text: "编辑属性",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("编辑属性")
				NewMemoEvent(true, NowMemoId)
			},
		},
		&widget.Button{
			Text: "另存为txt",
			//Icon: theme.NavigateNextIcon(),
			OnTapped: func() {
				logger.Debug("另存为txt")
				MemoSaveToTxt()
			},
		},
		layout.NewSpacer())

	MemoEntryContainer.Add(container.NewBorder(nil, entryLoremIpsumBtn, nil, nil, MemoEntry))
	MemoEntryContainer.Refresh()
}

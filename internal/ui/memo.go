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
var NowMemoId = ""

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
	MemoListShow()
	ListContainerTop := container.NewVBox(
		layout.NewSpacer(),
	)

	addMemoBtn := &widget.Button{
		Text: ML(MLTAddMemoBtn),
		Icon: theme.ContentAddIcon(),
		OnTapped: func() {
			logger.Debug("新建备忘录")
			NewMemoEvent(false, "")
		},
	}
	RegisterTranslatable(MLTAddMemoBtn, addMemoBtn)

	importTxtBtn := &widget.Button{
		Text: ML(MLTImportTxtBtn),
		Icon: theme.FolderOpenIcon(),
		OnTapped: func() {
			logger.Debug("导入本地txt")
			ImportTxtEvent()
		},
	}
	RegisterTranslatable(MLTImportTxtBtn, importTxtBtn)

	refreshBtn := &widget.Button{
		Text: ML(MLTRefresh),
		Icon: theme.ViewRefreshIcon(),
		OnTapped: func() {
			MemoListShow()
			d := dialog.NewInformation(MLGet(MLTDialogTipTitle), MLGet(MLTRefreshSuccess), MainWindow)
			d.Resize(fyne.NewSize(260, 120))
			d.Show()
		},
	}
	RegisterTranslatable(MLTRefresh, refreshBtn)

	ListContainerTop.Add(container.NewHBox(
		addMemoBtn,
		importTxtBtn,
		refreshBtn,
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

	refreshBtn := &widget.Button{
		Text: ML(MLTRefresh),
		Icon: theme.ViewRefreshIcon(),
		OnTapped: func() {
			newContent, newContentErr := data.GetMemoContent(NowMemoId)
			if newContentErr != nil {
				logger.Error(newContentErr)
				dialog.ShowError(newContentErr, MainWindow)
			}
			MemoEntry.SetText(newContent.String())
			MemoEntry.Refresh()
			d := dialog.NewInformation(MLGet(MLTDialogTipTitle), MLGet(MLTRefreshSuccess), MainWindow)
			d.Resize(fyne.NewSize(260, 120))
			d.Show()
		},
	}
	RegisterTranslatable(MLTRefresh, refreshBtn)

	copyBtn := &widget.Button{
		Text: ML(MLTCopy),
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			CopyMemoEvent(memoUrl)
		},
	}
	RegisterTranslatable(MLTCopy, copyBtn)

	openQrBtn := &widget.Button{
		Text: ML(MLTOpenQr),
		Icon: theme.ViewFullScreenIcon(),
		OnTapped: func() {
			OpenMemoEvent(memoUrl)
		},
	}
	RegisterTranslatable(MLTOpenQr, openQrBtn)

	delBtn := &widget.Button{
		Text: ML(MLTDel),
		Icon: theme.DeleteIcon(),
		OnTapped: func() {
			DelMemoEvent()
		},
	}
	RegisterTranslatable(MLTDel, delBtn)

	editPropertiesBtn := &widget.Button{
		Text: ML(MLTEditProperties),
		Icon: theme.DocumentCreateIcon(),
		OnTapped: func() {
			NewMemoEvent(true, NowMemoId)
		},
	}
	RegisterTranslatable(MLTEditProperties, editPropertiesBtn)

	saveAsTxtBtn := &widget.Button{
		Text: ML(MLTSaveAsTxt),
		Icon: theme.DownloadIcon(),
		OnTapped: func() {
			MemoSaveToTxt()
		},
	}
	RegisterTranslatable(MLTSaveAsTxt, saveAsTxtBtn)

	entryLoremIpsumBtn := container.NewHBox(layout.NewSpacer(),
		refreshBtn,
		copyBtn,
		openQrBtn,
		delBtn,
		editPropertiesBtn,
		saveAsTxtBtn,
		layout.NewSpacer())

	title := "标题"
	memoData, err := data.GetMemoInfo(NowMemoId)
	if err != nil {
		logger.Error(err)
	} else {
		title = memoData.Name
	}
	memoTitle := &widget.Label{
		Text: title,
	}

	MemoEntryContainer.Add(container.NewBorder(container.NewCenter(memoTitle), entryLoremIpsumBtn, nil, nil, MemoEntry))
	MemoEntryContainer.Refresh()
}

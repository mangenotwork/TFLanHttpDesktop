package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewSearchBox 创建搜索框组件
func NewSearchBox() *fyne.Container {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(ML(MLTEnterSearch))
	entry.OnChanged = func(s string) {
		logger.Debug("搜索 ", s)

		if s == "" {
			MemoListShow()
			return
		}

		popupShow(s)

	} // 支持回车搜索
	RegisterTranslatable(MLTEnterSearch, entry)
	entryContainer := container.NewStack(entry)
	return entryContainer
}

func popupShow(s string) {

	fc := data.TermExtract(s)

	ciList := make([]string, 0)

	result := make([]*data.CiList, 0)

	for _, v := range fc {
		item := data.MatchCi(v.Text)
		ciList = append(ciList, item...)
	}

	ciList = utils.SliceDeduplicate[string](ciList)

	for _, v := range ciList {
		item, _ := data.GetCiList(v)
		result = append(result, item...)
	}

	//result = utils.SliceDeduplicate[*data.CiList](result)

	deduplicate := make(map[string]struct{})

	dataList := make(map[int]*data.Memo)
	i := 0
	for _, v := range result {
		if _, ok := deduplicate[v.MemoId]; !ok {
			logger.Debug("v = ", v, v.MemoId)
			memoData, err := data.GetMemoInfo(v.MemoId)
			if err != nil {
				logger.Error("获取失败: err=", err)
				continue
			}
			dataList[i] = memoData
			deduplicate[v.MemoId] = struct{}{}
			i++
		}
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

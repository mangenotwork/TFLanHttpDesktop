package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"os"
	"time"
)

var authorityMap = map[int]string{
	1: "无权限",
	2: "只读",
	3: "可读写",
}

func NewMemoEvent(isEdit bool, memoId string) {
	var err error
	oldMemoData := &data.Memo{}
	if isEdit {
		oldMemoData, err = data.GetMemoInfo(memoId)
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
	}

	name := widget.NewEntry()
	password := widget.NewPasswordEntry()
	authorityValue := 3
	authority := widget.NewRadioGroup([]string{"无权限", "只读", "可读写"}, func(value string) {
		logger.Debug(value)
		switch value {
		case "无权限":
			authorityValue = 1
		case "只读":
			authorityValue = 2
		case "可读写":
			authorityValue = 3
		}
	})
	authority.Horizontal = true
	authority.SetSelected("可读写")
	authority.Required = true
	items := []*widget.FormItem{
		{Text: "标题", Widget: name, HintText: "标题，非必填"},
		{Text: "权限", Widget: authority, HintText: "该权限只针对三方设备"},
		{Text: "密码", Widget: password, HintText: "密码，非必填"},
	}

	dialogTitle := "新建备忘录"
	dialogConfirm := "创建"
	if isEdit {
		dialogTitle = fmt.Sprintf("编辑 - %s", oldMemoData.Name)
		dialogConfirm = "保存编辑"
		name.SetText(oldMemoData.Name)
		password.SetText(oldMemoData.Password)
		authority.SetSelected(authorityMap[oldMemoData.Authority])
	}

	passwordDialog := dialog.NewForm(dialogTitle, dialogConfirm, "取消", items, func(b bool) {

		logger.Debug("name = ", name.Text)
		logger.Debug("authorityValue = ", authorityValue)
		logger.Debug("password = ", password.Text)

		if len(name.Text) == 0 {
			name.Text = time.Now().Format(utils.TimeTemplate)
		}

		if isEdit {
			// 编辑
			_, err := data.SetMemoInfo(memoId, name.Text, authorityValue, password.Text)
			if err != nil {
				logger.Error(err)
				dialog.ShowError(err, MainWindow)
				return
			}
		} else {
			// 新建
			_, err := data.NewMemo(name.Text, authorityValue, password.Text)
			if err != nil {
				logger.Error(err)
				dialog.ShowError(err, MainWindow)
				return
			}
		}

		MemoListShow()

	}, MainWindow)

	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

func ImportTxtEvent() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		if reader == nil {
			logger.Debug("Cancelled")
			return
		}

		if reader == nil {
			logger.Debug("Cancelled")
			return
		}
		defer reader.Close()

		logger.Debug(reader.URI().Path())

		content, err := os.ReadFile(reader.URI().Path())
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}

		memoData, err := data.NewMemo(fmt.Sprintf("%s - %s - 导入", reader.URI().Name(), time.Now().Format(utils.TimeTemplate)), 3, "")
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}
		_, err = data.SetMemoContent(memoData.Id, string(content))
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}
		MemoListShow()

		return

	}, MainWindow)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fd.Show()
}

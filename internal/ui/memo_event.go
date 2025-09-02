package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/common/utils"
	"TFLanHttpDesktop/internal/data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"time"
)

func NewMemoEvent() {
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

	passwordDialog := dialog.NewForm("新建备忘录", "创建", "取消", items, func(b bool) {

		logger.Debug("name = ", name.Text)
		logger.Debug("authorityValue = ", authorityValue)
		logger.Debug("password = ", password.Text)

		if len(name.Text) == 0 {
			name.Text = time.Now().Format(utils.TimeTemplate)
		}

		_, err := data.NewMemo(name.Text, authorityValue, password.Text)
		if err != nil {
			logger.Error(err)
			dialog.ShowError(err, MainWindow)
			return
		}

		MemoListShow()

	}, MainWindow)

	passwordDialog.Resize(fyne.NewSize(500, 300))
	passwordDialog.Show()
}

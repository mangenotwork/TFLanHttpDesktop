package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ComponentSwitchLangBtn 切换中英文按钮
var ComponentSwitchLangBtn = widget.NewButton(T("switchLangBtn"), func() {
	// 切换语言（会触发全局事件）
	if T("switchLangBtn") == "切换到英文" {
		SetLang("en")
	} else {
		SetLang("zh-CN")
	}
})

// ComponentDialogContainer 语言偏好弹框
var ComponentDialogContainer *dialog.CustomDialog

func NewComponentDialogContainer() *dialog.CustomDialog {
	return dialog.NewCustom(T("DialogContainerTitle"), T("DialogContainerDismiss"), container.NewCenter(ComponentSwitchLangBtn), MainWindow)
}

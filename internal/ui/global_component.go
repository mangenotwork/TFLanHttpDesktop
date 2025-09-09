package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ComponentSwitchLangBtn 切换中英文按钮
var ComponentSwitchLangBtn = widget.NewButton(ML(MLTSwitchLangBtnText), func() {
	// 切换语言（会触发全局事件）
	if ML(MLTSwitchLangBtnText) == translations["zh-CN"][MLTSwitchLangBtnText] {
		SetLang("en")
	} else {
		SetLang("zh-CN")
	}
})

// ComponentDialogContainer 语言偏好弹框
var ComponentDialogContainer *dialog.CustomDialog

func NewComponentDialogContainer() *dialog.CustomDialog {
	return dialog.NewCustom(ML(MLTDialogContainerTitle), ML(MLTDialogContainerDismiss), container.NewCenter(ComponentSwitchLangBtn), MainWindow)
}

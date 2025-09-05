package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// 自定义文字工具栏项（实现 ToolbarItem 接口）
type TextToolbarItem struct {
	*widget.Button
}

// ToolbarObject 实现 ToolbarItem 接口
func (t *TextToolbarItem) ToolbarObject() fyne.CanvasObject {
	return t.Button
}

// 快捷创建文字工具栏项
func NewTextToolbarItem(label string, onTapped func()) *TextToolbarItem {
	return &TextToolbarItem{
		Button: widget.NewButton(label, onTapped),
	}
}

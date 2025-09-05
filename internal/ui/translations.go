package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"sync"
)

const (
	MLTAppTitle               = "MLTAppTitle"               // 应用标题
	MLTWelcomeText            = "MLTWelcomeText"            // 应用欢迎语
	MLTSwitchLangBtnText      = "MLTSwitchLangBtnText"      // 切换语言按钮文本
	MLTDialogContainerTitle   = "MLTDialogContainerTitle"   // 切换语言弹框标题
	MLTDialogContainerDismiss = "MLTDialogContainerDismiss" // 切换语言弹框关闭按钮文本
	MLTCopy                   = "MLTCopy"                   // 复制文案
	MLTSelectFile             = "MLTSelectFile"             // 选择文件文案
)

// 1. 翻译映射表
var translations = map[string]map[string]string{
	"en": {
		MLTAppTitle:               "TFLanHttpDesktop",
		MLTWelcomeText:            "Welcome to TFLanHttpDesktop, Generate a QR code or link for third-party devices within the local network to download or upload specified files, a desktop application, cross-platform, and green (portable, no installation required).",
		MLTSwitchLangBtnText:      "Switch to Chinese",
		MLTDialogContainerTitle:   "Lang",
		MLTDialogContainerDismiss: "Close",
		MLTCopy:                   "Copy",
		MLTSelectFile:             "Select File",
	},
	"zh-CN": {
		MLTAppTitle:               "TFLanHttpDesktop 内网传输工具",
		MLTWelcomeText:            "欢迎使用 TFLanHttpDesktop , Transfer Files from LAN Http Desktop, 生成二维码或链接提供给局域网内的三方设备下载指定文件和上传文件，桌面应用程序，跨平台，绿色免安装",
		MLTSwitchLangBtnText:      "切换到英文",
		MLTDialogContainerTitle:   "语言偏好",
		MLTDialogContainerDismiss: "关闭",
		MLTCopy:                   "复制",
		MLTSelectFile:             "选择文件",
	},
}

// 2. 全局事件总线（用于通知语言变化）
type EventBus struct {
	listeners []func() // 语言变化时的回调函数列表
	mu        sync.Mutex
}

var bus = &EventBus{}

func InitBus() {
	bus.OnLangChanged(func() {
		// 窗口标题刷新
		MainWindow.SetTitle(ML(MLTAppTitle))

		// 语言偏好
		ComponentSwitchLangBtn.SetText(ML(MLTSwitchLangBtnText))
		if ComponentDialogContainer != nil {
			ComponentDialogContainer.Hide()
		}
		ComponentDialogContainer = NewComponentDialogContainer()
		ComponentDialogContainer.Resize(fyne.NewSize(500, 600))
		ComponentDialogContainer.Show()

	})
}

// 注册语言变化监听器
func (b *EventBus) OnLangChanged(callback func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, callback)
}

// 触发语言变化事件（通知所有监听器）
func (b *EventBus) TriggerLangChanged() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, callback := range b.listeners {
		callback() // 调用所有注册的刷新函数
	}
}

// 3. 翻译管理器
var (
	currentLang = "zh-CN"
	langMu      sync.Mutex
)

// 设置当前语言并触发事件
func SetLang(lang string) {
	langMu.Lock()
	currentLang = lang
	langMu.Unlock()

	// 通知所有组件刷新
	bus.TriggerLangChanged()

	// 遍历所有注册的组件，自动更新文本（无需手动写每个组件的刷新逻辑）
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	for key, components := range translatableComponents.list {
		text := ML(key)
		for _, comp := range components {
			switch c := comp.(type) {
			case *widget.Label:
				c.SetText(text)
			case *widget.Button:
				c.SetText(text)
			case *widget.Entry:
				c.SetPlaceHolder(text) // Entry 适配占位文本
			}
		}
	}
}

// ML 多语言 multilingual
func ML(key string) string {
	langMu.Lock()
	lang := currentLang
	langMu.Unlock()

	if trans, ok := translations[lang]; ok {
		if text, ok := trans[key]; ok {
			return text
		}
	}
	return translations["zh-CN"][key] // 默认中文
}

// 存储需要翻译的组件：key是翻译键，value是组件和更新方法
var translatableComponents = struct {
	sync.Mutex
	list map[string][]interface{} // 支持多种组件类型（Label/Button/Entry等）
}{
	list: make(map[string][]interface{}),
}

// RegisterTranslatable 注册组件翻译
// 支持 Label/Button/Entry.PlaceHolder 等常见组件
func RegisterTranslatable(key string, component interface{}) {
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	translatableComponents.list[key] = append(translatableComponents.list[key], component)
}

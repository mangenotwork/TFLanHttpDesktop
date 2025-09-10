package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
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
	MLTDownloadTitle          = "MLTDownloadTitle"          // 下载标题
	MLTDownloadQrText         = "MLTDownloadQrText"         // 下载二维码说明
	MLTChoiceDownloadLabel    = "MLTChoiceDownloadLabel"    // 选择下载提示
	MLTDel                    = "MLTDel"                    // 删除
	MLTSetPassword            = "MLTSetPassword"            // 设置密码
	MLTLog                    = "MLTLog"                    // 日志
	MLTDialogCopyLinkErr      = "MLTDialogCopyLinkErr"      // 复制失败弹出层
	MLTDialogTipTitle         = "MLTDialogTipTitle"         // 提示框title
	MLTDialogCopyLinkSuccess  = "MLTDialogCopyLinkSuccess"  // 复制成功弹出层
	MLTDownloadDelSuccess     = "MLTDownloadDelSuccess"     // 删除下载成功
	MLTDownloadPasswordTitle  = "MLTDownloadPasswordTitle"  // 下载文件设置密码弹出层标题
	MLTSave                   = "MLTSave"                   // 保存
	MLTCancel                 = "MLTCancel"                 // 取消
	MLTClose                  = "MLTClose"                  // 关闭
	MLTAddMemoBtn             = "MLTAddMemoBtn"             // 共享备忘录
	MLTImportTxtBtn           = "MLTImportTxtBtn"           // 导入txt
	MLTRefreshSuccess         = "MLTRefreshSuccess"         // 刷新成功
	MLTRefresh                = "MLTRefresh"                // 刷新
	MLTOpenQr                 = "MLTOpenQr"                 // 打开二维码
	MLTEditProperties         = "MLTEditProperties"         // 编辑属性
	MLTSaveAsTxt              = "MLTSaveAsTxt"              // 另存为txt
	MLTNoPermission           = "MLTNoPermission"           // 无权限
	MLTReadOnly               = "MLTReadOnly"               // 只读
	MLTReadWrite              = "MLTReadWrite"              // 可读写
	MLTInputTitle             = "MLTInputTitle"             // 输入标题
	MLTInputTitleHint         = "MLTInputTitleHint"         // 输入标题提示
	MLTInputAuthority         = "MLTInputAuthority"         // 选择权限
	MLTInputAuthorityHint     = "MLTInputAuthorityHint"     // 选择权限提示
	MLTInputPassword          = "MLTInputPassword"          // 输入密码
	MLTInputPasswordHint      = "MLTInputPasswordHint"      // 输入密码提示
	MLTNewMemo                = "MLTNewMemo"                // 新建备忘录
	MLTSaveEdits              = "MLTSaveEdits"              // 保存编辑
	MLTCreate                 = "MLTCreate"                 // 创建
	MLTEditsMemo              = "MLTEditsMemo"              // 编辑备忘录
	MLTScanQr                 = "MLTScanQr"                 // 扫二维码
	MLTConfirmDeletion        = "MLTConfirmDeletion"        // 确认删除吗
	MLTSaveMemoSuccess        = "MLTSaveMemoSuccess"        // 另存备忘录成功
	MLTEnterSearch            = "MLTEnterSearch"            // 输入关键词搜索内容
	MLTOpen                   = "MLTOpen"                   // 打开
	MLTSettings               = "MLTSettings"               // 设置
	MLTAbout                  = "MLTAbout"                  // 关于
	MLTLanguage               = "MLTLanguage"               // 语言
	MLTSystemLog              = "MLTSystemLog"              // 系统日志
	MLTFile                   = "MLTFile"                   // 文件
	MLTHelp                   = "MLTHelp"                   // 帮助
	MLTDocumentation          = "MLTDocumentation"          // 使用文档
	MLTProjectAddress         = "MLTProjectAddress"         // 项目地址
	MLTNewVersion             = "MLTNewVersion"             // 新版本
	MLTContactTheAuthor       = "MLTContactTheAuthor"       // 联系作者
	MLUploadDir               = "MLUploadDir"               // 接收上传目录
	MLSpecifyUploadDir        = "MLSpecifyUploadDir"        // 指定接收上传目录
	MLUploadDirNow            = "MLUploadDirNow"            // 接收目录 %s
	MLUploadQrTip             = "MLUploadQrTip"             // 移动设备在同一WiFi内扫码上传
	MLUploadNotDir            = "MLUploadNotDir"            // 选择目录接收上传文件
	MLWelcome                 = "MLWelcome"                 // 欢迎
	MLWelcomeContent          = "MLWelcomeContent"          // 欢迎内容
	MLTAboutContent           = "MLTAboutContent"           // 关于
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
		MLTDownloadTitle:          "Download File: %s",
		MLTDownloadQrText:         "Scan and download codes on mobile devices within the same WiFi",
		MLTChoiceDownloadLabel:    "Select the file to provide for download",
		MLTDel:                    "Delete",
		MLTSetPassword:            "Password",
		MLTLog:                    "log",
		MLTDialogCopyLinkErr:      "Copy failed, link is empty",
		MLTDialogTipTitle:         "Tip",
		MLTDialogCopyLinkSuccess:  "Link copied to clipboard.\nLink: %s",
		MLTDownloadDelSuccess:     "The download link provided to the public for the deleted file has been deleted",
		MLTDownloadPasswordTitle:  "Set password for downloading files",
		MLTSave:                   "Save",
		MLTCancel:                 "Cancel",
		MLTClose:                  "Close",
		MLTAddMemoBtn:             "Shared Memo",
		MLTImportTxtBtn:           "Import Txt",
		MLTRefreshSuccess:         "Refresh Success",
		MLTRefresh:                "Refresh",
		MLTOpenQr:                 "Open QR",
		MLTEditProperties:         "Edit Properties",
		MLTSaveAsTxt:              "Save As Txt",
		MLTNoPermission:           "No Permission",
		MLTReadOnly:               "Read Only",
		MLTReadWrite:              "Read Write",
		MLTInputTitle:             "Title",
		MLTInputTitleHint:         "Input Title, not required",
		MLTInputAuthority:         "Authority",
		MLTInputAuthorityHint:     "This permission is only applicable to third-party devices",
		MLTInputPassword:          "Password",
		MLTInputPasswordHint:      "Input Password, not required",
		MLTNewMemo:                "Create Memo",
		MLTSaveEdits:              "Save Edits Memo",
		MLTCreate:                 "Create",
		MLTEditsMemo:              "Edits - %s",
		MLTScanQr:                 "Scan QR",
		MLTConfirmDeletion:        "Confirm Deletion?",
		MLTSaveMemoSuccess:        "Save successfully, save to:\n%s",
		MLTEnterSearch:            "Enter keywords to search for content",
		MLTOpen:                   "Open",
		MLTSettings:               "Settings",
		MLTAbout:                  "About",
		MLTLanguage:               "Language",
		MLTSystemLog:              "System Log",
		MLTFile:                   "File",
		MLTHelp:                   "Help",
		MLTDocumentation:          "Documentation",
		MLTProjectAddress:         "Project Address",
		MLTNewVersion:             "New Version",
		MLTContactTheAuthor:       "Contact The Author",
		MLUploadDir:               "Receive Upload Directory",
		MLSpecifyUploadDir:        "Specify Upload Dir",
		MLUploadDirNow:            "Receiving directory: %s",
		MLUploadQrTip:             "Scan and upload codes on mobile devices within the same WiFi",
		MLUploadNotDir:            "Select directory to receive uploaded files",
		MLWelcome:                 "Welcome to use TFLanHttpDesktop", // 欢迎
		MLWelcomeContent:          "Transfer Files from LAN Http Desktop, This application will start an HTTP service to transfer files through the HTTP protocol.\n\n[TFLanHttpDesktop on GitHub](https://github.com/mangenotwork/TFLanHttpDesktop)",
		MLTAboutContent:           "TFLanHttpDesktop 版本：%s\nTransfer Files from LAN Http Desktop, 用于局域网内指定文件生成二维码或链接提供给三方设备用局域网http协议下载文件，三方设备也可以上传文件，桌面应用程序，跨平台。\n开源项目: https://github.com/mangenotwork/TFLanHttpDesktop\n\n",
	},
	"zh-CN": {
		MLTAppTitle:               "TFLanHttpDesktop 内网传输工具",
		MLTWelcomeText:            "欢迎使用 TFLanHttpDesktop , Transfer Files from LAN Http Desktop, 生成二维码或链接提供给局域网内的三方设备下载指定文件和上传文件，桌面应用程序，跨平台，绿色免安装",
		MLTSwitchLangBtnText:      "切换到英文",
		MLTDialogContainerTitle:   "语言偏好",
		MLTDialogContainerDismiss: "关闭",
		MLTCopy:                   "复制",
		MLTSelectFile:             "选择文件",
		MLTDownloadTitle:          "下载文件: %s",
		MLTDownloadQrText:         "移动设备在同一WiFi内扫码下载",
		MLTChoiceDownloadLabel:    "选择提供下载的文件",
		MLTDel:                    "删除",
		MLTSetPassword:            "密码",
		MLTLog:                    "日志",
		MLTDialogCopyLinkErr:      "复制失败，链接为空",
		MLTDialogTipTitle:         "提示",
		MLTDialogCopyLinkSuccess:  "链接已复制到剪贴板;\n链接: %s",
		MLTDownloadDelSuccess:     "已删除文件对外提供的下载链接",
		MLTDownloadPasswordTitle:  "为下载文件设置密码",
		MLTSave:                   "保存",
		MLTCancel:                 "取消",
		MLTClose:                  "关闭",
		MLTAddMemoBtn:             "共享备忘录",
		MLTImportTxtBtn:           "导入txt",
		MLTRefreshSuccess:         "刷新成功",
		MLTRefresh:                "刷新",
		MLTOpenQr:                 "打开二维码",
		MLTEditProperties:         "编辑属性",
		MLTSaveAsTxt:              "另存为txt",
		MLTNoPermission:           "无权限",
		MLTReadOnly:               "只读",
		MLTReadWrite:              "可读写",
		MLTInputTitle:             "标题",
		MLTInputTitleHint:         "标题，非必填",
		MLTInputAuthority:         "权限",
		MLTInputAuthorityHint:     "该权限只针对三方设备",
		MLTInputPassword:          "密码",
		MLTInputPasswordHint:      "密码，非必填",
		MLTNewMemo:                "新建备忘录",
		MLTSaveEdits:              "保存编辑",
		MLTCreate:                 "创建",
		MLTEditsMemo:              "编辑 - %s",
		MLTScanQr:                 "扫码访问",
		MLTConfirmDeletion:        "确认删除吗?",
		MLTSaveMemoSuccess:        "另存成功, 另存至:\n%s",
		MLTEnterSearch:            "输入关键词搜索内容",
		MLTOpen:                   "打开",
		MLTSettings:               "设置",
		MLTAbout:                  "关于",
		MLTLanguage:               "语言",
		MLTSystemLog:              "系统日志",
		MLTFile:                   "文件",
		MLTHelp:                   "帮助",
		MLTDocumentation:          "使用文档",
		MLTProjectAddress:         "项目地址",
		MLTNewVersion:             "新版本",
		MLTContactTheAuthor:       "联系作者",
		MLUploadDir:               "接收上传目录",
		MLSpecifyUploadDir:        "指定接收上传目录",
		MLUploadDirNow:            "接收目录: %s",
		MLUploadQrTip:             "移动设备在同一WiFi内扫码上传",
		MLUploadNotDir:            "选择目录接收上传文件",
		MLWelcome:                 "欢迎使用TFLanHttpDesktop", // 欢迎
		MLWelcomeContent:          "该应用会启动一个http服务，通过http协议进行传输文件\n\n[TFLanHttpDesktop on GitHub](https://github.com/mangenotwork/TFLanHttpDesktop)",
		MLTAboutContent:           "TFLanHttpDesktop 版本：%s\nTransfer Files from LAN Http Desktop, 用于局域网内指定文件生成二维码或链接提供给三方设备用局域网http协议下载文件，三方设备也可以上传文件，桌面应用程序，跨平台。\n开源项目: https://github.com/mangenotwork/TFLanHttpDesktop\n\n",
	},
}

// EventBus 全局事件总线（用于通知语言变化）
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

// OnLangChanged 注册语言变化监听器
func (b *EventBus) OnLangChanged(callback func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, callback)
}

// TriggerLangChanged 触发语言变化事件（通知所有监听器）
func (b *EventBus) TriggerLangChanged() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, callback := range b.listeners {
		callback() // 调用所有注册的刷新函数
	}
}

// 3. 翻译管理器
var (
	CurrentLang = "zh-CN"
	langMu      sync.Mutex
)

// SetLang 设置当前语言并触发事件
func SetLang(lang string) {
	langMu.Lock()
	CurrentLang = lang
	langMu.Unlock()

	// 通知所有组件刷新
	bus.TriggerLangChanged()

	// 遍历所有注册的组件，自动更新文本（无需手动写每个组件的刷新逻辑）
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	for key, components := range translatableComponents.list {
		for _, comp := range components {
			text := ML(key, comp.Value...)
			switch c := comp.Component.(type) {
			case *widget.Label:
				c.SetText(text)
			case *widget.Button:
				c.SetText(text)
			case *widget.Entry:
				c.SetPlaceHolder(text)
			case *canvas.Text:
				c.Text = text
				c.Refresh()
			case *fyne.MenuItem:
				c.Label = text
				if FileMenu != nil {
					FileMenu.Refresh()
				}
				if HelpMenu != nil {
					HelpMenu.Refresh()
				}

			case *fyne.Menu:
				c.Label = text
				c.Refresh()

			}
		}
	}
}

// ML 多语言 multilingual
func ML(key string, val ...any) string {
	langMu.Lock()
	lang := CurrentLang
	langMu.Unlock()

	if trans, ok := translations[lang]; ok {
		if text, ok := trans[key]; ok {

			if len(val) > 0 {
				return fmt.Sprintf(text, val...)
			}

			return text
		}
	}

	if len(val) > 0 {
		return fmt.Sprintf(translations["zh-CN"][key], val...)
	}
	return translations["zh-CN"][key] // 默认中文
}

type translatableComponentsData struct {
	Component interface{}
	Value     []any
}

// 存储需要翻译的组件：key是翻译键，value是组件和更新方法
var translatableComponents = struct {
	sync.Mutex
	list map[string][]translatableComponentsData // 支持多种组件类型（Label/Button/Entry等）
}{
	list: make(map[string][]translatableComponentsData),
}

// RegisterTranslatable 注册组件翻译
// 支持 Label/Button/Entry.PlaceHolder 等常见组件
func RegisterTranslatable(key string, component interface{}, val ...any) {
	translatableComponents.Lock()
	defer translatableComponents.Unlock()
	translatableComponents.list[key] = append(translatableComponents.list[key], translatableComponentsData{
		Component: component,
		Value:     val,
	})
}

func MLGet(mlt string, val ...any) string {
	str, ok := translations[CurrentLang][mlt]
	if !ok {
		return ""
	}
	if len(val) > 0 {
		return fmt.Sprintf(str, val...)
	}
	return str
}

func DialogCopyErr() {
	dialog.NewInformation(translations[CurrentLang][MLTDialogTipTitle], translations[CurrentLang][MLTDialogCopyLinkErr], MainWindow).Show()
}

func DialogCopySuccess(val string) {
	dialog.NewInformation(translations[CurrentLang][MLTDialogTipTitle],
		fmt.Sprintf(translations[CurrentLang][MLTDialogCopyLinkSuccess], val),
		MainWindow).Show()
}

func DialogDelSuccess(mlt string) {
	dialog.NewInformation(translations[CurrentLang][MLTDialogTipTitle], translations[CurrentLang][mlt], MainWindow).Show()
}

package ui

import (
	"TFLanHttpDesktop/common/logger"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

var FileMenu *fyne.Menu
var HelpMenu *fyne.Menu

// MakeMenu 菜单
func MakeMenu() *fyne.MainMenu {
	newItem := fyne.NewMenuItem(ML(MLTOpen), nil)
	RegisterTranslatable(MLTOpen, newItem)
	fileItem := fyne.NewMenuItem("打开图片", func() {
		//openImgFile(MainWindow)
	})
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("打开目录", func() {
		logger.Debug("Menu New->Directory")
		//openFile(MainWindow)
	})
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		dirItem,
		fileItem,
	)

	openSettings := func() {
		w := MainApp.NewWindow(MLGet(MLTSettings))
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	operationLog := func() {
		logList, _ := data.GetOperationLog()
		logger.Debug("logList", logList)
		content := container.NewVBox()
		for _, v := range logList {
			content.Add(widget.NewLabel(fmt.Sprintf("%s | %s", v.Time, v.Event)))
		}
		downloadDialog := dialog.NewCustom("下载日志", "关闭", container.NewScroll(content), MainWindow)
		downloadDialog.Resize(fyne.NewSize(500, 600))
		downloadDialog.Show()
	}
	showAbout := func() {
		w := MainApp.NewWindow(MLGet(MLTAbout))
		w.SetContent(widget.NewLabel("TFLanHttpDesktop\nTransfer Files from LAN Http Desktop, 用于局域网内指定文件生成二维码或链接提供给三方设备用局域网http协议下载文件，三方设备也可以上传文件，桌面应用程序，跨平台。"))
		w.Show()
	}

	lang := func() {
		ComponentDialogContainer = NewComponentDialogContainer()
		ComponentDialogContainer.Resize(fyne.NewSize(500, 600))
		ComponentDialogContainer.Show()
	}

	aboutItem := fyne.NewMenuItem(ML(MLTAbout), showAbout)
	RegisterTranslatable(MLTAbout, aboutItem)

	settingsItem := fyne.NewMenuItem(ML(MLTSettings), openSettings)
	RegisterTranslatable(MLTSettings, settingsItem)

	langItem := fyne.NewMenuItem(ML(MLTLanguage), lang)
	RegisterTranslatable(MLTLanguage, langItem)

	operationLogItem := fyne.NewMenuItem(ML(MLTSystemLog), operationLog)
	RegisterTranslatable(MLTSystemLog, operationLogItem)

	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	MainWindow.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	documentation := fyne.NewMenuItem(ML(MLTDocumentation), func() {
		u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
		_ = MainApp.OpenURL(u)
	})
	RegisterTranslatable(MLTDocumentation, documentation)

	projectAddress := fyne.NewMenuItem(ML(MLTProjectAddress), func() {
		u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
		_ = MainApp.OpenURL(u)
	})
	RegisterTranslatable(MLTProjectAddress, projectAddress)

	newVersion := fyne.NewMenuItem(ML(MLTNewVersion), func() {
		u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
		_ = MainApp.OpenURL(u)
	})
	RegisterTranslatable(MLTNewVersion, newVersion)

	contactTheAuthor := fyne.NewMenuItem(ML(MLTContactTheAuthor), func() {
		u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
		_ = MainApp.OpenURL(u)
	})
	RegisterTranslatable(MLTContactTheAuthor, contactTheAuthor)

	HelpMenu = fyne.NewMenu(ML(MLTHelp), documentation, projectAddress, newVersion, contactTheAuthor)
	RegisterTranslatable(MLTHelp, HelpMenu)

	FileMenu = fyne.NewMenu(ML(MLTFile), newItem)
	RegisterTranslatable(MLTFile, FileMenu)

	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		FileMenu.Items = append(FileMenu.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	FileMenu.Items = append(FileMenu.Items, operationLogItem)
	FileMenu.Items = append(FileMenu.Items, langItem)
	FileMenu.Items = append(FileMenu.Items, aboutItem)

	menu := fyne.NewMainMenu(
		FileMenu,
		HelpMenu,
	)

	return menu
}

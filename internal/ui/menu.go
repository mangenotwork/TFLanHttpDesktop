package ui

import (
	"TFLanHttpDesktop/common/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

// MakeMenu 菜单
func MakeMenu() *fyne.MainMenu {
	newItem := fyne.NewMenuItem("打开", nil)
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
		w := MainApp.NewWindow("设置")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	showAbout := func() {
		w := MainApp.NewWindow("关于")
		w.SetContent(widget.NewLabel("TFLanHttpDesktop\nTransfer Files from LAN Http Desktop, 用于局域网内指定文件生成二维码或链接提供给三方设备用局域网http协议下载文件，三方设备也可以上传文件，桌面应用程序，跨平台。"))
		w.Show()
	}
	aboutItem := fyne.NewMenuItem("关于", showAbout)
	settingsItem := fyne.NewMenuItem("设置", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	MainWindow.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	helpMenu := fyne.NewMenu("帮助",
		fyne.NewMenuItem("使用文档", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = MainApp.OpenURL(u)
		}),
		fyne.NewMenuItem("项目地址", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = MainApp.OpenURL(u)
		}),
		fyne.NewMenuItem("新版本", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = MainApp.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("联系作者", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/MyPicViu")
			_ = MainApp.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("文件", newItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	return fyne.NewMainMenu(
		file,
		helpMenu,
	)
}

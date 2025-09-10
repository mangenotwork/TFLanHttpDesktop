package ui

import (
	"TFLanHttpDesktop/common/define"
	"TFLanHttpDesktop/internal/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

var FileMenu *fyne.Menu
var HelpMenu *fyne.Menu

// MakeMenu 菜单
func MakeMenu() *fyne.MainMenu {

	openSettings := func() {
		w := MainApp.NewWindow(MLGet(MLTSettings))
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	operationLog := func() {
		logList, _ := data.GetOperationLog()

		content := container.NewVBox()
		for _, v := range logList {
			content.Add(widget.NewLabel(fmt.Sprintf("%s | %s", v.Time, v.Event)))
		}
		downloadDialog := dialog.NewCustom("系统日志", "关闭", container.NewScroll(content), MainWindow)
		downloadDialog.Resize(fyne.NewSize(define.Level1Width, define.Level1Height))
		downloadDialog.Show()
	}
	showAbout := func() {
		w := MainApp.NewWindow(MLGet(MLTAbout))
		c := container.NewVBox(layout.NewSpacer(), widget.NewLabel(
			fmt.Sprintf(MLGet(MLTAboutContent, define.Version)),
		))
		c.Add(widget.NewButton("https://github.com/mangenotwork/TFLanHttpDesktop", func() {
			u, _ := url.Parse("https://github.com/mangenotwork/TFLanHttpDesktop")
			_ = MainApp.OpenURL(u)
		}))
		c.Add(layout.NewSpacer())
		c.Resize(fyne.NewSize(500, 600))
		w.SetContent(c)
		w.Show()
	}

	lang := func() {
		ComponentDialogContainer = NewComponentDialogContainer()
		ComponentDialogContainer.Resize(fyne.NewSize(500, 600))
		ComponentDialogContainer.Show()
	}

	downloadItem := fyne.NewMenuItem(ML(MLTSelectFile), func() {
		DownloadEvent()
	})
	downloadItem.Icon = theme.FileIcon()
	RegisterTranslatable(MLTSelectFile, downloadItem)

	uploadItem := fyne.NewMenuItem(ML(MLUploadDir), func() {
		UploadEvent()
	})
	uploadItem.Icon = theme.FolderIcon()
	RegisterTranslatable(MLUploadDir, uploadItem)

	memoItem := fyne.NewMenuItem(ML(MLTAddMemoBtn), func() {
		NewMemoEvent(false, "")
	})
	memoItem.Icon = theme.ContentAddIcon()
	RegisterTranslatable(MLTAddMemoBtn, memoItem)

	importMemoItem := fyne.NewMenuItem(ML(MLTImportTxtBtn), func() {
		ImportTxtEvent()
	})
	importMemoItem.Icon = theme.FolderOpenIcon()
	RegisterTranslatable(MLTImportTxtBtn, importMemoItem)

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

	FileMenu = fyne.NewMenu(ML(MLTFile), downloadItem, uploadItem, memoItem, importMemoItem)
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

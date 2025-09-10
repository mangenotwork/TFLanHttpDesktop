package main

import (
	"TFLanHttpDesktop/common/logger"
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

//go:embed fonts/NotoSansCJKsc-Regular.otf
var fontFiles embed.FS

// 自定义主题，使用嵌入的中文字体
type embeddedFontTheme struct {
	baseTheme fyne.Theme
	fontRes   fyne.Resource // 嵌入的字体资源
}

// 初始化主题，加载嵌入的中文字体
func newEmbeddedFontTheme() *embeddedFontTheme {
	// 读取嵌入的字体文件（替换为你的字体文件名）
	fontData, err := fontFiles.ReadFile("fonts/NotoSansCJKsc-Regular.otf")
	if err != nil {
		fyne.LogError("无法读取嵌入的字体文件", err)
		return &embeddedFontTheme{baseTheme: theme.DefaultTheme()}
	}

	// 验证字体数据非空
	if len(fontData) == 0 {
		fyne.LogError("嵌入的字体文件为空", nil)
		return &embeddedFontTheme{baseTheme: theme.DefaultTheme()}
	}

	return &embeddedFontTheme{
		baseTheme: theme.DefaultTheme(),
		fontRes:   fyne.NewStaticResource("NotoSansCJKsc-Regular.otf", fontData),
	}
}

func (e *embeddedFontTheme) Font(style fyne.TextStyle) fyne.Resource {
	if e.fontRes != nil {
		logger.Debug("使用嵌入字体")
		return e.fontRes // 忽略样式，所有控件都使用中文字体
	}
	//logger.Debug("使用默认字体")
	return theme.DefaultTextFont() // 失败时回退到默认字体
}

func (e *embeddedFontTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return e.baseTheme.Color(name, variant)
}
func (e *embeddedFontTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return e.baseTheme.Icon(name)
}
func (e *embeddedFontTheme) Size(name fyne.ThemeSizeName) float32 {
	return e.baseTheme.Size(name)
}

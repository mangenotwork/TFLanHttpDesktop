package main

import (
	"embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.variant)
}

// 嵌入字体文件（//go:embed 指令必须放在变量定义前，且路径相对于当前文件）
// 注意：fonts目录下的所有.ttf文件都会被嵌入
//
//go:embed fonts/NotoSans-Regular.ttf
var fontFiles embed.FS

// 自定义主题，使用嵌入的中文字体
type embeddedFontTheme struct {
	baseTheme fyne.Theme
	fontRes   fyne.Resource // 嵌入的字体资源
}

// 初始化主题，加载嵌入的中文字体
func newEmbeddedFontTheme() *embeddedFontTheme {
	// 读取嵌入的字体文件（替换为你的字体文件名）
	fontData, err := fontFiles.ReadFile("fonts/NotoSans-Regular.ttf")
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
		fontRes:   fyne.NewStaticResource("embedded-font", fontData),
	}
}

// 实现Theme接口的Font方法，强制所有控件使用嵌入的字体
func (e *embeddedFontTheme) Font(style fyne.TextStyle) fyne.Resource {
	if e.fontRes != nil {
		return e.fontRes // 忽略样式，所有控件都使用中文字体
	}
	return e.baseTheme.Font(style) // 失败时回退到默认字体
}

// 继承默认主题的其他配置（颜色、图标、尺寸）
func (e *embeddedFontTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return e.baseTheme.Color(name, variant)
}
func (e *embeddedFontTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return e.baseTheme.Icon(name)
}
func (e *embeddedFontTheme) Size(name fyne.ThemeSizeName) float32 {
	return e.baseTheme.Size(name)
}

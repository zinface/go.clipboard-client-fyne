//	Fyne 默认主题缺少中文字体支持，这里嵌入自定义字体资源以便支持中文。
package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/icons/icon.png
var icon_systray []byte
var resourceIconSystray = fyne.NewStaticResource("icon.png", icon_systray)

//go:embed assets/fonts/wqy-microhei.ttc
var font_wqy_microhei []byte
var resourceFontFyne = fyne.NewStaticResource("wqy-microhei.ttc", font_wqy_microhei)

type UnicodeTheme struct{}

func (t *UnicodeTheme) Font(s fyne.TextStyle) fyne.Resource {
	// return theme.DefaultTheme().Font(s)
	return resourceFontFyne
}

func (t *UnicodeTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (t *UnicodeTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *UnicodeTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

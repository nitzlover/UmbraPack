package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type umbraTheme struct{}

func newUmbraTheme() fyne.Theme {
	return &umbraTheme{}
}

func (t *umbraTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 15, G: 15, B: 18, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 35, G: 38, B: 46, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 85, G: 88, B: 96, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 55, G: 58, B: 65, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 230, G: 230, B: 234, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 111, G: 127, B: 255, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 45, G: 48, B: 56, A: 255}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 111, G: 127, B: 255, A: 128}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 65, G: 68, B: 75, A: 200}
	default:
		return theme.DarkTheme().Color(name, variant)
	}
}

func (t *umbraTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (t *umbraTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (t *umbraTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameText:
		return 14
	default:
		return theme.DarkTheme().Size(name)
	}
}

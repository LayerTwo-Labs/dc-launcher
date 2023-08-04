package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed data/images/start.svg
var startIconBytes []byte
var startIconRes fyne.Resource

//go:embed data/images/stop.svg
var stopIconBytes []byte
var stopIconRes fyne.Resource

//go:embed data/images/mine.svg
var mineIconBytes []byte
var mineIconRes fyne.Resource

//go:embed data/images/deposit.svg
var depositIconBytes []byte
var depositIconRes fyne.Resource

//go:embed data/images/withdraw.svg
var withdrawIconBytes []byte
var withdrawIconRes fyne.Resource

//go:embed data/images/search.svg
var searchIconBytes []byte
var searchIconRes fyne.Resource

//go:embed data/images/updown.svg
var upDownIconBytes []byte
var upDownIconRes fyne.Resource

//go:embed data/images/parent.svg
var parentIconBytes []byte
var parentIconRes fyne.Resource

//go:embed data/images/home.svg
var homeIconBytes []byte
var homeIconRes fyne.Resource

//go:embed data/images/calculator.svg
var calculatorIconBytes []byte
var calculatorIconRes fyne.Resource

//go:embed data/images/github.svg
var githubIconBytes []byte
var githubIconRes fyne.Resource

//go:embed data/images/website.svg
var websiteIconBytes []byte
var websiteIconRes fyne.Resource

const (
	StartIcon      fyne.ThemeIconName = "start.svg"
	StopIcon       fyne.ThemeIconName = "stop.svg"
	MineIcon       fyne.ThemeIconName = "mine.svg"
	DepositIcon    fyne.ThemeIconName = "deposit.svg"
	WithdrawIcon   fyne.ThemeIconName = "withdraw.svg"
	SearchIcon     fyne.ThemeIconName = "search.svg"
	UpDownIcon     fyne.ThemeIconName = "updown.svg"
	ParentIcon     fyne.ThemeIconName = "parent.svg"
	HomeIcon       fyne.ThemeIconName = "home.svg"
	CalculatorIcon fyne.ThemeIconName = "calculator.svg"
	GithubIcon     fyne.ThemeIconName = "github.svg"
	WebsiteIcon    fyne.ThemeIconName = "website.svg"
)

var darkScheme = map[fyne.ThemeColorName]color.Color{
	theme.ColorBlue:                  color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @blue_3
	theme.ColorBrown:                 color.RGBA{0x98, 0x6a, 0x44, 0xff}, // Adwaita color name @brown_3
	theme.ColorGray:                  color.RGBA{0x5e, 0x5c, 0x64, 0xff}, // Adwaita color name @dark_2
	theme.ColorGreen:                 color.RGBA{0x26, 0xa2, 0x69, 0xff}, // Adwaita color name @green_5
	theme.ColorNameBackground:        color.RGBA{0x24, 0x24, 0x24, 0xff}, // Adwaita color name @window_bg_color
	theme.ColorNameButton:            color.RGBA{0x30, 0x30, 0x30, 0xff}, // Adwaita color name @headerbar_bg_color
	theme.ColorNameError:             color.RGBA{0xc0, 0x1c, 0x28, 0xff}, // Adwaita color name @error_bg_color
	theme.ColorNameForeground:        color.RGBA{0xef, 0xef, 0xef, 0xff}, // Adwaita color name @window_fg_color
	theme.ColorNameInputBackground:   color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameMenuBackground:    color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameOverlayBackground: color.RGBA{0x1e, 0x1e, 0x1e, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNamePrimary:           color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameSelection:         color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameShadow:            color.RGBA{0x00, 0x00, 0x00, 0x5b}, // Adwaita color name @shade_color
	theme.ColorNameSuccess:           color.RGBA{0x26, 0xa2, 0x69, 0xff}, // Adwaita color name @success_bg_color
	theme.ColorNameWarning:           color.RGBA{0xcd, 0x93, 0x09, 0xff}, // Adwaita color name @warning_bg_color
	theme.ColorOrange:                color.RGBA{0xff, 0x78, 0x00, 0xff}, // Adwaita color name @orange_3
	theme.ColorPurple:                color.RGBA{0x91, 0x41, 0xac, 0xff}, // Adwaita color name @purple_3
	theme.ColorRed:                   color.RGBA{0xc0, 0x1c, 0x28, 0xff}, // Adwaita color name @red_4
	theme.ColorYellow:                color.RGBA{0xf6, 0xd3, 0x2d, 0xff}, // Adwaita color name @yellow_3
}

var lightScheme = map[fyne.ThemeColorName]color.Color{
	theme.ColorBlue:                  color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @blue_3
	theme.ColorBrown:                 color.RGBA{0x98, 0x6a, 0x44, 0xff}, // Adwaita color name @brown_3
	theme.ColorGray:                  color.RGBA{0x5e, 0x5c, 0x64, 0xff}, // Adwaita color name @dark_2
	theme.ColorGreen:                 color.RGBA{0x2e, 0xc2, 0x7e, 0xff}, // Adwaita color name @green_4
	theme.ColorNameBackground:        color.RGBA{0xe9, 0xec, 0xef, 0xff}, // Adwaita color name @window_bg_color
	theme.ColorNameButton:            color.RGBA{0xeb, 0xeb, 0xeb, 0xff}, // Adwaita color name @headerbar_bg_color
	theme.ColorNameError:             color.RGBA{0xe0, 0x1b, 0x24, 0xff}, // Adwaita color name @error_bg_color
	theme.ColorNameForeground:        color.RGBA{0x3d, 0x3d, 0x3d, 0xff}, // Adwaita color name @window_fg_color
	theme.ColorNameInputBackground:   color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameMenuBackground:    color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNameOverlayBackground: color.RGBA{0xff, 0xff, 0xff, 0xff}, // Adwaita color name @view_bg_color
	theme.ColorNamePrimary:           color.RGBA{0x35, 0x84, 0xe4, 0xff}, // Adwaita color name @accent_bg_color
	theme.ColorNameShadow:            color.RGBA{0x00, 0x00, 0x00, 0x11}, // Adwaita color name @shade_color
	theme.ColorNameSuccess:           color.RGBA{0x2e, 0xc2, 0x7e, 0xff}, // Adwaita color name @success_bg_color
	theme.ColorNameWarning:           color.RGBA{0xe5, 0xa5, 0x0a, 0xff}, // Adwaita color name @warning_bg_color
	theme.ColorOrange:                color.RGBA{0xff, 0x78, 0x00, 0xff}, // Adwaita color name @orange_3
	theme.ColorPurple:                color.RGBA{0x91, 0x41, 0xac, 0xff}, // Adwaita color name @purple_3
	theme.ColorRed:                   color.RGBA{0xe0, 0x1b, 0x24, 0xff}, // Adwaita color name @red_3
	theme.ColorYellow:                color.RGBA{0xf6, 0xd3, 0x2d, 0xff}, // Adwaita color name @yellow_3
}

type CustomTheme struct{}

func NewCustomTheme() *CustomTheme {
	t := CustomTheme{}
	startIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(StartIcon), startIconBytes))
	stopIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(StopIcon), stopIconBytes))
	mineIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(MineIcon), mineIconBytes))
	depositIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(DepositIcon), depositIconBytes))
	withdrawIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(WithdrawIcon), withdrawIconBytes))
	searchIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(SearchIcon), searchIconBytes))
	upDownIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(UpDownIcon), upDownIconBytes))
	parentIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(ParentIcon), parentIconBytes))
	homeIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(HomeIcon), homeIconBytes))
	calculatorIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(CalculatorIcon), calculatorIconBytes))
	githubIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(GithubIcon), githubIconBytes))
	websiteIconRes = theme.NewThemedResource(fyne.NewStaticResource(string(WebsiteIcon), websiteIconBytes))
	return &t
}

func (t CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch variant {
	case theme.VariantLight:
		if c, ok := lightScheme[name]; ok {
			return c
		}
	case theme.VariantDark:
		if c, ok := darkScheme[name]; ok {
			return c
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	switch name {
	case StartIcon:
		return startIconRes
	case StopIcon:
		return stopIconRes
	case MineIcon:
		return mineIconRes
	case DepositIcon:
		return depositIconRes
	case WithdrawIcon:
		return withdrawIconRes
	case SearchIcon:
		return searchIconRes
	case UpDownIcon:
		return upDownIconRes
	case ParentIcon:
		return parentIconRes
	case HomeIcon:
		return homeIconRes
	case CalculatorIcon:
		return calculatorIconRes
	case GithubIcon:
		return githubIconRes
	case WebsiteIcon:
		return websiteIconRes
	default:
		return theme.DefaultTheme().Icon(name)
	}
}

func (t CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

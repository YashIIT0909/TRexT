package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Theme defines color scheme for the application
type Theme struct {
	Name       string
	Background tcell.Color
	Foreground tcell.Color
	Border     tcell.Color
	Title      tcell.Color
	Accent     tcell.Color
	Success    tcell.Color
	Warning    tcell.Color
	Error      tcell.Color
	SelectedBg tcell.Color
	SelectedFg tcell.Color
}

// DefaultTheme returns the default color theme
func DefaultTheme() *Theme {
	return &Theme{
		Name:       "default",
		Background: tcell.ColorDefault,
		Foreground: tcell.ColorWhite,
		Border:     tcell.ColorDarkCyan,
		Title:      tcell.NewRGBColor(0, 255, 255),
		Accent:     tcell.ColorDarkCyan,
		Success:    tcell.ColorGreen,
		Warning:    tcell.ColorYellow,
		Error:      tcell.ColorRed,
		SelectedBg: tcell.ColorDarkCyan,
		SelectedFg: tcell.ColorWhite,
	}
}

// DraculaTheme returns a Dracula-inspired theme
func DraculaTheme() *Theme {
	return &Theme{
		Name:       "dracula",
		Background: tcell.NewRGBColor(40, 42, 54),
		Foreground: tcell.NewRGBColor(248, 248, 242),
		Border:     tcell.NewRGBColor(189, 147, 249),
		Title:      tcell.NewRGBColor(139, 233, 253),
		Accent:     tcell.NewRGBColor(189, 147, 249),
		Success:    tcell.NewRGBColor(80, 250, 123),
		Warning:    tcell.NewRGBColor(255, 184, 108),
		Error:      tcell.NewRGBColor(255, 85, 85),
		SelectedBg: tcell.NewRGBColor(68, 71, 90),
		SelectedFg: tcell.NewRGBColor(248, 248, 242),
	}
}

// ApplyTheme applies a theme to tview styles
func ApplyTheme(theme *Theme) {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    theme.Background,
		ContrastBackgroundColor:     theme.SelectedBg,
		MoreContrastBackgroundColor: theme.Accent,
		BorderColor:                 theme.Border,
		TitleColor:                  theme.Title,
		GraphicsColor:               theme.Border,
		PrimaryTextColor:            theme.Foreground,
		SecondaryTextColor:          tcell.ColorGray,
		TertiaryTextColor:           theme.Accent,
		InverseTextColor:            theme.Background,
		ContrastSecondaryTextColor:  theme.SelectedFg,
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "dracula":
		return DraculaTheme()
	default:
		return DefaultTheme()
	}
}

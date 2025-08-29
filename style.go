package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// background color
	nvimBg        = lipgloss.Color("#1a1b26") // fst background
	nvimBgDark    = lipgloss.Color("#16161e") // darker
	nvimBgFloat   = lipgloss.Color("#1f2335") // floating window
	nvimBgSidebar = lipgloss.Color("#16161e") // Sidebar

	// text color
	nvimFg       = lipgloss.Color("#c0caf5") // fst text
	nvimFgDark   = lipgloss.Color("#a9b1d6") // snd text
	nvimFgDarker = lipgloss.Color("#737aa8") // thd text

	// accent
	nvimBlue    = lipgloss.Color("#7aa2f7")
	nvimCyan    = lipgloss.Color("#7dcfff")
	nvimGreen   = lipgloss.Color("#9ece6a")
	nvimYellow  = lipgloss.Color("#e0af68")
	nvimOrange  = lipgloss.Color("#ff9e64")
	nvimRed     = lipgloss.Color("#f7768e")
	nvimPurple  = lipgloss.Color("#bb9af7")
	nvimMagenta = lipgloss.Color("#ff757f")

	// border color
	nvimBorder      = lipgloss.Color("#3b4261")
	nvimBorderFocus = lipgloss.Color("#7aa2f7")
)

// interface style
var (
	// basic background style
	baseStyle = lipgloss.NewStyle().
			Background(nvimBg).
			Foreground(nvimFg)

	// tabline
	headerBarStyle = lipgloss.NewStyle().
			Background(nvimBgDark).
			Foreground(nvimBlue).
			Bold(true).
			Padding(0, 2)

	// sidebar
	sidebarStyle = lipgloss.NewStyle().
			Background(nvimBgSidebar).
			Foreground(nvimFg).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(nvimBorder).
			Padding(1, 1)

	// principal content
	mainContentStyle = lipgloss.NewStyle().
				Background(nvimBg).
				Foreground(nvimFg).
				Padding(1, 2)

	// statusline nvim like
	statusLineStyle = lipgloss.NewStyle().
			Background(nvimBgDark).
			Foreground(nvimFg).
			Padding(0, 2)

	// select elmt
	selectedItemStyle = lipgloss.NewStyle().
				Background(nvimBgFloat).
				Foreground(nvimYellow).
				Bold(true).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(nvimFgDark).
			Padding(0, 1)

	// notif
	successNotifStyle = lipgloss.NewStyle().
				Background(nvimGreen).
				Foreground(nvimBg).
				Bold(true).
				Padding(0, 2)

	errorNotifStyle = lipgloss.NewStyle().
			Background(nvimRed).
			Foreground(nvimBg).
			Bold(true).
			Padding(0, 2)

	// section title
	titleStyle = lipgloss.NewStyle().
			Foreground(nvimBlue).
			Bold(true).
			Underline(true)

	// stat
	statLabelStyle = lipgloss.NewStyle().
			Foreground(nvimFgDark)

	statValueStyle = lipgloss.NewStyle().
			Foreground(nvimYellow).
			Bold(true)

	// help text
	helpTextStyle = lipgloss.NewStyle().
			Foreground(nvimFgDarker).
			Italic(true)
)

func getEventIconAndColor(eventType string) (string, lipgloss.Color) {
	switch eventType {
	case "PushEvent":
		return "✏", nvimGreen
	case "IssuesEvent":
		return "☒", nvimRed
	case "WatchEvent":
		return "☆", nvimYellow
	case "ForkEvent":
		return "⑂", nvimPurple
	case "CreateEvent":
		return "﹢", nvimBlue
	case "DeleteEvent":
		return "␀", nvimRed
	case "PullRequestEvent":
		return "♺", nvimCyan
	case "ReleaseEvent":
		return "𝌚", nvimGreen
	case "PublicEvent":
		return "℗", nvimBlue
	default:
		return "≝", nvimFgDark
	}
}

package initiative

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type colors struct {
	subdued lipgloss.AdaptiveColor
}

type theme struct {
	colors colors

	container lipgloss.Style

	list list.Styles
	help help.Styles
}

func newTheme() theme {
	colors := colors{
		subdued: lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"},
	}

	container := lipgloss.NewStyle().Margin(0, 1)

	list := list.DefaultStyles()
	list.Title = lipgloss.NewStyle()
	list.TitleBar = lipgloss.NewStyle()
	list.StatusBar = lipgloss.NewStyle().Foreground(colors.subdued).PaddingTop(1)
	list.NoItems = lipgloss.NewStyle().Height(0).Width(0)

	help := help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(colors.subdued),
		ShortDesc:      lipgloss.NewStyle().Foreground(colors.subdued),
		ShortSeparator: lipgloss.NewStyle().Foreground(colors.subdued),
		Ellipsis:       lipgloss.NewStyle().Foreground(colors.subdued),
		FullKey:        lipgloss.NewStyle().Foreground(colors.subdued),
		FullDesc:       lipgloss.NewStyle().Foreground(colors.subdued),
		FullSeparator:  lipgloss.NewStyle().Foreground(colors.subdued),
	}

	return theme{
		colors:    colors,
		container: container,

		list: list,
		help: help,
	}
}

var currentTheme = newTheme()

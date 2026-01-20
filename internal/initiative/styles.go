package initiative

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type colors struct {
	text    lipgloss.AdaptiveColor
	muted   lipgloss.AdaptiveColor
	primary lipgloss.AdaptiveColor
	accent  lipgloss.Color
	success lipgloss.AdaptiveColor
	error   lipgloss.AdaptiveColor
	light   lipgloss.AdaptiveColor
}

type theme struct {
	colors colors

	container lipgloss.Style

	list list.Styles
	help help.Styles
	form *huh.Theme
}

func newTheme() theme {
	// Colors inspired by the Charm theme
	colors := colors{
		text:    lipgloss.AdaptiveColor{Light: "235", Dark: "252"},
		muted:   lipgloss.AdaptiveColor{Light: "248", Dark: "238"},
		primary: lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"},
		accent:  lipgloss.Color("#F780E2"),
		success: lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"},
		error:   lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"},
		light:   lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"},
	}

	container := lipgloss.NewStyle().Margin(0, 1)

	listStyles := list.DefaultStyles()
	listStyles.Title = lipgloss.NewStyle()
	listStyles.TitleBar = lipgloss.NewStyle()
	listStyles.StatusBar = lipgloss.NewStyle().Foreground(colors.muted).PaddingTop(1)
	listStyles.NoItems = lipgloss.NewStyle().Height(0).Width(0)

	helpStyles := help.Styles{
		ShortKey:       lipgloss.NewStyle().Foreground(colors.muted),
		ShortDesc:      lipgloss.NewStyle().Foreground(colors.muted),
		ShortSeparator: lipgloss.NewStyle().Foreground(colors.muted),
		Ellipsis:       lipgloss.NewStyle().Foreground(colors.muted),
		FullKey:        lipgloss.NewStyle().Foreground(colors.muted),
		FullDesc:       lipgloss.NewStyle().Foreground(colors.muted),
		FullSeparator:  lipgloss.NewStyle().Foreground(colors.muted),
	}

	formTheme := newFormTheme(colors)

	return theme{
		colors:    colors,
		container: container,

		list: listStyles,
		help: helpStyles,
		form: formTheme,
	}
}

func newFormTheme(colors colors) *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(colors.primary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(colors.primary).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(colors.error)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(colors.error)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(colors.accent)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(colors.accent)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(colors.accent)
	t.Focused.Option = t.Focused.Option.Foreground(colors.text)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(colors.accent)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(colors.success)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(colors.text)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(colors.light).Background(colors.accent)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(colors.text).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(colors.success)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(colors.muted)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(colors.accent)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}

var currentTheme = newTheme()

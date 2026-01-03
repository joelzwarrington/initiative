package components

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EmptyState struct {
	help help.Model

	keyMap  help.KeyMap
	message string

	height int
	width  int
}

func NewEmptyState(message string, keyMap help.KeyMap, width int, height int) *EmptyState {
	return &EmptyState{
		help: help.New(),

		keyMap:  keyMap,
		message: message,

		width:  width,
		height: height,
	}
}

func (e *EmptyState) Init() tea.Cmd {
	return nil
}

func (e *EmptyState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e, nil
}

func (e *EmptyState) SetSize(width int, height int) {
	e.width = width
	e.height = height

	e.help.Width = width
}

func (e *EmptyState) SetShowFullHelp(v bool) {
	e.help.ShowAll = v
}

func (e *EmptyState) View() string {
	// Create message with top-left positioning and padding
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(1, 2, 0, 2) // 1 top, 2 left/right, 0 bottom

	message := messageStyle.Render(e.message)

	// Create content area that takes full space with message at top-left - no help component
	contentStyle := lipgloss.NewStyle().
		Height(e.height).
		Width(e.width).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top)

	return contentStyle.Render(message)
}

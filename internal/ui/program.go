package ui

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
)

func NewProgram(d *data.Data) *tea.Program {
	g := newGame(d)
	return tea.NewProgram(g)
}

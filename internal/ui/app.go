package ui

import (
	"initiative/internal/data"
	"initiative/internal/ui/views"

	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = (*app)(nil)

type app struct {
	views.GameModel
}

func newApp(data *data.Data) app {
	app := app{
		GameModel: views.NewGameModel(data),
	}

	return app
}

func NewProgram(appData *data.Data) *tea.Program {
	app := newApp(appData)
	return tea.NewProgram(app)
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// todo: add global ctrl c handler

	return m.GameModel.Update(msg)
}

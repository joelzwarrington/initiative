package ui

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

var _ tea.Model = (*GameFormModel)(nil)

type GameFormModel struct {
	game *data.Game

	form *huh.Form

	name *string

	width  int
	height int
}

func newGameForm(game *data.Game) *GameFormModel {
	m := GameFormModel{}
	m.SetGame(game)

	return &m
}

func (m GameFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m GameFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		// m.state = stateDone
	}

	return m, cmd
}

func (m GameFormModel) View() string {
	return m.form.View()
}

func (m GameFormModel) SetGame(g *data.Game) {
	m.game = g

	name := ""

	if game := m.game; game != nil {
		name = game.Name
	}

	m.name = &name

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Placeholder("Lorem Ipsum").
				Value(m.name),
		),
	)
}

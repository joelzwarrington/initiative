package ui

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type gameSavedMsg struct {
	uuid string
	name string
}

var _ tea.Model = (*GameFormModel)(nil)

type GameFormModel struct {
	uuid string
	game *data.Game

	form *huh.Form

	name *string

	width  int
	height int
}

func newGameForm(uuid string, game *data.Game) *GameFormModel {
	m := GameFormModel{}
	m.SetGame(uuid, game)

	return &m
}

func (m *GameFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *GameFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m, tea.Cmd(func() tea.Msg {
			return gameSavedMsg{
				uuid: m.uuid,
				name: *m.name,
			}
		})
	}

	return m, cmd
}

func (m GameFormModel) View() string {
	if m.form.State == huh.StateCompleted {
		return "Submitted!"
	}

	return m.form.View()
}

func (m *GameFormModel) SetGame(uuid string, g *data.Game) {
	m.uuid = uuid
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
				Value(m.name),
		),
	)
}

package ui

import (
	"initiative/internal/data"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type characterSavedMsg struct {
	uuid string
	name string
}

type characterDeselectedMsg struct{}

var _ tea.Model = (*CharacterFormModel)(nil)

type CharacterFormModel struct {
	uuid      string
	character *data.Character

	form *huh.Form

	name *string

	width  int
	height int
}

func newCharacterForm(uuid string, character *data.Character) *CharacterFormModel {
	m := CharacterFormModel{}
	m.SetCharacter(uuid, character)

	return &m
}

func (m *CharacterFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CharacterFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle escape key to go back
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, key.NewBinding(key.WithKeys("esc"))) {
			return m, tea.Cmd(func() tea.Msg {
				return characterDeselectedMsg{}
			})
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m, tea.Cmd(func() tea.Msg {
			return characterSavedMsg{
				uuid: m.uuid,
				name: *m.name,
			}
		})
	}

	return m, cmd
}

func (m CharacterFormModel) View() string {
	if m.form.State == huh.StateCompleted {
		return "Submitted!"
	}

	return m.form.View()
}

func (m *CharacterFormModel) SetCharacter(uuid string, c *data.Character) {
	m.uuid = uuid
	m.character = c

	name := ""

	if character := m.character; character != nil {
		name = character.Name
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
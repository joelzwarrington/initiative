package ui

import (
	"initiative/internal/data"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type encounterSavedMsg struct {
	uuid       string
	summary    string
	characters []string
}

type encounterDeselectedMsg struct{}

var _ tea.Model = (*EncounterFormModel)(nil)

type EncounterFormModel struct {
	uuid      string
	encounter *data.Encounter
	game      *data.Game

	form *huh.Form

	summary             *string
	participatingChars  []string

	width  int
	height int
}

func newEncounterForm(game *data.Game) *EncounterFormModel {
	m := EncounterFormModel{
		game: game,
	}
	m.SetEncounter("", nil)

	return &m
}

func (m *EncounterFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *EncounterFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle escape key to go back
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, key.NewBinding(key.WithKeys("esc"))) {
			return m, tea.Cmd(func() tea.Msg {
				return encounterDeselectedMsg{}
			})
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m, tea.Cmd(func() tea.Msg {
			return encounterSavedMsg{
				uuid:       m.uuid,
				summary:    *m.summary,
				characters: m.participatingChars,
			}
		})
	}

	return m, cmd
}

func (m EncounterFormModel) View() string {
	if m.form.State == huh.StateCompleted {
		return "Submitted!"
	}

	return m.form.View()
}

func (m *EncounterFormModel) SetEncounter(uuid string, e *data.Encounter) {
	m.uuid = uuid
	m.encounter = e

	summary := ""
	var participatingChars []string

	if encounter := m.encounter; encounter != nil {
		summary = encounter.Summary
		participatingChars = encounter.Characters
	}

	m.summary = &summary
	m.participatingChars = participatingChars

	// Build character options from game characters
	var characterOptions []huh.Option[string]
	if m.game != nil {
		for _, character := range m.game.Characters {
			characterOptions = append(characterOptions, huh.NewOption(character.Name, character.Name))
		}
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Key("summary").
				Title("Encounter Summary").
				Placeholder("Describe what happens in this encounter...").
				Value(m.summary),

			huh.NewMultiSelect[string]().
				Key("characters").
				Title("Characters").
				Description("Choose which characters participate in this encounter").
				Options(characterOptions...).
				Value(&m.participatingChars),
		),
	)
}
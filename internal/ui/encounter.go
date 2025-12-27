package ui

import (
	"initiative/internal/data"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type encounterView int

const (
	noEncounter encounterView = iota
	encounterFormView
	activeEncounter
)

// Encounter represents an active encounter (ephemeral, not saved)
type Encounter struct {
	Summary    string  
	Characters []string
}

type EncounterModel struct {
	game *data.Game

	currentView encounterView
	encounter   *Encounter
	form        *EncounterFormModel

	help   help.Model
	keyMap encounterKeyMap

	width  int
	height int
}

type encounterKeyMap struct {
	tab      key.Binding
	shiftTab key.Binding
	back     key.Binding
	quit     key.Binding
	help     key.Binding
}

func newEncounterKeyMap() encounterKeyMap {
	return encounterKeyMap{
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		shiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (k encounterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "stop")),
		k.help,
	}
}

func (k encounterKeyMap) FullHelp() [][]key.Binding {
	encounterHelp := []key.Binding{
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new encounter")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "stop encounter")),
	}
	
	gameNavHelp := []key.Binding{
		k.tab, k.shiftTab, k.back, k.quit,
	}
	
	return [][]key.Binding{
		append(encounterHelp, gameNavHelp...),
	}
}

func newEncounter(game *data.Game, width, height int) *EncounterModel {
	return &EncounterModel{
		game:        game,
		currentView: noEncounter,
		help:        help.New(),
		keyMap:      newEncounterKeyMap(),
		width:       width,
		height:      height,
	}
}

func (m *EncounterModel) Init() tea.Cmd {
	return nil
}

func (m *EncounterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *EncounterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case encounterSavedMsg:
		{
			m.encounter = &Encounter{
				Summary:    msg.summary,
				Characters: msg.characters,
			}
			m.currentView = activeEncounter
			return m, nil
		}

	case encounterDeselectedMsg:
		{
			m.currentView = noEncounter
			m.form = nil
			return m, nil
		}

	case tea.KeyMsg:
		// Handle global encounter tab keys first
		switch {
		case key.Matches(msg, m.keyMap.tab):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: 1}
			}
		case key.Matches(msg, m.keyMap.shiftTab):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: -1}
			}
		case key.Matches(msg, m.keyMap.back):
			return m, func() tea.Msg {
				return gameBackMsg{}
			}
		case key.Matches(msg, m.keyMap.quit):
			return m, func() tea.Msg {
				return gameQuitMsg{}
			}
		case key.Matches(msg, m.keyMap.help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			// Start new encounter - only if no active encounter
			if m.currentView == noEncounter && len(m.game.Characters) > 0 {
				m.form = newEncounterForm(m.game)
				m.currentView = encounterFormView
				return m, m.form.Init()
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			// Stop encounter - only if there's an active encounter
			if m.currentView == activeEncounter {
				m.encounter = nil
				m.currentView = noEncounter
				return m, nil
			}
		}
	}

	switch m.currentView {
	case encounterFormView:
		if m.form != nil {
			var cmd tea.Cmd
			f, cmd := m.form.Update(msg)
			if fm, ok := f.(*EncounterFormModel); ok {
				m.form = fm
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m EncounterModel) View() string {
	switch m.currentView {
	case noEncounter:
		if len(m.game.Characters) == 0 {
			return "No characters available.\n\nAdd characters first to start an encounter."
		}
		return "No active encounter.\n\nPress 'n' to start a new encounter."
	case encounterFormView:
		if m.form != nil {
			return m.form.View()
		}
		return "Loading form..."
	case activeEncounter:
		return m.viewActiveEncounter()
	default:
		return "Unknown encounter state"
	}
}

func (m EncounterModel) viewActiveEncounter() string {
	if m.encounter == nil {
		return "No active encounter"
	}
	
	content := "🎲 Active Encounter\n\n"
	
	if m.encounter.Summary != "" {
		content += m.encounter.Summary + "\n\n"
	}
	
	if len(m.encounter.Characters) > 0 {
		content += "Participants:\n"
		for _, charName := range m.encounter.Characters {
			content += "• " + charName + "\n"
		}
		content += "\n"
	}
	
	content += "Press 'esc' to stop the encounter"
	
	return content
}
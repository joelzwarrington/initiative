package views

import (
	"fmt"
	"initiative/internal/data"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InitiativeGroup represents a group of characters with the same initiative
type InitiativeGroup struct {
	Initiative int
	Characters []data.Character
}

// Encounter represents a running encounter
type Encounter struct {
	Name             string
	InitiativeGroups []InitiativeGroup
}


// encounterState represents the current state of the encounter model
type encounterState int

const (
	noEncounter encounterState = iota
	runningEncounter
)

// Encounter keymap
type encounterKeyMap struct {
	startEncounter key.Binding
	stopEncounter  key.Binding
	back           key.Binding
	quit           key.Binding
}


func newEncounterKeyMap() *encounterKeyMap {
	return &encounterKeyMap{
		startEncounter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new encounter"),
		),
		stopEncounter: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop encounter"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k *encounterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.startEncounter, k.stopEncounter}
}

func (k *encounterKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.startEncounter, k.stopEncounter, k.back, k.quit},
	}
}

// EncounterModel manages encounter creation and display
type EncounterModel struct {
	game  *data.Game
	state encounterState
	help  help.Model

	keyMap *encounterKeyMap

	// Current encounter
	encounter *Encounter

	width  int
	height int
}

// NewEncounterModel creates a new encounter model
func NewEncounterModel(game *data.Game) EncounterModel {
	return EncounterModel{
		game:   game,
		state:  noEncounter,
		help:   help.New(),
		keyMap: newEncounterKeyMap(),
	}
}

// Init initializes the encounter model and returns any initial commands
func (m EncounterModel) Init() tea.Cmd {
	return nil
}

// SetSize sets the size of the encounter model
func (m *EncounterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// IsCreating returns true if the encounter model is currently creating an encounter
func (m EncounterModel) IsCreating() bool {
	return false // No creating state for now
}

// Update handles updates for the encounter model
func (m EncounterModel) Update(msg tea.Msg) (EncounterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case noEncounter:
			switch {
			case key.Matches(msg, m.keyMap.startEncounter):
				// TODO: Implement encounter creation later
				return m, nil
			}

		case runningEncounter:
			switch {
			case key.Matches(msg, m.keyMap.stopEncounter):
				m.state = noEncounter
				m.encounter = nil
				return m, nil
			}
		}
	}

	return m, nil
}



// View renders the encounter model
func (m EncounterModel) View() string {
	// Handle case where game is nil (like in tests)
	if m.game == nil {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)
		return emptyStyle.Render("No game loaded.")
	}

	var content string

	switch m.state {
	case noEncounter:
		if len(m.game.Characters) == 0 {
			content = "No characters available.\n\nAdd characters first to start an encounter."
		} else {
			content = "No encounter started yet."
		}

		// Style the content in gray and center it
		grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
		styledContent := grayStyle.Render(content)
		
		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(styledContent)

	case runningEncounter:
		if m.encounter == nil {
			return "No encounter data"
		}

		var s strings.Builder

		// Encounter title
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Margin(0, 0, 1, 0)
		s.WriteString(titleStyle.Render(m.encounter.Name))
		s.WriteString("\n")

		// Initiative order
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Margin(1, 0, 1, 0)
		s.WriteString(headerStyle.Render("Initiative Order:"))
		s.WriteString("\n")

		// List initiative groups
		for i, group := range m.encounter.InitiativeGroups {
			initiativeStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("220"))

			characterStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

			s.WriteString(fmt.Sprintf("%d. ", i+1))
			s.WriteString(initiativeStyle.Render(fmt.Sprintf("Initiative %d: ", group.Initiative)))

			var charNames []string
			for _, char := range group.Characters {
				charNames = append(charNames, char.Name)
			}
			s.WriteString(characterStyle.Render(strings.Join(charNames, ", ")))
			s.WriteString("\n")
		}

		return s.String()
	}

	return "Unknown state"
}

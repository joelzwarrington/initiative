package views

import (
	"fmt"
	"initiative/internal/data"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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

type encounterFormState int

const (
	statusNormal encounterFormState = iota
	stateDone
)

type EncounterFormModel struct {
	game *data.Game

	state encounterFormState
	form  *huh.Form
}

func newEncounterFormModel(game *data.Game) EncounterFormModel {
	m := EncounterFormModel{game: game}

	// Build character options from game characters
	var characterOptions []huh.Option[string]
	if game != nil {
		for _, character := range game.Characters {
			characterOptions = append(characterOptions, huh.NewOption(character.Name, character.Name))
		}
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("summary").
				Title("Encounter Summary").
				Placeholder("Epic battle at the castle gates..."),

			huh.NewMultiSelect[string]().
				Key("characters").
				Title("Characters").
				Description("Choose which characters participate in this encounter").
				Options(characterOptions...),

			huh.NewConfirm().
				Key("done").
				Title("Ready to roll?").
				Affirmative("Yes").
				Negative("Wait, no"),
		),
	)

	return m
}

func (m EncounterFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m EncounterFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fmt.Fprintf(f, "EncounterFormModel.Update: msg=%T %+v, form.State=%d\n", msg, msg, m.form.State)
		f.Close()
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fmt.Fprintf(f, "EncounterFormModel: After update, form.State=%d, returning cmd=%T\n", m.form.State, cmd)
		f.Close()
	}

	if m.form.State == huh.StateCompleted {
		m.state = stateDone
	}

	return m, cmd
}

func (m EncounterFormModel) View() string {
	return m.form.View()
}

// encounterState represents the current state of the encounter model
type encounterState int

const (
	noEncounter encounterState = iota
	creatingEncounter
	formCompleted
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
	encounter     *Encounter
	encounterForm *EncounterFormModel

	// Form completion data
	completedSummary    string
	completedCharacters []string

	width  int
	height int
}

// NewEncounterModel creates a new encounter model
func NewEncounterModel(game *data.Game) EncounterModel {
	formModel := newEncounterFormModel(game)
	return EncounterModel{
		game:          game,
		state:         noEncounter,
		help:          help.New(),
		keyMap:        newEncounterKeyMap(),
		encounterForm: &formModel,
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

// Update handles updates for the encounter model
func (m EncounterModel) Update(msg tea.Msg) (EncounterModel, tea.Cmd) {
	// Debug logging for all messages in encounter model
	if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		fmt.Fprintf(f, "EncounterModel.Update: state=%d, msg=%T %+v\n", m.state, msg, msg)
		f.Close()
	}

	switch m.state {
	case noEncounter:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keyMap.startEncounter):
				if len(m.game.Characters) > 0 {
					m.state = creatingEncounter
					return m, m.encounterForm.Init()
				}
				return m, nil
			default:
				// Don't handle other keys in noEncounter state - let parent handle them
				return m, nil
			}
		}

	case creatingEncounter:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		case tea.KeyMsg:
			if m.encounterForm != nil {
				if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					fmt.Fprintf(f, "EncounterModel: Passing KeyMsg to form: %+v\n", msg)
					f.Close()
				}
				
				var cmd tea.Cmd
				updatedForm, cmd := m.encounterForm.Update(msg)
				if formModel, ok := updatedForm.(EncounterFormModel); ok {
					m.encounterForm = &formModel
				}

				if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					fmt.Fprintf(f, "EncounterModel: Form returned cmd=%T, form.state=%d\n", cmd, m.encounterForm.state)
					f.Close()
				}

				// Check if form is completed
				if m.encounterForm.state == stateDone {
					m.state = formCompleted
					m.encounter = &Encounter{
						Name: "Lorem Isum",
					}
				}

				// In creating state, always consume the message to prevent parent handling
				return m, cmd
			}
		default:
			// Pass non-KeyMsg to form as well
			if m.encounterForm != nil {
				if f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					fmt.Fprintf(f, "EncounterModel: Passing non-KeyMsg to form: %T %+v\n", msg, msg)
					f.Close()
				}
				
				var cmd tea.Cmd
				updatedForm, cmd := m.encounterForm.Update(msg)
				if formModel, ok := updatedForm.(EncounterFormModel); ok {
					m.encounterForm = &formModel
				}

				return m, cmd
			}
		}

	case formCompleted:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		case tea.KeyMsg:
			// Any key press returns to no encounter state
			m.state = noEncounter
			return m, nil
		}

	case runningEncounter:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keyMap.stopEncounter):
				m.state = noEncounter
				m.encounter = nil
				return m, nil
			default:
				// Don't handle other keys in runningEncounter state - let parent handle them
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

	case creatingEncounter:
		if m.encounterForm != nil {
			return m.encounterForm.View()
		}
		return "Form loading..."

	case formCompleted:
		// Display the completion message with instruction
		instructionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
		instruction := instructionStyle.Render("\n\nPress any key to continue...")

		return lipgloss.NewStyle().
			Height(m.height).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Foo" + instruction)

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

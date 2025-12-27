package ui

import (
	"initiative/internal/data"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameTab int

const (
	encounterTab gameTab = iota
	charactersTab
	statsTab
	logTab
)

type gameDataChangedMsg struct{}

var _ tea.Model = (*GameModel)(nil)

type GameModel struct {
	game *data.Game

	help help.Model

	activeTab gameTab
	keyMap    gameKeyMap

	characterModel *CharacterModel
	encounterModel *EncounterModel

	width  int
	height int
}

type gameKeyMap struct {
	tab      key.Binding
	shiftTab key.Binding
	back     key.Binding
	quit     key.Binding
}

func newGameKeyMap() gameKeyMap {
	return gameKeyMap{
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
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k gameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.tab, k.shiftTab, k.back, k.quit}
}

func (k gameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.tab, k.shiftTab, k.back, k.quit},
	}
}


func newGame(game *data.Game, width, height int) *GameModel {
	return &GameModel{
		game:           game,
		help:           help.New(),
		keyMap:         newGameKeyMap(),
		characterModel: newCharacter(game, width, height),
		encounterModel: newEncounter(game, width, height),
		width:          width,
		height:         height,
	}
}

func (m *GameModel) Init() tea.Cmd {
	return nil
}

func (m *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.broadcastWindowResize(msg)
	case tea.KeyMsg:
		// Only handle emergency quit globally
		if key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c"))) {
			return m, tea.Quit
		}
		// Everything else delegated to active tab
		return m.delegateToActiveTab(msg)
	case characterDataChangedMsg:
		// Character data was modified, bubble up to app level to save
		return m, func() tea.Msg {
			return gameDataChangedMsg{}
		}
	case tabNavigationMsg:
		// Handle tab navigation from sub-models
		if msg.direction == 1 {
			m.activeTab = (m.activeTab + 1) % 4
		} else {
			m.activeTab = (m.activeTab - 1 + 4) % 4
		}
		return m, nil
	case gameBackMsg:
		// Handle back navigation from sub-models
		return m, func() tea.Msg {
			return gameDeselectedMsg{}
		}
	case gameQuitMsg:
		// Handle quit from sub-models
		return m, tea.Quit
	default:
		return m.delegateToActiveTab(msg)
	}
}

func (m *GameModel) broadcastWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	
	// Pass window size to all sub-models
	if m.characterModel != nil {
		charModel, _ := m.characterModel.Update(msg)
		if cm, ok := charModel.(*CharacterModel); ok {
			m.characterModel = cm
		}
	}
	if m.encounterModel != nil {
		m.encounterModel.SetSize(msg.Width, msg.Height)
	}
	
	return m, nil
}

func (m *GameModel) delegateToActiveTab(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case encounterTab:
		if m.encounterModel != nil {
			encModel, cmd := m.encounterModel.Update(msg)
			if em, ok := encModel.(*EncounterModel); ok {
				m.encounterModel = em
			}
			return m, cmd
		}
	case charactersTab:
		if m.characterModel != nil {
			charModel, cmd := m.characterModel.Update(msg)
			if cm, ok := charModel.(*CharacterModel); ok {
				m.characterModel = cm
			}
			return m, cmd
		}
	case statsTab:
		return m.handlePlaceholderTab(msg)
	case logTab:
		return m.handlePlaceholderTab(msg)
	}
	return m, nil
}

func (m *GameModel) handlePlaceholderTab(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle basic navigation for unimplemented tabs
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("tab"))):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: 1}
			}
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("shift+tab"))):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: -1}
			}
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("esc"))):
			return m, func() tea.Msg {
				return gameBackMsg{}
			}
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("q"))):
			return m, func() tea.Msg {
				return gameQuitMsg{}
			}
		case key.Matches(keyMsg, key.NewBinding(key.WithKeys("?"))):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}
	return m, nil
}

func (m GameModel) View() string {
	if m.game == nil {
		return "No game selected"
	}

	var (
		sections    []string
		availHeight = m.height
	)

	// Define colors
	highlight := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	// Define borders for tabs
	inactiveTabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	activeTabBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	// Tab styles
	tab := lipgloss.NewStyle().
		Border(inactiveTabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1).
		MarginTop(0)

	activeTab := tab.Border(activeTabBorder, true)

	// Render tabs
	tabs := []string{"Encounter", "Characters", "Stats", "Log"}
	var renderedTabs []string

	for i, t := range tabs {
		var style lipgloss.Style
		if gameTab(i) == m.activeTab {
			style = activeTab
		} else {
			style = tab
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	// Join tabs together
	tabsJoined := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Calculate remaining width and add border line to fill
	tabsWidth := lipgloss.Width(tabsJoined)
	remainingWidth := m.width - tabsWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	// Create the remaining border line
	borderStyle := lipgloss.NewStyle().Foreground(highlight)
	remainingBorder := borderStyle.Render(strings.Repeat("─", remainingWidth))

	// Join tabs with the remaining border
	tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, tabsJoined, remainingBorder)

	// Add tab row to sections and subtract its height
	sections = append(sections, tabRow)
	availHeight -= lipgloss.Height(tabRow)

	// Render help and subtract its height
	var helpView string
	if m.activeTab == charactersTab && m.characterModel != nil {
		// Use character model's help
		helpView = m.characterModel.help.View(m.characterModel.keyMap)
	} else if m.activeTab == encounterTab && m.encounterModel != nil {
		// Use encounter model's help
		helpView = m.encounterModel.help.View(m.encounterModel.keyMap)
	} else {
		// Use game-level help for other tabs
		helpView = m.help.View(m.keyMap)
	}
	availHeight -= lipgloss.Height(helpView)

	// Tab content - placeholder for now
	var contentArea string
	switch m.activeTab {
	case encounterTab:
		if m.encounterModel != nil {
			// Set the size for the encounter model
			m.encounterModel.SetSize(m.width, availHeight)
			contentArea = m.encounterModel.View()
		} else {
			content := "Encounter model not initialized"
			contentArea = lipgloss.NewStyle().
				Height(availHeight).
				Width(m.width).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render(content)
		}
	case charactersTab:
		if m.characterModel != nil {
			// Set the size for the character model
			m.characterModel.SetSize(m.width, availHeight)
			contentArea = m.characterModel.View()
		} else {
			content := "Character model not initialized"
			contentArea = lipgloss.NewStyle().
				Height(availHeight).
				Width(m.width).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render(content)
		}
	case statsTab:
		content := "Stats tab content - placeholder"
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	case logTab:
		content := "Log tab content - placeholder"
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	}

	sections = append(sections, contentArea)
	sections = append(sections, helpView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

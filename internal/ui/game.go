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

var _ tea.Model = (*GameModel)(nil)

type GameModel struct {
	game *data.Game

	help help.Model

	activeTab gameTab
	keyMap    gameKeyMap

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
		game:   game,
		help:   help.New(),
		keyMap: newGameKeyMap(),
		width:  width,
		height: height,
	}
}

func (m *GameModel) Init() tea.Cmd {
	return nil
}

func (m *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.tab):
			m.activeTab = (m.activeTab + 1) % 4
			return m, nil
		case key.Matches(msg, m.keyMap.shiftTab):
			m.activeTab = (m.activeTab - 1 + 4) % 4
			return m, nil
		case key.Matches(msg, m.keyMap.back):
			return m, tea.Cmd(func() tea.Msg {
				return gameDeselectedMsg{}
			})
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
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
	helpView := m.help.View(m.keyMap)
	availHeight -= lipgloss.Height(helpView)

	// Tab content - placeholder for now
	var contentArea string
	switch m.activeTab {
	case encounterTab:
		content := "Encounter tab content - placeholder"
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	case charactersTab:
		content := "Characters tab content - placeholder"
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
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

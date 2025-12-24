package views

import (
	"fmt"
	"initiative/internal/data"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages for delegate communication
type deleteGameMsg struct{ game *data.Game }
type startEditingMsg struct{ index int }
type startNewGameMsg struct{}
type showGameMsg struct{ game *data.Game }

// Command functions
func deleteGameCmd(game *data.Game) tea.Cmd {
	return func() tea.Msg { return deleteGameMsg{game} }
}

func startEditingCmd(index int) tea.Cmd {
	return func() tea.Msg { return startEditingMsg{index} }
}

func startNewGameCmd() tea.Cmd {
	return func() tea.Msg { return startNewGameMsg{} }
}

func showGameCmd(game *data.Game) tea.Cmd {
	return func() tea.Msg { return showGameMsg{game} }
}

type viewState int

const (
	listView viewState = iota
	gameView
)

// Delegate keymap - for item-specific actions
type delegateKeyMap struct {
	choose key.Binding
	edit   key.Binding
	delete key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	}
}

func (k delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.choose, k.edit, k.delete}
}

func (k delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.choose, k.edit, k.delete},
	}
}

// List keymap - for list-wide actions
type listKeyMap struct {
	add    key.Binding
	back   key.Binding
	quit   key.Binding
	save   key.Binding
	cancel key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		save: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "save"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.add}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.add},
	}
}

// Game view keymap - for when viewing a specific game
type gameViewKeyMap struct {
	back     key.Binding
	quit     key.Binding
	tab      key.Binding
	shiftTab key.Binding
}

func newGameViewKeyMap() *gameViewKeyMap {
	return &gameViewKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		shiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
	}
}

func (k gameViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.tab, k.shiftTab, k.back, k.quit}
}

func (k gameViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.tab, k.shiftTab, k.back, k.quit},
	}
}

// Character tab keymap that combines game view and character-specific keys
type characterTabKeyMap struct {
	*gameViewKeyMap
	addCharacter key.Binding
}

func (k *characterTabKeyMap) ShortHelp() []key.Binding {
	return append(k.gameViewKeyMap.ShortHelp(), k.addCharacter)
}

func (k *characterTabKeyMap) FullHelp() [][]key.Binding {
	gameHelp := k.gameViewKeyMap.FullHelp()
	return [][]key.Binding{
		append(gameHelp[0], k.addCharacter),
	}
}

type GameModel struct {
	data *data.Data

	state       viewState
	currentGame *data.Game
	activeTab   int

	list          list.Model
	characterModel CharacterModel
	delegate      *gameDelegate
	help          help.Model

	listKeyMap     *listKeyMap
	delegateKeyMap *delegateKeyMap
	gameViewKeyMap *gameViewKeyMap

	width  int
	height int
}

func NewGameModel(d *data.Data) GameModel {
	listKeyMap := newListKeyMap()
	delegateKeyMap := newDelegateKeyMap()
	gameViewKeyMap := newGameViewKeyMap()
	delegate := newGameDelegate(delegateKeyMap)

	gameList := list.New(
		gamesToListItems(d.Games),
		delegate,
		0,
		0,
	)
	gameList.Title = "Games"
	gameList.SetShowTitle(true)
	gameList.SetShowStatusBar(false)
	gameList.SetFilteringEnabled(false)
	gameList.SetShowHelp(true)

	// Set up additional help keys for the list
	gameList.AdditionalShortHelpKeys = func() []key.Binding {
		// Show both list and delegate keys in short help
		return append(listKeyMap.ShortHelp(), delegateKeyMap.ShortHelp()...)
	}

	gameList.AdditionalFullHelpKeys = func() []key.Binding {
		return append(listKeyMap.FullHelp()[0], delegateKeyMap.FullHelp()[0]...)
	}

	return GameModel{
		data:           d,
		list:           gameList,
		delegate:       delegate,
		state:          listView,
		currentGame:    nil,
		listKeyMap:     listKeyMap,
		delegateKeyMap: delegateKeyMap,
		gameViewKeyMap: gameViewKeyMap,
		help:           help.New(),
	}
}

func (m GameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		if m.currentGame != nil {
			m.characterModel.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	case deleteGameMsg:
		// Handle delete from delegate
		currentIndex := m.list.Index()
		for i := range m.data.Games {
			if &m.data.Games[i] == msg.game {
				m.data.Games = append(m.data.Games[:i], m.data.Games[i+1:]...)
				m.list.SetItems(gamesToListItems(m.data.Games))

				// Set selection to game above (or same index if deleting first item)
				newIndex := currentIndex
				if currentIndex > 0 && currentIndex >= len(m.data.Games) {
					newIndex = len(m.data.Games) - 1
				} else if currentIndex > 0 {
					newIndex = currentIndex - 1
				}

				if len(m.data.Games) > 0 {
					m.list.Select(newIndex)
				}
				break
			}
		}
		m.data.Save()
		return m, nil

	case showGameMsg:
		// Handle show game from delegate
		m.currentGame = msg.game
		m.state = gameView
		// Initialize character model for this game
		m.characterModel = NewCharacterModel(msg.game)
		if m.width > 0 && m.height > 0 {
			m.characterModel.SetSize(m.width, m.height)
		}
		return m, nil

	case startEditingMsg:
		// Handle start editing from delegate
		selectedIndex := msg.index
		if selectedIndex >= 0 && selectedIndex < len(m.data.Games) {
			// Create list items with editing state
			items := make([]list.Item, len(m.data.Games))
			for i := range m.data.Games {
				isEditing := i == selectedIndex
				var textInputPtr *textinput.Model
				if isEditing {
					textInput := newTextInput()
					textInput.SetValue(m.data.Games[selectedIndex].Name)
					textInput.Focus()
					textInputPtr = &textInput
				}
				items[i] = gameListItem{&m.data.Games[i], isEditing, textInputPtr}
			}
			m.list.SetItems(items)
			return m, textinput.Blink
		}

	case startNewGameMsg:
		// Handle new game from delegate
		newGame := data.Game{Name: ""}
		m.data.Games = append(m.data.Games, newGame)
		newIndex := len(m.data.Games) - 1

		// Create list items with editing state
		items := make([]list.Item, len(m.data.Games))
		for i := range m.data.Games {
			isEditing := i == newIndex
			var textInputPtr *textinput.Model
			if isEditing {
				textInput := newTextInput()
				textInput.SetValue("")
				textInput.Focus()
				textInputPtr = &textInput
			}
			items[i] = gameListItem{&m.data.Games[i], isEditing, textInputPtr}
		}
		m.list.SetItems(items)
		m.list.Select(newIndex)
		return m, textinput.Blink

	case tea.KeyMsg:
		switch m.state {
		case listView:
			if m.isEditing() {
				// Get the current editing item
				selectedIndex := m.list.Index()
				if selectedItem := m.list.SelectedItem(); selectedItem != nil {
					if gameItem, ok := selectedItem.(gameListItem); ok && gameItem.isEditing && gameItem.textInput != nil {
						switch {
						case key.Matches(msg, m.listKeyMap.save):
							name := strings.TrimSpace(gameItem.textInput.Value())
							if name != "" {
								// Finish editing - update the game
								if selectedIndex >= 0 && selectedIndex < len(m.data.Games) {
									m.data.Games[selectedIndex].Name = name
								}
								gameItem.textInput.Blur()
								m.list.SetItems(gamesToListItems(m.data.Games))
								m.data.Save()
								return m, nil
							}
							return m, nil
						case key.Matches(msg, m.listKeyMap.cancel):
							// Cancel editing
							if selectedIndex >= 0 {
								// If this was a new empty game, remove it
								if selectedIndex < len(m.data.Games) && m.data.Games[selectedIndex].Name == "" {
									m.data.Games = append(m.data.Games[:selectedIndex], m.data.Games[selectedIndex+1:]...)
									// Update selection to previous item if possible
									if len(m.data.Games) > 0 && selectedIndex > 0 {
										m.list.Select(selectedIndex - 1)
									}
								}
							}
							gameItem.textInput.Blur()
							m.list.SetItems(gamesToListItems(m.data.Games))
							return m, nil
						default:
							// Update the text input and refresh the list items
							var cmd tea.Cmd
							*gameItem.textInput, cmd = gameItem.textInput.Update(msg)

							// Update list items to show the current text input value
							items := make([]list.Item, len(m.data.Games))
							for i := range m.data.Games {
								isEditing := i == selectedIndex
								var textInputPtr *textinput.Model
								if isEditing {
									textInputPtr = gameItem.textInput
								}
								items[i] = gameListItem{&m.data.Games[i], isEditing, textInputPtr}
							}
							m.list.SetItems(items)

							return m, cmd
						}
					}
				}
			} else {
				// List-wide keys that aren't handled by delegate
				switch {
				case key.Matches(msg, m.listKeyMap.add):
					return m, startNewGameCmd()
				case key.Matches(msg, m.listKeyMap.quit):
					return m, tea.Quit
				default:
					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
			}

		case gameView:
			switch {
			case key.Matches(msg, m.gameViewKeyMap.back):
				m.state = listView
				return m, nil
			case key.Matches(msg, m.gameViewKeyMap.quit):
				return m, tea.Quit
			case key.Matches(msg, m.gameViewKeyMap.tab):
				m.activeTab = (m.activeTab + 1) % 4
				return m, nil
			case key.Matches(msg, m.gameViewKeyMap.shiftTab):
				m.activeTab = (m.activeTab - 1 + 4) % 4
				return m, nil
			default:
				// Delegate to character model when on Characters tab
				if m.activeTab == 1 { // Characters tab
					var cmd tea.Cmd
					m.characterModel, cmd = m.characterModel.Update(msg)
					// Save the data after character changes
					m.data.Save()
					return m, cmd
				}
			}
		}
	}

	return m, nil
}

func (m GameModel) View() string {
	switch m.state {
	case listView:
		return m.listView()
	case gameView:
		return m.gameView()
	default:
		return "Unknown view state"
	}
}

func (m GameModel) listView() string {
	// Set help style with bottom padding
	m.list.Styles.HelpStyle = m.list.Styles.HelpStyle.PaddingBottom(1)
	return m.list.View()
}

func (m GameModel) gameView() string {
	if m.currentGame == nil {
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
		Padding(0, 1)

	activeTab := tab.Border(activeTabBorder, true)

	// Render tabs
	tabs := []string{"Encounter", "Characters", "Stats", "Log"}
	var renderedTabs []string

	for i, t := range tabs {
		var style lipgloss.Style
		if i == m.activeTab {
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

	// Help view with bottom padding - combine keys based on active tab
	var helpKeyMap help.KeyMap
	if m.activeTab == 1 { // Characters tab
		// Create combined keymap for character tab
		helpKeyMap = &characterTabKeyMap{
			gameViewKeyMap: m.gameViewKeyMap,
			addCharacter: key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "add character"),
			),
		}
	} else {
		helpKeyMap = m.gameViewKeyMap
	}
	
	helpStyle := lipgloss.NewStyle().PaddingBottom(1)
	helpView := helpStyle.Render(m.help.View(helpKeyMap))
	availHeight -= lipgloss.Height(helpView)

	// Tab content
	var contentArea string
	switch m.activeTab {
	case 0: // Encounter
		content := "No encounter started yet."
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	case 1: // Characters  
		// Use character model for character tab
		// Set the character model size to available space
		m.characterModel.SetSize(m.width, availHeight)
		contentArea = m.characterModel.View()
	case 2: // Stats
		content := "No stats to display."
		contentArea = lipgloss.NewStyle().
			Height(availHeight).
			Width(m.width).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	case 3: // Log
		content := "No log entries yet."
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

// gameListItem wraps data.Game to implement list.Item interface
type gameListItem struct {
	*data.Game
	isEditing bool
	textInput *textinput.Model
}

func (g gameListItem) Title() string       { return g.Name }
func (g gameListItem) FilterValue() string { return g.Name }

func gamesToListItems(games []data.Game) []list.Item {
	items := make([]list.Item, len(games))
	for i := range games {
		items[i] = gameListItem{&games[i], false, nil}
	}
	return items
}

// Helper function to create a new text input
func newTextInput() textinput.Model {
	textInput := textinput.New()
	textInput.Placeholder = "Name..."
	textInput.CharLimit = 50
	textInput.Width = 30
	textInput.Prompt = ""
	textInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	textInput.Cursor.SetMode(cursor.CursorBlink)
	return textInput
}

// Helper method to check if we're currently editing
func (m GameModel) isEditing() bool {
	if selectedItem := m.list.SelectedItem(); selectedItem != nil {
		if gameItem, ok := selectedItem.(gameListItem); ok {
			return gameItem.isEditing
		}
	}
	return false
}

// Helper method to get the editing item
func (m GameModel) getEditingItem() *gameListItem {
	if selectedItem := m.list.SelectedItem(); selectedItem != nil {
		if gameItem, ok := selectedItem.(gameListItem); ok && gameItem.isEditing {
			return &gameItem
		}
	}
	return nil
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type gameDelegate struct {
	keyMap *delegateKeyMap
}

func newGameDelegate(keyMap *delegateKeyMap) *gameDelegate {
	return &gameDelegate{keyMap: keyMap}
}

func (d *gameDelegate) Height() int  { return 1 }
func (d *gameDelegate) Spacing() int { return 0 }
func (d *gameDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, d.keyMap.choose):
			if selectedItem := m.SelectedItem(); selectedItem != nil {
				if gameItem, ok := selectedItem.(gameListItem); ok {
					return showGameCmd(gameItem.Game)
				}
			}
		case key.Matches(keyMsg, d.keyMap.delete):
			if selectedItem := m.SelectedItem(); selectedItem != nil {
				if gameItem, ok := selectedItem.(gameListItem); ok {
					return deleteGameCmd(gameItem.Game)
				}
			}
		case key.Matches(keyMsg, d.keyMap.edit):
			if selectedItem := m.SelectedItem(); selectedItem != nil {
				if _, ok := selectedItem.(gameListItem); ok {
					return startEditingCmd(m.Index())
				}
			}
		}
	}
	return nil
}
func (d *gameDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	g, ok := listItem.(gameListItem)
	if !ok {
		return
	}

	var str string
	if g.isEditing && g.textInput != nil {
		// When editing, use the text input from the item
		str = fmt.Sprintf("%d. %s", index+1, g.textInput.View())
	} else {
		str = fmt.Sprintf("%d. %s", index+1, g.Name)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

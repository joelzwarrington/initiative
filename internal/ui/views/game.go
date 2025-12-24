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
	back key.Binding
	quit key.Binding
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
	}
}

func (k gameViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.back, k.quit}
}

func (k gameViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.back, k.quit},
	}
}

type GameModel struct {
	data *data.Data

	state       viewState
	currentGame *data.Game

	list     list.Model
	delegate *gameDelegate
	help     help.Model

	listKeyMap     *listKeyMap
	delegateKeyMap *delegateKeyMap
	gameViewKeyMap *gameViewKeyMap
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
		m.list.SetSize(msg.Width, msg.Height)
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
	return m.list.View()
}

func (m GameModel) gameView() string {
	var s strings.Builder

	if m.currentGame == nil {
		return "No game selected"
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	title := titleStyle.Render(m.currentGame.Name)
	s.WriteString(title)
	s.WriteString("\n\n")

	helpView := m.help.View(m.gameViewKeyMap)
	s.WriteString(helpView)

	return s.String()
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

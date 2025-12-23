package views

import (
	"fmt"
	"initiative/internal/data"
	"initiative/internal/ui/messages"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &GamePageModel{}
var _ tea.Model = &GameListModel{}


// gameListItem wraps data.Game to implement list.Item interface
type gameListItem struct {
	*data.Game
	isNew     bool
	isEditing bool
	textInput *textinput.Model
}

func (g gameListItem) Title() string       { return g.Name }
func (g gameListItem) FilterValue() string { return g.Name }

func gamesToListItems(games []data.Game) []list.Item {
	items := make([]list.Item, len(games))
	for i := range games {
		items[i] = gameListItem{&games[i], false, false, nil}
	}
	return items
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type gameDelegate struct{}

func newGameDelegate() *gameDelegate {
	return &gameDelegate{}
}

func (d *gameDelegate) Height() int                             { return 1 }
func (d *gameDelegate) Spacing() int                            { return 0 }
func (d *gameDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d *gameDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	g, ok := listItem.(gameListItem)
	if !ok {
		return
	}

	var str string
	if g.isEditing && g.textInput != nil {
		// When editing, show the text input with number prefix and proper selection indicator
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

type GameListModel struct {
	list         list.Model
	games        *[]data.Game
	currentGame  **data.Game
	delegate     *gameDelegate
	editing      bool
	editingIndex int
	textInput    textinput.Model
}

func NewGameListModel(games *[]data.Game, currentGame **data.Game) *GameListModel {
	textInput := textinput.New()
	textInput.Placeholder = "Enter game name..."
	textInput.CharLimit = 50
	textInput.Width = 30
	textInput.Prompt = ""  // Remove the default "> " prompt
	textInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	gameList := &GameListModel{
		games:        games,
		currentGame:  currentGame,
		editing:      false,
		editingIndex: -1,
		textInput:    textInput,
	}

	gameList.delegate = newGameDelegate()

	gameList.list = list.New(
		gamesToListItems(*games),
		gameList.delegate,
		0,
		0,
	)

	gameList.list.Title = "Games"
	gameList.list.SetShowTitle(true)
	gameList.list.SetShowStatusBar(false)
	gameList.list.SetFilteringEnabled(false)
	gameList.list.SetShowHelp(true)

	gameList.list.AdditionalShortHelpKeys = func() []key.Binding {
		if gameList.editing {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
			}
		}
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	gameList.list.AdditionalFullHelpKeys = func() []key.Binding {
		if gameList.editing {
			return []key.Binding{
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
				key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
			}
		}
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		}
	}

	return gameList
}

func (m *GameListModel) setSize(width, height int) {
	m.list.SetSize(width, height)
}


func (m *GameListModel) RefreshItems() {
	m.list.SetItems(gamesToListItems(*m.games))
}

func (m *GameListModel) selectedGame() *data.Game {
	if selectedItem := m.list.SelectedItem(); selectedItem != nil {
		if gameItem, ok := selectedItem.(gameListItem); ok {
			return gameItem.Game
		}
	}
	return nil
}

func (m *GameListModel) deleteSelectedGame() {
	if selectedItem := m.list.SelectedItem(); selectedItem != nil {
		if gameItem, ok := selectedItem.(gameListItem); ok {
			currentIndex := m.list.Index()
			// Find and remove the game from the slice
			for i := range *m.games {
				if &(*m.games)[i] == gameItem.Game {
					*m.games = append((*m.games)[:i], (*m.games)[i+1:]...)
					
					// Update list items
					m.list.SetItems(gamesToListItems(*m.games))
					
					// Set selection to game above (or same index if deleting first item)
					newIndex := currentIndex
					if currentIndex > 0 && currentIndex >= len(*m.games) {
						newIndex = len(*m.games) - 1
					} else if currentIndex > 0 {
						newIndex = currentIndex - 1
					}
					
					// Only set cursor if there are still games
					if len(*m.games) > 0 {
						m.list.Select(newIndex)
					}
					break
				}
			}
		}
	}
}

func (m *GameListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *GameListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		if m.editing {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				name := strings.TrimSpace(m.textInput.Value())
				if name != "" {
					m.finishEditing(name)
					return m, messages.SaveData()
				}
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.cancelEditing()
				return m, nil
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		} else {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
				m.startNewGame()
				return m, textinput.Blink
			case key.Matches(msg, key.NewBinding(key.WithKeys("e"))):
				if len(*m.games) > 0 {
					m.startEditingSelected()
					return m, textinput.Blink
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
				return m, tea.Quit
			case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
				if len(*m.games) > 0 {
					m.deleteSelectedGame()
					return m, messages.SaveData()
				}
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				if selectedGame := m.selectedGame(); selectedGame != nil {
					return m, messages.NavigateToShowGame(selectedGame)
				}
			default:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
		}
	}
	
	// Only update list if not editing
	if !m.editing {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

func (m *GameListModel) startNewGame() {
	newGame := data.Game{Name: ""}
	*m.games = append(*m.games, newGame)
	
	newIndex := len(*m.games) - 1
	
	m.editing = true
	m.editingIndex = newIndex
	m.textInput.SetValue("")
	m.textInput.Focus()
	
	// Create list items with editing state
	items := make([]list.Item, len(*m.games))
	for i := range *m.games {
		isEditing := i == newIndex
		var textInputPtr *textinput.Model
		if isEditing {
			textInputPtr = &m.textInput
		}
		items[i] = gameListItem{&(*m.games)[i], i == newIndex, isEditing, textInputPtr}
	}
	
	m.list.SetItems(items)
	m.list.Select(newIndex)
}

func (m *GameListModel) startEditingSelected() {
	selectedIndex := m.list.Index()
	if selectedIndex < 0 || selectedIndex >= len(*m.games) {
		return
	}
	
	m.editing = true
	m.editingIndex = selectedIndex
	m.textInput.SetValue((*m.games)[selectedIndex].Name)
	m.textInput.Focus()
	
	// Create list items with editing state
	items := make([]list.Item, len(*m.games))
	for i := range *m.games {
		isEditing := i == selectedIndex
		var textInputPtr *textinput.Model
		if isEditing {
			textInputPtr = &m.textInput
		}
		items[i] = gameListItem{&(*m.games)[i], false, isEditing, textInputPtr}
	}
	
	m.list.SetItems(items)
}

func (m *GameListModel) finishEditing(name string) {
	if m.editingIndex >= 0 && m.editingIndex < len(*m.games) {
		(*m.games)[m.editingIndex].Name = name
	}
	
	m.editing = false
	m.editingIndex = -1
	m.textInput.Blur()
	
	// Refresh items without editing state
	m.RefreshItems()
}

func (m *GameListModel) cancelEditing() {
	if m.editingIndex >= 0 {
		items := m.list.Items()
		if m.editingIndex < len(items) {
			if gameItem, ok := items[m.editingIndex].(gameListItem); ok && gameItem.isNew {
				// Remove the new game that was being created
				*m.games = (*m.games)[:len(*m.games)-1]
				
				// Update selection to previous item if possible
				if len(*m.games) > 0 && m.editingIndex > 0 {
					m.list.Select(m.editingIndex - 1)
				}
			}
		}
	}
	
	m.editing = false
	m.editingIndex = -1
	m.textInput.Blur()
	
	// Refresh items without editing state
	m.RefreshItems()
}

func (m *GameListModel) View() string {
	return m.list.View()
}

type GamePageModel struct {
	currentGame *data.Game
}

func NewGamePageModel() *GamePageModel {
	return &GamePageModel{}
}

func (m *GamePageModel) SetCurrentGame(game *data.Game) {
	m.currentGame = game
}

func (m *GamePageModel) Init() tea.Cmd {
	return nil
}

func (m *GamePageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			return m, messages.NavigateToGameList()
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *GamePageModel) View() string {
	var s strings.Builder

	if m.currentGame == nil {
		return "No game selected"
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	title := titleStyle.Render(m.currentGame.Name)
	s.WriteString(title)
	s.WriteString("\n\n")

	helpModel := help.New()
	helpView := helpModel.View(gamePageKeyMap{})
	s.WriteString(helpView)

	return s.String()
}

type gamePageKeyMap struct{}

func (k gamePageKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

func (k gamePageKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		},
	}
}

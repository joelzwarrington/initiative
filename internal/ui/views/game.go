package views

import (
	"fmt"
	"initiative/internal/data"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &GamePageModel{}
var _ tea.Model = &GameNewFormModel{}
var _ tea.Model = &GameListModel{}

// Navigation commands that views can return to communicate with the app
type NavigateToGameListMsg struct{}
type NavigateToNewGameFormMsg struct{}
type NavigateToShowGameMsg struct {
	Game *data.Game
}
type SaveDataMsg struct{}

func NavigateToGameList() tea.Cmd {
	return func() tea.Msg { return NavigateToGameListMsg{} }
}

func NavigateToNewGameForm() tea.Cmd {
	return func() tea.Msg { return NavigateToNewGameFormMsg{} }
}

func NavigateToShowGame(game *data.Game) tea.Cmd {
	return func() tea.Msg { return NavigateToShowGameMsg{Game: game} }
}

func SaveData() tea.Cmd {
	return func() tea.Msg { return SaveDataMsg{} }
}

// gameListItem wraps data.Game to implement list.Item interface
type gameListItem struct {
	*data.Game
}

func (g gameListItem) Title() string       { return g.Name }
func (g gameListItem) FilterValue() string { return g.Name }

func gamesToListItems(games []data.Game) []list.Item {
	items := make([]list.Item, len(games))
	for i, game := range games {
		items[i] = gameListItem{&game}
	}
	return items
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type gameDelegate struct{}

func (d gameDelegate) Height() int                             { return 1 }
func (d gameDelegate) Spacing() int                            { return 0 }
func (d gameDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gameDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	g, ok := listItem.(gameListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, g.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type GameListModel struct {
	list        list.Model
	games       *[]data.Game
	currentGame **data.Game
}

func NewGameListModel(games *[]data.Game, currentGame **data.Game) *GameListModel {
	gameList := &GameListModel{
		games:       games,
		currentGame: currentGame,
	}

	gameList.list = list.New(
		gamesToListItems(*games),
		gameDelegate{},
		0,
		0,
	)

	gameList.list.Title = "Games"
	gameList.list.SetShowTitle(true)
	gameList.list.SetShowStatusBar(false)
	gameList.list.SetFilteringEnabled(false)
	gameList.list.SetShowHelp(true)

	gameList.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		}
	}

	gameList.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
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

func (m *GameListModel) Init() tea.Cmd {
	return nil
}

func (m *GameListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			return m, NavigateToNewGameForm()
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if selectedGame := m.selectedGame(); selectedGame != nil {
				return m, NavigateToShowGame(selectedGame)
			}
		}
	}
	
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *GameListModel) View() string {
	return m.list.View()
}

type GameNewFormModel struct {
	form  huh.Form
	games *[]data.Game
}

func NewGameNewFormModel(games *[]data.Game) *GameNewFormModel {
	gameForm := &GameNewFormModel{
		games: games,
	}

	gameForm.form = *huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Validate(func(str string) error {
					if str == "" {
						return nil
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return gameForm
}

func (m *GameNewFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *GameNewFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.reset()
			return m, tea.Batch(m.Init(), NavigateToGameList())
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = *f
	}

	if m.isCompleted() {
		m.addGame()
		m.reset()
		return m, tea.Batch(m.Init(), SaveData(), NavigateToGameList())
	}

	return m, cmd
}

func (m *GameNewFormModel) View() string {
	return m.form.View()
}

func (m *GameNewFormModel) isCompleted() bool {
	return m.form.State == huh.StateCompleted
}

func (m *GameNewFormModel) getGameName() string {
	return m.form.GetString("name")
}

func (m *GameNewFormModel) addGame() {
	name := m.getGameName()
	newGame := data.Game{Name: name}
	*m.games = append(*m.games, newGame)
}

func (m *GameNewFormModel) reset() {
	m.form = *huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Validate(func(str string) error {
					if str == "" {
						return nil
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())
	m.form.Init()
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
			return m, NavigateToGameList()
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

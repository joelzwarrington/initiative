package ui

import (
	"fmt"
	"initiative/internal/data"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = (*GameListModel)(nil)

type GameListModel struct {
	list list.Model

	width  int
	height int
}

type viewGameMsg struct {
	uuid string
}

type editGameMsg struct {
	uuid string
}

type deleteGameMsg struct {
	index int
	uuid  string
}

func newGameList(games map[string]data.Game) *GameListModel {
	items := []list.Item{}
	for uuid, game := range games {
		items = append(items, gameItem{uuid: uuid, Game: game})
	}

	keys := newGameItemKeyMap()
	// Disable item actions if no games
	if len(items) == 0 {
		keys.view.SetEnabled(false)
		keys.edit.SetEnabled(false)
		keys.delete.SetEnabled(false)
	}
	delegate := &gameItemDelegate{keys: keys}
	l := list.New(items, delegate, 0, 0)
	l.Title = "Games"
	l.SetStatusBarItemName("game", "games")
	l.KeyMap = newGameListKeyMap()
	
	additionalKeys := newAdditionalGameListKeyMap()
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{additionalKeys.newGame}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{additionalKeys.newGame}
	}

	return &GameListModel{list: l}
}

// newGameListKeyMap returns a modified default set of keybindings.
func newGameListKeyMap() list.KeyMap {
	keyMap := list.DefaultKeyMap()

	// Disable GoToStart and GoToEnd
	keyMap.GoToStart = key.NewBinding(key.WithDisabled())
	keyMap.GoToEnd = key.NewBinding(key.WithDisabled())

	// Update cursor keys to show only arrows in help
	keyMap.CursorUp = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	)
	keyMap.CursorDown = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	)

	return keyMap
}

func (m *GameListModel) Init() tea.Cmd {
	return nil
}

func (m *GameListModel) RemoveGame(uuid string) {
	items := m.list.Items()
	var index int

	for i, item := range items {
		if gameItem, ok := item.(gameItem); ok && gameItem.uuid == uuid {
			m.list.RemoveItem(i)
			index = i
			break
		}
	}

	if index == m.list.Index() {
		m.list.Select(max(0, index-1))
	}
}

func (m *GameListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)

		return m, nil
	case tea.KeyMsg:
		additionalKeys := newAdditionalGameListKeyMap()
		if key.Matches(msg, additionalKeys.newGame) {
			return m, tea.Cmd(func() tea.Msg {
				return editGameMsg{uuid: ""}
			})
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m GameListModel) View() string {
	return m.list.View()
}

type gameItem struct {
	uuid string
	data.Game
}

func (g gameItem) FilterValue() string { return g.Name }

type gameItemDelegate struct {
	keys gameItemKeyMap
}

func (d gameItemDelegate) Height() int  { return 1 }
func (d gameItemDelegate) Spacing() int { return 0 }
func (d *gameItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// Update key states based on item count
	hasItems := len(m.Items()) > 0
	d.keys.view.SetEnabled(hasItems)
	d.keys.edit.SetEnabled(hasItems)
	d.keys.delete.SetEnabled(hasItems)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.view):
			if item, ok := m.SelectedItem().(gameItem); ok {
				return tea.Cmd(func() tea.Msg {
					return viewGameMsg{uuid: item.uuid}
				})
			}

		case key.Matches(msg, d.keys.edit):
			if item, ok := m.SelectedItem().(gameItem); ok {
				return tea.Cmd(func() tea.Msg {
					return editGameMsg{uuid: item.uuid}
				})
			}

		case key.Matches(msg, d.keys.delete):
			if item, ok := m.SelectedItem().(gameItem); ok {
				return tea.Cmd(func() tea.Msg {
					return deleteGameMsg{uuid: item.uuid}
				})
			}
		}
	}

	return nil
}

func (d gameItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(gameItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (d gameItemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		d.keys.view,
		d.keys.edit,
		d.keys.delete,
	}
}

func (d gameItemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.keys.view,
			d.keys.edit,
			d.keys.delete,
		},
	}
}

type gameItemKeyMap struct {
	view   key.Binding
	edit   key.Binding
	delete key.Binding
}

func (d gameItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.view,
		d.edit,
		d.delete,
	}
}

func (d gameItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.view,
			d.edit,
			d.delete,
		},
	}
}

func newGameItemKeyMap() gameItemKeyMap {
	return gameItemKeyMap{
		view: key.NewBinding(
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

type additionalGameListKeyMap struct {
	newGame key.Binding
}

func newAdditionalGameListKeyMap() additionalGameListKeyMap {
	return additionalGameListKeyMap{
		newGame: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
	}
}

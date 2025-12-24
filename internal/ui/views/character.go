package views

import (
	"initiative/internal/data"
	"io"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Character delegate keymap
type characterDelegateKeyMap struct {
	choose key.Binding
	delete key.Binding
}

func newCharacterDelegateKeyMap() *characterDelegateKeyMap {
	return &characterDelegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	}
}

// Character list keymap
type characterListKeyMap struct {
	add  key.Binding
	back key.Binding
	quit key.Binding
}

func newCharacterListKeyMap() *characterListKeyMap {
	return &characterListKeyMap{
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
	}
}

// CharacterModel represents the character list component
type CharacterModel struct {
	game *data.Game
	list list.Model
	help help.Model

	listKeyMap     *characterListKeyMap
	delegateKeyMap *characterDelegateKeyMap

	width  int
	height int
}

// NewCharacterModel creates a new character model
func NewCharacterModel(game *data.Game) CharacterModel {
	listKeyMap := newCharacterListKeyMap()
	delegateKeyMap := newCharacterDelegateKeyMap()
	delegate := newCharacterDelegate(delegateKeyMap)

	characterList := list.New(
		charactersToListItems(game.Characters),
		delegate,
		0,
		0,
	)
	characterList.Title = "Characters"
	characterList.SetShowTitle(false)
	characterList.SetShowStatusBar(false)
	characterList.SetFilteringEnabled(false)
	characterList.SetShowHelp(false) // We'll handle help ourselves

	return CharacterModel{
		game:           game,
		list:           characterList,
		help:           help.New(),
		listKeyMap:     listKeyMap,
		delegateKeyMap: delegateKeyMap,
	}
}

// SetSize sets the size of the character model
func (m *CharacterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// Update handles updates for the character model
func (m CharacterModel) Update(msg tea.Msg) (CharacterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.listKeyMap.add):
			// Add new character
			newChar := data.Character{Name: "New Character"}
			m.game.Characters = append(m.game.Characters, newChar)
			m.list.SetItems(charactersToListItems(m.game.Characters))
			// Make sure the list has proper size after adding items
			if m.width > 0 && m.height > 0 {
				m.list.SetSize(m.width, m.height)
			}
			// Select the new character
			m.list.Select(len(m.game.Characters) - 1)
			return m, nil
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the character model
func (m CharacterModel) View() string {
	if len(m.game.Characters) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true).
			Padding(2, 4)
		return emptyStyle.Render("No characters added yet.\n\nPress 'a' to add a character.")
	}

	return m.list.View()
}

// characterListItem wraps data.Character to implement list.Item interface
type characterListItem struct {
	*data.Character
}

func (c characterListItem) Title() string       { return c.Name }
func (c characterListItem) FilterValue() string { return c.Name }

func charactersToListItems(characters []data.Character) []list.Item {
	items := make([]list.Item, len(characters))
	for i := range characters {
		items[i] = characterListItem{&characters[i]}
	}
	return items
}

// Character delegate
type characterDelegate struct {
	keyMap *characterDelegateKeyMap
}

func newCharacterDelegate(keyMap *characterDelegateKeyMap) *characterDelegate {
	return &characterDelegate{keyMap: keyMap}
}

func (d *characterDelegate) Height() int  { return 1 }
func (d *characterDelegate) Spacing() int { return 0 }

func (d *characterDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d *characterDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	c, ok := listItem.(characterListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, c.Name)

	itemStyle := lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle := lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
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

var _ tea.Model = (*CharacterListModel)(nil)

type CharacterListModel struct {
	list list.Model

	width  int
	height int
}

type characterEditedMsg struct {
	uuid string
}

type characterDeletedMsg struct {
	uuid string
}

func newCharacterList(characters map[string]data.Character) *CharacterListModel {
	items := []list.Item{}
	for uuid, character := range characters {
		items = append(items, characterItem{uuid: uuid, Character: character})
	}

	keys := newCharacterItemKeyMap()
	// Disable item actions if no characters
	if len(items) == 0 {
		keys.view.SetEnabled(false)
		keys.edit.SetEnabled(false)
		keys.delete.SetEnabled(false)
	}
	delegate := &characterItemDelegate{keys: keys}
	l := list.New(items, delegate, 0, 0)
	l.Title = "Characters"
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetStatusBarItemName("character", "characters")
	l.KeyMap = newCharacterListKeyMap()
	l.SetShowHelp(false) // Hide the default help
	
	additionalKeys := newAdditionalCharacterListKeyMap()
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{additionalKeys.newCharacter}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{additionalKeys.newCharacter}
	}

	return &CharacterListModel{list: l}
}

// newCharacterListKeyMap returns a modified default set of keybindings.
func newCharacterListKeyMap() list.KeyMap {
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

func (m *CharacterListModel) Init() tea.Cmd {
	return nil
}

func (m *CharacterListModel) RemoveCharacter(uuid string) {
	items := m.list.Items()
	var index int

	for i, item := range items {
		if characterItem, ok := item.(characterItem); ok && characterItem.uuid == uuid {
			m.list.RemoveItem(i)
			index = i
			break
		}
	}

	if index == m.list.Index() {
		m.list.Select(max(0, index-1))
	}
}

func (m *CharacterListModel) UpdateCharacter(uuid string, character data.Character) {
	items := m.list.Items()
	
	for i, item := range items {
		if existingItem, ok := item.(characterItem); ok && existingItem.uuid == uuid {
			// Update existing item
			updatedItem := characterItem{
				uuid: uuid, 
				Character: character,
			}
			m.list.SetItem(i, updatedItem)
			return
		}
	}
	
	// If not found, add as new item
	newItem := characterItem{
		uuid: uuid, 
		Character: character,
	}
	m.list.InsertItem(len(items), newItem)
}

func (m *CharacterListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)

		return m, nil
	case tea.KeyMsg:
		additionalKeys := newAdditionalCharacterListKeyMap()
		if key.Matches(msg, additionalKeys.newCharacter) {
			return m, tea.Cmd(func() tea.Msg {
				return characterEditedMsg{uuid: ""}
			})
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m CharacterListModel) View() string {
	return m.list.View()
}

type characterItem struct {
	uuid string
	data.Character
}

func (c characterItem) FilterValue() string { return c.Name }

type characterItemDelegate struct {
	keys characterItemKeyMap
}

func (d characterItemDelegate) Height() int  { return 1 }
func (d characterItemDelegate) Spacing() int { return 0 }
func (d *characterItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// Update key states based on item count
	hasItems := len(m.Items()) > 0
	d.keys.view.SetEnabled(hasItems)
	d.keys.edit.SetEnabled(hasItems)
	d.keys.delete.SetEnabled(hasItems)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.edit):
			if item, ok := m.SelectedItem().(characterItem); ok {
				return tea.Cmd(func() tea.Msg {
					return characterEditedMsg{uuid: item.uuid}
				})
			}

		case key.Matches(msg, d.keys.delete):
			if item, ok := m.SelectedItem().(characterItem); ok {
				return tea.Cmd(func() tea.Msg {
					return characterDeletedMsg{uuid: item.uuid}
				})
			}
		}
	}

	return nil
}

func (d characterItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(characterItem)
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

func (d characterItemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		d.keys.view,
		d.keys.edit,
		d.keys.delete,
	}
}

func (d characterItemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.keys.view,
			d.keys.edit,
			d.keys.delete,
		},
	}
}

type characterItemKeyMap struct {
	view   key.Binding
	edit   key.Binding
	delete key.Binding
}

func (d characterItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.view,
		d.edit,
		d.delete,
	}
}

func (d characterItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.view,
			d.edit,
			d.delete,
		},
	}
}

func newCharacterItemKeyMap() characterItemKeyMap {
	return characterItemKeyMap{
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

type additionalCharacterListKeyMap struct {
	newCharacter key.Binding
}

func newAdditionalCharacterListKeyMap() additionalCharacterListKeyMap {
	return additionalCharacterListKeyMap{
		newCharacter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
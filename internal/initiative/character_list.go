package initiative

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

type characterList struct {
	list       list.Model
	characters *map[string]dnd.Character

	width  int
	height int
}

func newCharacterList(characters *map[string]dnd.Character, width int, height int) *characterList {
	cl := &characterList{
		characters: characters,
		width:      width,
		height:     height,
	}

	items := cl.toItems()

	delegate := &characterItemDelegate{
		keys: newCharacterItemKeyMap(),
	}

	// Update key states based on item count
	hasItems := len(items) > 0
	delegate.keys.view.SetEnabled(hasItems)
	delegate.keys.edit.SetEnabled(hasItems)
	delegate.keys.delete.SetEnabled(hasItems)

	l := list.New(items, delegate, width, height)
	l.SetStatusBarItemName("character", "characters")
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.StatusMessageLifetime = time.Duration(1500) * time.Millisecond
	l.StatusMessageLocation = list.InStatusBar

	// Hide the "No items" message
	styles := l.Styles
	styles.NoItems = styles.NoItems.Width(0).Height(0)
	l.Styles = styles

	// Set up key map
	keyMap := list.DefaultKeyMap()
	keyMap.GoToStart = key.NewBinding(key.WithDisabled())
	keyMap.GoToEnd = key.NewBinding(key.WithDisabled())
	keyMap.CursorUp = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	)
	keyMap.CursorDown = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	)
	l.KeyMap = keyMap

	cl.list = l
	return cl
}

func (c *characterList) Init() tea.Cmd {
	return nil
}

func (c *characterList) toItems() []list.Item {
	items := []list.Item{}

	if c.characters != nil {
		for uuid, character := range *c.characters {
			items = append(items, characterItem{uuid: uuid, Character: character})
		}
	}

	// Sort items by character name (case-insensitive)
	sort.Slice(items, func(i, j int) bool {
		charI := items[i].(characterItem)
		charJ := items[j].(characterItem)
		return strings.ToLower(charI.GetName()) < strings.ToLower(charJ.GetName())
	})

	return items
}

func (c *characterList) findCharacterIndex(uuid string, items []list.Item) int {
	for i, item := range items {
		if charItem, ok := item.(characterItem); ok && charItem.uuid == uuid {
			return i
		}
	}
	return -1
}

func (c *characterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case deleteCharacterMsg:
		if c.characters != nil {
			if character, exists := (*c.characters)[msg.uuid]; exists {
				leaveMessage := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(character.GetName() + " has left the party!")
				cmds = append(cmds, c.list.NewStatusMessage(leaveMessage))
			}
			delete(*c.characters, msg.uuid)
		}

		// Rebuild sorted items and maintain selection
		currentIndex := c.list.Index()
		items := c.toItems()
		c.list.SetItems(items)

		// Adjust selection after deletion
		if len(items) > 0 {
			newIndex := min(currentIndex, len(items)-1)
			c.list.Select(newIndex)
		}

		return c, tea.Batch(cmds...)

	case characterUpdatedMsg:
		if c.characters != nil {
			(*c.characters)[msg.uuid] = msg.character

			items := c.toItems()
			c.list.SetItems(items)

			// Find and select the updated character
			if newIndex := c.findCharacterIndex(msg.uuid, items); newIndex != -1 {
				c.list.Select(newIndex)
			}
		}

		return c, nil

	case characterAddedMsg:
		if c.characters == nil {
			newCharacters := make(map[string]dnd.Character)
			c.characters = &newCharacters
		}

		(*c.characters)[msg.uuid] = msg.character

		// Rebuild sorted items and move selection to new character
		items := c.toItems()
		c.list.SetItems(items)

		// Find and select the new character
		if newIndex := c.findCharacterIndex(msg.uuid, items); newIndex != -1 {
			c.list.Select(newIndex)
		}

		joinMessage := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(msg.character.GetName() + " has joined the party!")
		cmds = append(cmds, c.list.NewStatusMessage(joinMessage))

		return c, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c *characterList) View() string {
	c.list.SetHeight(c.height)
	c.list.SetWidth(c.width)
	return c.list.View()
}

func (c *characterList) SetSize(width int, height int) {
	c.width = width
	c.height = height
	c.list.SetWidth(width)
	c.list.SetHeight(height)
}

func (c *characterList) ShortHelp() []key.Binding {
	// Get list navigation keys
	listKeys := []key.Binding{
		c.list.KeyMap.CursorUp,
		c.list.KeyMap.CursorDown,
	}

	// Add delegate keys for character actions (we know we have a characterItemDelegate)
	delegate := newCharacterItemKeyMap()
	if len(c.list.Items()) > 0 {
		listKeys = append(listKeys, delegate.view, delegate.edit, delegate.delete)
	}

	// Add list-specific keys
	if c.list.FilterState() != list.Filtering {
		listKeys = append(listKeys, c.list.KeyMap.Filter)
	}

	return listKeys
}

func (c *characterList) FullHelp() [][]key.Binding {
	// Get list navigation keys
	navKeys := []key.Binding{
		c.list.KeyMap.CursorUp,
		c.list.KeyMap.CursorDown,
	}

	actionKeys := []key.Binding{
		c.list.KeyMap.Filter,
	}

	// Add delegate keys for character actions if items exist
	if len(c.list.Items()) > 0 {
		delegate := newCharacterItemKeyMap()
		actionKeys = append(actionKeys, delegate.view, delegate.edit, delegate.delete)
	}

	return [][]key.Binding{navKeys, actionKeys}
}

// Messages
type viewCharacterMsg struct {
	uuid string
}

type editCharacterMsg struct {
	uuid string
}

type deleteCharacterMsg struct {
	uuid string
}

type characterUpdatedMsg struct {
	uuid      string
	character dnd.Character
}

type characterAddedMsg struct {
	uuid      string
	character dnd.Character
}

// List item for characters
type characterItem struct {
	uuid string
	dnd.Character
}

func (c characterItem) FilterValue() string { return c.GetName() }

// List delegate for characters
type characterItemDelegate struct {
	keys characterItemKeyMap
}

func (c characterItemDelegate) Height() int  { return 1 }
func (c characterItemDelegate) Spacing() int { return 0 }

func (c *characterItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// Update key states based on item count
	hasItems := len(m.Items()) > 0
	c.keys.view.SetEnabled(hasItems)
	c.keys.edit.SetEnabled(hasItems)
	c.keys.delete.SetEnabled(hasItems)

	item, ok := m.SelectedItem().(characterItem)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.view):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return viewCharacterMsg{uuid: item.uuid}
				})
			}

		case key.Matches(msg, c.keys.edit):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return editCharacterMsg{uuid: item.uuid}
				})
			}

		case key.Matches(msg, c.keys.delete):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return deleteCharacterMsg{uuid: item.uuid}
				})
			}
		}
	}

	return nil
}

func (c characterItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(characterItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.GetName())

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (c characterItemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		c.keys.view,
		c.keys.edit,
		c.keys.delete,
	}
}

func (c characterItemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			c.keys.view,
			c.keys.edit,
			c.keys.delete,
		},
	}
}

// Key mappings for character items
type characterItemKeyMap struct {
	view   key.Binding
	edit   key.Binding
	delete key.Binding
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

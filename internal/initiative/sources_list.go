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

type sourcesList struct {
	list    list.Model
	sources map[string]*dnd.Source

	width  int
	height int
}

func newSourcesList(sources map[string]*dnd.Source, width int, height int) *sourcesList {
	sl := &sourcesList{
		sources: sources,
		width:   width,
		height:  height,
	}

	items := sl.toItems()

	delegate := &sourceItemDelegate{
		keys: newSourceItemKeyMap(),
	}

	// Update key states based on item count
	hasItems := len(items) > 0
	delegate.keys.view.SetEnabled(hasItems)

	l := list.New(items, delegate, width, height)
	l.Title = "Sources"
	l.SetStatusBarItemName("source", "sources")
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.StatusMessageLifetime = time.Duration(1500) * time.Millisecond
	l.Styles = currentTheme.list

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

	sl.list = l
	return sl
}

func (s *sourcesList) Init() tea.Cmd {
	return nil
}

func (s *sourcesList) toItems() []list.Item {
	items := []list.Item{}

	for _, source := range s.sources {
		items = append(items, sourceItem{Source: *source})
	}

	// Sort items by source name (case-insensitive)
	sort.Slice(items, func(i, j int) bool {
		sourceI := items[i].(sourceItem)
		sourceJ := items[j].(sourceItem)
		return strings.ToLower(sourceI.Meta.Name) < strings.ToLower(sourceJ.Meta.Name)
	})

	return items
}

func (s *sourcesList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *sourcesList) View() string {
	s.list.SetHeight(s.height)
	s.list.SetWidth(s.width)
	return s.list.View()
}

func (s *sourcesList) SetSize(width int, height int) {
	s.width = width
	s.height = height
	s.list.SetWidth(width)
	s.list.SetHeight(height)
}

func (s *sourcesList) ShortHelp() []key.Binding {
	// Get list navigation keys
	listKeys := []key.Binding{
		s.list.KeyMap.CursorUp,
		s.list.KeyMap.CursorDown,
	}

	// Add delegate keys for source actions
	delegate := newSourceItemKeyMap()
	if len(s.list.Items()) > 0 {
		listKeys = append(listKeys, delegate.view)
	}

	// Add list-specific keys
	if s.list.FilterState() != list.Filtering {
		listKeys = append(listKeys, s.list.KeyMap.Filter)
	}

	return listKeys
}

func (s *sourcesList) FullHelp() [][]key.Binding {
	// Get list navigation keys
	navKeys := []key.Binding{
		s.list.KeyMap.CursorUp,
		s.list.KeyMap.CursorDown,
	}

	actionKeys := []key.Binding{
		s.list.KeyMap.Filter,
	}

	// Add delegate keys for source actions if items exist
	if len(s.list.Items()) > 0 {
		delegate := newSourceItemKeyMap()
		actionKeys = append(actionKeys, delegate.view)
	}

	return [][]key.Binding{navKeys, actionKeys}
}

// List item for sources
type sourceItem struct {
	dnd.Source
}

func (s sourceItem) FilterValue() string {
	return s.Meta.Key + " " + s.Meta.Name
}

// List delegate for sources
type sourceItemDelegate struct {
	keys sourceItemKeyMap
}

func (s sourceItemDelegate) Height() int  { return 1 }
func (s sourceItemDelegate) Spacing() int { return 0 }

func (s *sourceItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// Update key states based on item count
	hasItems := len(m.Items()) > 0
	s.keys.view.SetEnabled(hasItems)

	item, ok := m.SelectedItem().(sourceItem)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.view):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return viewSourceMsg{key: item.Meta.Key}
				})
			}
		}
	}

	return nil
}

func (s sourceItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(sourceItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Meta.Name)

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (s sourceItemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{
		s.keys.view,
	}
}

func (s sourceItemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			s.keys.view,
		},
	}
}

// Key mappings for source items
type sourceItemKeyMap struct {
	view key.Binding
}

func newSourceItemKeyMap() sourceItemKeyMap {
	return sourceItemKeyMap{
		view: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
	}
}

func (d sourceItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.view,
	}
}

func (d sourceItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.view,
		},
	}
}

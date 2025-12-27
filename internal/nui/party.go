package nui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

var _ tea.Model = (*encounter)(nil)

type partyView int

const (
	partyList partyView = iota
	partyDetail
	partyForm
)

type party struct {
	skeleton *skeleton.Skeleton
	party    *map[string]Character

	view partyView

	list       list.Model
	listKeys   additionalGameListKeyMap
	form       *huh.Form
	detailKeys partyDetailKeyMap
	help       help.Model

	// the uuid of the character currently being viewed or edited
	character string
}

func newParty(s *skeleton.Skeleton, p *map[string]Character) *party {
	items := []list.Item{}

	if p != nil {
		for uuid, character := range *p {
			items = append(
				items,
				characterItem{uuid: uuid, Character: character},
			)
		}
	}

	characterItemKeyMap := newCharacterItemKeyMap()
	if len(items) == 0 {
		characterItemKeyMap.view.SetEnabled(false)
		characterItemKeyMap.edit.SetEnabled(false)
		characterItemKeyMap.delete.SetEnabled(false)
	}

	l := list.New(items, &characterItemDelegate{keys: characterItemKeyMap}, s.GetContentWidth(), s.GetContentHeight())

	additionalPartyListKeyMap := newAdditionalPartyListKeyMap()
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{additionalPartyListKeyMap.newCharacter}
	}
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{additionalPartyListKeyMap.newCharacter}
	}
	l.KeyMap = newPartyListKeyMap()

	l.SetStatusBarItemName("character", "characters")
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	return &party{
		skeleton: s,
		party:    p,

		view: partyList,

		list:       l,
		listKeys:   additionalPartyListKeyMap,
		detailKeys: newPartyDetailKeyMap(),
		help:       help.New(),
	}
}

func (p party) Init() tea.Cmd {
	return nil
}

func (p party) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, p.listKeys.newCharacter) && p.listKeys.newCharacter.Enabled() && p.view == partyList {
			return p, tea.Cmd(func() tea.Msg {
				return editCharacterMsg{uuid: ""}
			})
		}
	case viewCharacterMsg:
		{
			p.character = msg.uuid
			p.view = partyDetail

			var characterName string
			if p.party != nil {
				if character, exists := (*p.party)[msg.uuid]; exists {
					characterName = character.Name()
				}
			}

			p.skeleton.UpdatePageTitle("party", "Party > "+characterName)
			return p, nil
		}
	case editCharacterMsg:
		{
			var name string
			if msg.uuid != "" && p.party != nil {
				if character, exists := (*p.party)[msg.uuid]; exists {
					name = character.Name()
				}
			}
			p.form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Key("name").
						Title("Name").
						Value(&name),
				),
			)
			p.character = msg.uuid
			p.view = partyForm
			return p, p.form.Init()
		}
	case deleteCharacterMsg:
		{
			if p.party != nil {
				delete(*p.party, msg.uuid)
			}

			items := p.list.Items()
			for i, item := range items {
				if charItem, ok := item.(characterItem); ok && charItem.uuid == msg.uuid {
					p.list.RemoveItem(i)
					p.list.Select(max(0, i-1))
					break
				}
			}
			return p, nil
		}
	}

	switch p.view {
	case partyList:
		{
			var cmd tea.Cmd
			p.list, cmd = p.list.Update(msg)

			return p, cmd
		}
	case partyDetail:
		{
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, p.detailKeys.back):
					p.character = ""
					p.view = partyList
					p.skeleton.UpdatePageTitle("party", "Party")
					return p, nil
				}
			}
		}
	case partyForm:
		{
			form, cmd := p.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				p.form = f
			}

			if p.form.State == huh.StateCompleted {
				name := p.form.GetString("name")

				if p.character != "" {
					// 1. editing existing character
					if p.party != nil {
						// Update the character in the map
						character := (*p.party)[p.character]
						character.name = name
						(*p.party)[p.character] = character

						// Find and update the corresponding list item with the updated character
						items := p.list.Items()
						for i, item := range items {
							if charItem, ok := item.(characterItem); ok && charItem.uuid == p.character {
								updatedItem := characterItem{uuid: p.character, Character: character}
								p.list.SetItem(i, updatedItem)
								break
							}
						}
					}
				} else {
					// 2. adding new character - generate new UUID
					character := Character{name: name}
					uuid := fmt.Sprintf("char_%d", len(*p.party))
					if p.party == nil {
						newParty := make(map[string]Character)
						p.party = &newParty
					}
					(*p.party)[uuid] = character
					newItem := characterItem{uuid: uuid, Character: character}
					p.list.InsertItem(len(p.list.Items()), newItem)
				}

				p.view = partyList
				p.character = ""
			}

			return p, cmd
		}
	}

	return p, nil
}

func (p party) View() string {
	switch p.view {
	case partyList:
		p.list.SetHeight(p.skeleton.GetContentHeight())
		p.list.SetWidth(p.skeleton.GetContentWidth())
		return p.list.View()
	case partyDetail:
		var characterName string
		if p.party != nil {
			if character, exists := (*p.party)[p.character]; exists {
				characterName = character.Name()
			}
		}

		// Calculate available height for content
		helpStyle := lipgloss.NewStyle().PaddingBottom(1)
		helpView := helpStyle.Render(p.help.View(p.detailKeys))
		availHeight := p.skeleton.GetContentHeight() - lipgloss.Height(helpView)

		// Create main content area
		content := fmt.Sprintf("Viewing character: %s", characterName)
		contentArea := lipgloss.NewStyle().
			Height(availHeight).
			Width(p.skeleton.GetContentWidth()).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top).
			Render(content)

		return lipgloss.JoinVertical(lipgloss.Left, contentArea, helpView)
	case partyForm:
		p.form.WithHeight(p.skeleton.GetContentHeight()).WithWidth(p.skeleton.GetContentWidth())
		return p.form.View()
	}

	return ""
}

func newPartyListKeyMap() list.KeyMap {
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

type additionalGameListKeyMap struct {
	newCharacter key.Binding
}

func newAdditionalPartyListKeyMap() additionalGameListKeyMap {
	return additionalGameListKeyMap{
		newCharacter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
	}
}

type partyDetailKeyMap struct {
	back key.Binding
}

func newPartyDetailKeyMap() partyDetailKeyMap {
	return partyDetailKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k partyDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.back}
}

func (k partyDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.back},
	}
}

// ------- characterItem
var _ list.Item = (*characterItem)(nil)

type characterItem struct {
	uuid string
	Character
}

func (c characterItem) FilterValue() string { return c.Name() }

// -------- characterItemDelegate
type characterItemDelegate struct {
	keys characterItemKeyMap
}

type viewCharacterMsg struct {
	uuid string
}

type editCharacterMsg struct {
	uuid string
}

type deleteCharacterMsg struct {
	uuid string
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

	str := fmt.Sprintf("%d. %s", index+1, i.Name())

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

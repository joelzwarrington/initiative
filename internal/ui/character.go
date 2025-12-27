package ui

import (
	"initiative/internal/data"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type characterView int

const (
	characterListView characterView = iota
	characterEditView
)

type characterDataChangedMsg struct{}

type tabNavigationMsg struct {
	direction int // 1 for next, -1 for previous
}

type gameBackMsg struct{}

type gameQuitMsg struct{}

var _ tea.Model = (*CharacterModel)(nil)

type CharacterModel struct {
	game *data.Game

	currentView characterView

	list *CharacterListModel
	form *CharacterFormModel

	help   help.Model
	keyMap characterKeyMap

	width  int
	height int
}

type characterKeyMap struct {
	tab      key.Binding
	shiftTab key.Binding
	back     key.Binding
	quit     key.Binding
	help     key.Binding
}

func newCharacterKeyMap() characterKeyMap {
	return characterKeyMap{
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		shiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev tab"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (k characterKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		k.help,
	}
}

func (k characterKeyMap) FullHelp() [][]key.Binding {
	characterHelp := []key.Binding{
		key.NewBinding(key.WithKeys("↑", "k"), key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithKeys("↓", "j"), key.WithHelp("↓/j", "down")),
		key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	}
	
	gameNavHelp := []key.Binding{
		k.tab, k.shiftTab, k.back, k.quit,
	}
	
	return [][]key.Binding{
		append(characterHelp, gameNavHelp...),
	}
}

func newCharacter(game *data.Game, width, height int) *CharacterModel {
	return &CharacterModel{
		game:        game,
		currentView: characterListView,
		list:        newCharacterList(game.Characters),
		help:        help.New(),
		keyMap:      newCharacterKeyMap(),
		width:       width,
		height:      height,
	}
}

func (m *CharacterModel) Init() tea.Cmd {
	return nil
}

func (m *CharacterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.list != nil {
		m.list.list.SetSize(width, height)
	}
}

func (m *CharacterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case characterEditedMsg:
		{
			var character *data.Character
			if msg.uuid != "" {
				if char, exists := m.game.Characters[msg.uuid]; exists {
					character = &char
				}
			}
			m.form = newCharacterForm(msg.uuid, character)
			m.currentView = characterEditView
			return m, m.form.Init()
		}

	case characterSavedMsg:
		{
			if m.game.Characters == nil {
				m.game.Characters = make(map[string]data.Character)
			}

			character := data.Character{
				Name: msg.name,
			}

			characterUUID := msg.uuid
			if characterUUID == "" {
				characterUUID = uuid.New().String()
			}

			m.game.Characters[characterUUID] = character

			// Update the list with the new/updated character
			m.list.UpdateCharacter(characterUUID, character)

			m.currentView = characterListView
			// Return a command to save data  
			return m, func() tea.Msg {
				return characterDataChangedMsg{}
			}
		}

	case characterDeletedMsg:
		{
			if m.game.Characters != nil {
				delete(m.game.Characters, msg.uuid)
			}
			m.list.RemoveCharacter(msg.uuid)
			m.currentView = characterListView
			// Return a command to save data  
			return m, func() tea.Msg {
				return characterDataChangedMsg{}
			}
		}

	case characterDeselectedMsg:
		{
			m.currentView = characterListView
			m.form = nil
			return m, nil
		}

	case tea.KeyMsg:
		// Handle global character tab keys
		switch {
		case key.Matches(msg, m.keyMap.tab):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: 1}
			}
		case key.Matches(msg, m.keyMap.shiftTab):
			return m, func() tea.Msg {
				return tabNavigationMsg{direction: -1}
			}
		case key.Matches(msg, m.keyMap.back):
			return m, func() tea.Msg {
				return gameBackMsg{}
			}
		case key.Matches(msg, m.keyMap.quit):
			return m, func() tea.Msg {
				return gameQuitMsg{}
			}
		case key.Matches(msg, m.keyMap.help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}

	switch m.currentView {
	case characterListView:
		var cmd tea.Cmd
		l, cmd := m.list.Update(msg)
		if lm, ok := l.(*CharacterListModel); ok {
			m.list = lm
		}
		return m, cmd
	case characterEditView:
		var cmd tea.Cmd
		f, cmd := m.form.Update(msg)
		if fm, ok := f.(*CharacterFormModel); ok {
			m.form = fm
		}
		return m, cmd
	default:
		return m, nil
	}
}

func (m CharacterModel) View() string {
	switch m.currentView {
	case characterListView:
		return m.list.View()
	case characterEditView:
		return m.form.View()
	default:
		return "Unknown character view state"
	}
}


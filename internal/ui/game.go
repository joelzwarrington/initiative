package ui

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
)

type GameModelView int

const (
	GameView GameModelView = iota
	GameListView
	GameEditView
)

var _ tea.Model = (*GameModel)(nil)

type GameModel struct {
	*data.Data
	game *data.Game

	currentView GameModelView

	list *GameListModel
	form *GameFormModel
}

func newGame(d *data.Data) GameModel {
	var games map[string]data.Game

	if d != nil {
		games = d.Games
	}

	return GameModel{
		Data: d,
		game: nil, // there is no selected game when program begins

		currentView: GameListView,

		list: newGameList(games),
	}
}

func (m GameModel) Init() tea.Cmd {
	return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewGameMsg:
		{
			m.game = m.getGame(msg.uuid)
			m.currentView = GameView
			return m, nil
		}
	case editGameMsg:
		{
			m.game = m.getGame(msg.uuid)
			m.currentView = GameEditView
			m.form = newGameForm(m.game)

			return m, nil
		}
	case deleteGameMsg:
		{
			m.game = nil
			if m.Data != nil && m.Games != nil {
				delete(m.Games, msg.uuid)
				m.Save()
			}
			m.list.RemoveGameByUUID(msg.uuid)
			return m, nil
		}
	}

	switch m.currentView {
	case GameListView:
		var cmd tea.Cmd
		l, cmd := m.list.Update(msg)
		if lm, ok := l.(*GameListModel); ok {
			m.list = lm
		}
		return m, cmd
	case GameView:
		return m, nil
	default:
		return m, nil
	}
}

func (m GameModel) View() string {
	switch m.currentView {
	case GameListView:
		return m.list.View()
	case GameView:
		if m.game != nil {
			return "Game:" + m.game.Name
		}

		return "Game: None selected"
	default:
		return "No view"
	}
}

func (m GameModel) getGame(uuid string) *data.Game {
	if uuid == "" || m.Data == nil || m.Games == nil {
		return nil
	}

	if game, exists := m.Games[uuid]; exists {
		return &game
	}

	return nil
}

package ui

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type gameDeselectedMsg struct{}

type appView int

const (
	GameView appView = iota
	GameListView
	GameEditView
)

var _ tea.Model = (*AppModel)(nil)

type AppModel struct {
	*data.Data
	currentGame     *data.Game
	currentGameUUID string

	currentView appView

	game *GameModel
	list *GameListModel
	form *GameFormModel

	width  int
	height int
}

func newApp(d *data.Data) AppModel {
	var games map[string]data.Game

	if d != nil {
		games = d.Games
	}

	return AppModel{
		Data: d,
		game: nil, // there is no selected game when app begins

		currentView: GameListView,

		list: newGameList(games),
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		break
	case gameViewedMsg:
		{
			m.currentGame = m.getGame(msg.uuid)
			m.currentGameUUID = msg.uuid
			m.currentView = GameView
			m.game = newGame(m.currentGame, m.width, m.height)
			return m, m.game.Init()
		}
	case gameEditedMsg:
		{
			m.currentGame = m.getGame(msg.uuid)
			m.form = newGameForm(msg.uuid, m.currentGame)
			m.currentView = GameEditView

			return m, m.form.Init()
		}
	case gameSavedMsg:
		{
			if m.Data != nil {
				if m.Games == nil {
					m.Games = make(map[string]data.Game)
				}

				game := data.Game{
					Name: msg.name,
				}

				gameUUID := msg.uuid
				if gameUUID == "" {
					gameUUID = uuid.New().String()
				}

				m.Games[gameUUID] = game
				m.Save()

				// Update the list with the new/updated game
				m.list.UpdateGame(gameUUID, game)

				m.currentGame = &game
			}
			m.currentView = GameView
			m.game = newGame(m.currentGame, m.width, m.height)
			return m, m.game.Init()
		}
	case gameDeletedMsg:
		{
			m.game = nil
			if m.Data != nil && m.Games != nil {
				delete(m.Games, msg.uuid)
				m.Save()
			}
			m.list.RemoveGame(msg.uuid)
			return m, nil
		}
	case gameDeselectedMsg:
		{
			m.currentView = GameListView
			m.currentGame = nil
			m.currentGameUUID = ""
			m.game = nil
			m.form = nil
			return m, nil
		}
	case gameDataChangedMsg:
		{
			// Game data (including characters) has changed, save it
			if m.Data != nil {
				// IMPORTANT: Copy the current game data back to the map
				// since we're working with a pointer to a copy
				if m.currentGame != nil && m.currentGameUUID != "" {
					m.Games[m.currentGameUUID] = *m.currentGame
				}
				m.Save()
			}
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
		g, cmd := m.game.Update(msg)
		if g, ok := g.(*GameModel); ok {
			m.game = g
		}
		return m, cmd
	case GameEditView:
		var cmd tea.Cmd
		f, cmd := m.form.Update(msg)
		if fm, ok := f.(*GameFormModel); ok {
			m.form = fm
		}
		return m, cmd
	default:
		return m, nil
	}
}

func (m AppModel) View() string {
	switch m.currentView {
	case GameListView:
		return m.list.View()
	case GameEditView:
		return m.form.View()
	case GameView:
		return m.game.View()
	default:
		return "No view"
	}
}

func (m AppModel) getGame(uuid string) *data.Game {
	if uuid == "" || m.Data == nil || m.Games == nil {
		return nil
	}

	if game, exists := m.Games[uuid]; exists {
		return &game
	}

	return nil
}

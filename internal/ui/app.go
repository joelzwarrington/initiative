package ui

import (
	"initiative/internal/game"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	currentGame *game.Game
	games       []game.Game

	currentView view
	gameList    gameList
}

func newApp() app {
	games := []game.Game{}

	return app{
		currentView: GameList,
		games:       []game.Game{},
		gameList:    newGameList(nil, games),
	}
}

func NewProgram() *tea.Program {
	app := newApp()
	return tea.NewProgram(app)
}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.gameList.SetSize(msg.Width, msg.Height)
		return a, nil
	case tea.KeyMsg:

		switch a.currentView {
		case GameList:
			{
				switch {
				case key.Matches(msg, keys.Help):
					a.gameList.SetShowHelp(!a.gameList.ShowHelp())
					return a, nil
				case key.Matches(msg, keys.New):
					a.currentView = NewGameForm
					return a, nil
				}
			}
		case NewGameForm:
			{
				switch {
				case key.Matches(msg, keys.Back):
					a.currentView = GameList
					return a, nil
				}
			}
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a app) View() string {
	switch a.currentView {
	case GameList:
		return a.gameList.View()
	case NewGameForm:
		return "New game form!"
	}

	return ""
}

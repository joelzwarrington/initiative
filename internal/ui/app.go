package ui

import (
	"initiative/internal/game"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type app struct {
	currentGame *game.Game
	games       []game.Game

	currentView view
	gameList    gameList
	gameForm    gameForm
}

func newApp() app {
	games := []game.Game{}

	return app{
		currentView: GameList,
		games:       games,
		gameList:    newGameList(nil, games),
		gameForm:    newGameForm(),
	}
}

func NewProgram() *tea.Program {
	app := newApp()
	return tea.NewProgram(app)
}

func (a app) Init() tea.Cmd {
	a.gameForm.Init()
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch a.currentView {
	case GameList:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			a.gameList.SetSize(msg.Width, msg.Height)
			return a, nil
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.New):
				a.currentView = NewGameForm
				return a, nil
			case key.Matches(msg, keys.Quit):
				return a, tea.Quit
			}
		}
		var cmd tea.Cmd
		a.gameList.Model, cmd = a.gameList.Update(msg)
		return a, cmd
	case NewGameForm:
		form, cmd := a.gameForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			a.gameForm.Form = *f
		}

		if a.gameForm.State == huh.StateCompleted {
			name := a.gameForm.GetString("name")
			newGame := game.Game{Name: name}
			a.games = append(a.games, newGame)
			a.gameList.SetItems(game.ToListItems(a.games))
			a.gameForm = newGameForm()
			a.gameForm.Init()
			a.currentView = GameList
			return a, nil
		}

		return a, cmd
	}
	return a, nil
}

func (a app) View() string {
	switch a.currentView {
	case GameList:
		return a.gameList.View()
	case NewGameForm:
		return a.gameForm.View()
	}

	return ""
}

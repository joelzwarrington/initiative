package ui

import (
	"initiative/internal/data"
	"initiative/internal/ui/views"

	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	GameList view = iota
	NewGameForm
	ShowGame
)


type app struct {
	currentGame *data.Game
	games       []data.Game

	currentView     view
	gameListModel   *views.GameListModel
	gameFormModel   *views.GameNewFormModel
	gamePageModel   *views.GamePageModel
}

func newApp() app {
	games := []data.Game{}
	var currentGame *data.Game

	return app{
		currentView:     GameList,
		games:           games,
		currentGame:     currentGame,
		gameListModel:   views.NewGameListModel(&games, &currentGame),
		gameFormModel:   views.NewGameNewFormModel(&games),
		gamePageModel:   views.NewGamePageModel(),
	}
}

func NewProgram() *tea.Program {
	app := newApp()
	return tea.NewProgram(app)
}

func (a app) Init() tea.Cmd {
	return a.gameFormModel.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle navigation messages first
	switch msg := msg.(type) {
	case views.NavigateToGameListMsg:
		a.currentView = GameList
		a.gameListModel.RefreshItems()
		return a, nil
	case views.NavigateToNewGameFormMsg:
		a.currentView = NewGameForm
		return a, nil
	case views.NavigateToShowGameMsg:
		a.currentGame = msg.Game
		a.gamePageModel.SetCurrentGame(msg.Game)
		a.currentView = ShowGame
		return a, nil
	}

	// Delegate to current view
	var cmd tea.Cmd
	switch a.currentView {
	case GameList:
		var model tea.Model
		model, cmd = a.gameListModel.Update(msg)
		a.gameListModel = model.(*views.GameListModel)
	case ShowGame:
		var model tea.Model
		model, cmd = a.gamePageModel.Update(msg)
		a.gamePageModel = model.(*views.GamePageModel)
	case NewGameForm:
		var model tea.Model
		model, cmd = a.gameFormModel.Update(msg)
		a.gameFormModel = model.(*views.GameNewFormModel)
	}

	return a, cmd
}

func (a app) View() string {
	switch a.currentView {
	case GameList:
		return a.gameListModel.View()
	case NewGameForm:
		return a.gameFormModel.View()
	case ShowGame:
		return a.gamePageModel.View()
	}

	return ""
}

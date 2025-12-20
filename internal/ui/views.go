package ui

import (
	"initiative/internal/game"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type view int

const (
	GameList view = iota
	NewGameForm
)

type gameList struct {
	list.Model
	currentGame *game.Game
	games       *[]game.Game
}

func newGameList(currentGame *game.Game, games []game.Game) gameList {
	gameList := gameList{
		Model: list.New(
			game.ToListItems(games),
			list.NewDefaultDelegate(),
			0,
			0,
		),
		currentGame: currentGame,
	}

	gameList.Title = "Games"
	gameList.SetShowTitle(true)
	gameList.SetShowStatusBar(false)

	gameList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.New,
		}
	}

	return gameList
}

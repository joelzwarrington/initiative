package ui

import (
	"initiative/internal/game"
	"testing"
)

func TestViewConstants(t *testing.T) {
	if GameList != 0 {
		t.Error("GameList should be 0")
	}
	if NewGameForm != 1 {
		t.Error("NewGameForm should be 1")
	}
}

func TestNewGameList(t *testing.T) {

	t.Run("assigns items to games", func(t *testing.T) {
		games := []game.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
		}

		gameList := newGameList(nil, games)

		if gameList.Title != "Games" {
			t.Errorf("Expected title 'Games', got %s", gameList.Title)
		}

		items := gameList.Items()
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}

		gameList = newGameList(nil, []game.Game{})
		items = gameList.Items()
		if len(items) != 0 {
			t.Errorf("Expected 0 items, got %d", len(items))
		}
	})

	t.Run("assigns current game", func(t *testing.T) {
		games := []game.Game{
			{Name: "Game 1"},
		}

		gameList := newGameList(&games[0], games)
		if gameList.currentGame == nil {
			t.Error("Expected currentGame to be set on the list")
		}
	})
}

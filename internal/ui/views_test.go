package ui

import (
	"initiative/internal/game"
	"strings"
	"testing"
)

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

func TestShowGame(t *testing.T) {
	t.Run("shows no game selected when currentGame is nil", func(t *testing.T) {
		testApp := app{currentGame: nil}
		result := testApp.showGame()

		if result != "No game selected" {
			t.Errorf("Expected 'No game selected', got %s", result)
		}
	})

	t.Run("shows game name and help when currentGame is set", func(t *testing.T) {
		testGame := &game.Game{Name: "Test Game"}
		testApp := app{currentGame: testGame}
		result := testApp.showGame()

		if result == "" {
			t.Error("Expected non-empty result when currentGame is set")
		}

		if !strings.Contains(result, "Test Game") {
			t.Error("Expected game name to be in output")
		}

		if !strings.Contains(result, "esc") {
			t.Error("Expected help text to include escape key")
		}
	})
}

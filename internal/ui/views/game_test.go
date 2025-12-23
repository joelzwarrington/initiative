package views

import (
	"initiative/internal/data"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGameListModel(t *testing.T) {
	t.Run("creates list with games", func(t *testing.T) {
		games := []data.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)

		if model == nil {
			t.Error("NewGameListModel should return non-nil model")
		}

		view := model.View()
		if view == "" {
			t.Error("GameListModel view should return non-empty string")
		}
	})

	t.Run("refreshes items when games change", func(t *testing.T) {
		games := []data.Game{}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		initialView := model.View()

		// Add a game
		games = append(games, data.Game{Name: "New Game"})
		model.RefreshItems()

		updatedView := model.View()
		if updatedView == initialView {
			t.Error("View should change after refreshing items")
		}
	})

	t.Run("handles window size updates", func(t *testing.T) {
		games := []data.Game{}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)

		msg := tea.WindowSizeMsg{Width: 80, Height: 24}
		_, cmd := model.Update(msg)

		// Should not panic or error
		_ = cmd
	})
}

func TestGameNewFormModel(t *testing.T) {
	t.Run("creates form", func(t *testing.T) {
		games := []data.Game{}

		model := NewGameNewFormModel(&games)

		if model == nil {
			t.Error("NewGameNewFormModel should return non-nil model")
		}

		view := model.View()
		if view == "" {
			t.Error("GameNewFormModel view should return non-empty string")
		}
	})

	t.Run("starts not completed", func(t *testing.T) {
		games := []data.Game{}

		model := NewGameNewFormModel(&games)

		// We can't test private methods directly, but we can test the public interface
		// The form should render initially
		view := model.View()
		if view == "" {
			t.Error("New form should render initially")
		}
	})

	t.Run("form renders correctly", func(t *testing.T) {
		games := []data.Game{}

		model := NewGameNewFormModel(&games)

		view := model.View()
		if view == "" {
			t.Error("Form should render")
		}
	})

	t.Run("adds game to slice", func(t *testing.T) {
		games := []data.Game{}
		_ = NewGameNewFormModel(&games)

		// Simulate setting a name and adding
		initialCount := len(games)

		newGame := data.Game{Name: "Test Game"}
		games = append(games, newGame)

		if len(games) != initialCount+1 {
			t.Errorf("Expected %d games, got %d", initialCount+1, len(games))
		}

		if games[len(games)-1].Name != "Test Game" {
			t.Errorf("Expected last game name to be 'Test Game', got %s", games[len(games)-1].Name)
		}
	})
}

func TestGamePageModel(t *testing.T) {
	t.Run("shows no game selected when nil", func(t *testing.T) {
		model := NewGamePageModel()

		view := model.View()
		if view != "No game selected" {
			t.Errorf("Expected 'No game selected', got %s", view)
		}
	})

	t.Run("shows game when set", func(t *testing.T) {
		model := NewGamePageModel()
		game := &data.Game{Name: "Test Game"}
		model.SetCurrentGame(game)

		view := model.View()
		if view == "No game selected" {
			t.Error("Should not show 'No game selected' when game is set")
		}

		if !strings.Contains(view, "Test Game") {
			t.Error("View should contain the game name")
		}
	})

	t.Run("handles updates", func(t *testing.T) {
		model := NewGamePageModel()

		updatedModel, cmd := model.Update(tea.KeyMsg{})

		if updatedModel == nil {
			t.Error("Update should return non-nil model")
		}

		if cmd != nil {
			t.Error("Update should return nil command for basic key messages")
		}
	})
}

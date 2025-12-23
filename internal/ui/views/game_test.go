package views

import (
	"initiative/internal/data"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGameListModel(t *testing.T) {
	t.Run("creates list with games", func(t *testing.T) {
		games := []data.Game{{Name: "Game 1"}, {Name: "Game 2"}}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)

		if model == nil {
			t.Error("NewGameListModel should return non-nil model")
		}
		if model.View() == "" {
			t.Error("GameListModel view should return non-empty string")
		}
	})

	t.Run("handles window size updates", func(t *testing.T) {
		games := []data.Game{}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)

		_, cmd := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		_ = cmd // Should not panic
	})

	t.Run("deletes selected game", func(t *testing.T) {
		games := []data.Game{{Name: "Game 1"}, {Name: "Game 2"}, {Name: "Game 3"}}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)
		
		model.list.Select(1) // Select "Game 2"
		model.deleteSelectedGame()
		
		if len(games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(games))
		}
		if games[0].Name != "Game 1" || games[1].Name != "Game 3" {
			t.Error("Wrong game was deleted")
		}
	})

	t.Run("starts new game editing", func(t *testing.T) {
		games := []data.Game{{Name: "Existing Game"}}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)
		
		model.startNewGame()
		
		if len(games) != 2 {
			t.Errorf("Expected 2 games after starting new game, got %d", len(games))
		}
		if !model.editing {
			t.Error("Model should be in editing mode")
		}
		if model.editingIndex != 1 {
			t.Errorf("Expected editing index 1, got %d", model.editingIndex)
		}
	})

	t.Run("starts editing selected game", func(t *testing.T) {
		games := []data.Game{{Name: "Original Name"}}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)
		
		model.startEditingSelected()
		
		if !model.editing {
			t.Error("Model should be in editing mode")
		}
		if model.editingIndex != 0 {
			t.Errorf("Expected editing index 0, got %d", model.editingIndex)
		}
		if model.textInput.Value() != "Original Name" {
			t.Errorf("Expected text input to contain 'Original Name', got '%s'", model.textInput.Value())
		}
	})

	t.Run("finishes editing with new name", func(t *testing.T) {
		games := []data.Game{{Name: "Old Name"}}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)
		
		model.startEditingSelected()
		model.finishEditing("New Name")
		
		if model.editing {
			t.Error("Model should not be in editing mode after finishing")
		}
		if games[0].Name != "New Name" {
			t.Errorf("Expected game name 'New Name', got '%s'", games[0].Name)
		}
	})

	t.Run("cancels new game editing", func(t *testing.T) {
		games := []data.Game{{Name: "Existing Game"}}
		var currentGame *data.Game
		model := NewGameListModel(&games, &currentGame)
		
		model.startNewGame()
		model.cancelEditing()
		
		if len(games) != 1 {
			t.Errorf("Expected 1 game after canceling new game, got %d", len(games))
		}
		if model.editing {
			t.Error("Model should not be in editing mode after canceling")
		}
	})
}

// Note: GameNewFormModel functionality is now integrated into GameListModel

func TestGamePageModel(t *testing.T) {
	t.Run("shows no game selected initially", func(t *testing.T) {
		model := NewGamePageModel()
		if model.View() != "No game selected" {
			t.Error("Should show 'No game selected' when no game is set")
		}
	})

	t.Run("displays game name when set", func(t *testing.T) {
		model := NewGamePageModel()
		model.SetCurrentGame(&data.Game{Name: "Test Game"})

		view := model.View()
		if !strings.Contains(view, "Test Game") {
			t.Error("View should contain the game name")
		}
	})

	t.Run("handles updates without errors", func(t *testing.T) {
		model := NewGamePageModel()
		updatedModel, cmd := model.Update(tea.KeyMsg{})

		if updatedModel == nil {
			t.Error("Update should return non-nil model")
		}
		_ = cmd // Should not panic
	})
}

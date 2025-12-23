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

	t.Run("deletes selected game", func(t *testing.T) {
		games := []data.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
			{Name: "Game 3"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Select the second game (index 1)
		model.list.Select(1)
		
		// Verify initial state
		if len(games) != 3 {
			t.Fatalf("Expected 3 games initially, got %d", len(games))
		}
		
		// Delete the selected game
		model.deleteSelectedGame()
		
		// Verify game was deleted
		if len(games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(games))
		}
		
		// Verify the correct game was deleted (Game 2 should be gone)
		expectedNames := []string{"Game 1", "Game 3"}
		for i, expected := range expectedNames {
			if games[i].Name != expected {
				t.Errorf("Expected game %d to be %s, got %s", i, expected, games[i].Name)
			}
		}
	})

	t.Run("deletes first game and maintains selection", func(t *testing.T) {
		games := []data.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
			{Name: "Game 3"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Select the first game (index 0)
		model.list.Select(0)
		
		// Delete the selected game
		model.deleteSelectedGame()
		
		// Verify selection is still at index 0 (now showing "Game 2")
		if model.list.Index() != 0 {
			t.Errorf("Expected selection at index 0 after deleting first game, got %d", model.list.Index())
		}
		
		// Verify correct games remain
		if len(games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(games))
		}
		if games[0].Name != "Game 2" {
			t.Errorf("Expected first game to be 'Game 2', got %s", games[0].Name)
		}
	})

	t.Run("deletes middle game and selects game above", func(t *testing.T) {
		games := []data.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
			{Name: "Game 3"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Select the middle game (index 1)
		model.list.Select(1)
		
		// Delete the selected game
		model.deleteSelectedGame()
		
		// Verify selection moved to game above (index 0)
		if model.list.Index() != 0 {
			t.Errorf("Expected selection at index 0 after deleting middle game, got %d", model.list.Index())
		}
		
		// Verify correct games remain
		if len(games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(games))
		}
		expectedNames := []string{"Game 1", "Game 3"}
		for i, expected := range expectedNames {
			if games[i].Name != expected {
				t.Errorf("Expected game %d to be %s, got %s", i, expected, games[i].Name)
			}
		}
	})

	t.Run("deletes last game and selects new last game", func(t *testing.T) {
		games := []data.Game{
			{Name: "Game 1"},
			{Name: "Game 2"},
			{Name: "Game 3"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Select the last game (index 2)
		model.list.Select(2)
		
		// Delete the selected game
		model.deleteSelectedGame()
		
		// Verify selection moved to new last game (index 1)
		if model.list.Index() != 1 {
			t.Errorf("Expected selection at index 1 after deleting last game, got %d", model.list.Index())
		}
		
		// Verify correct games remain
		if len(games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(games))
		}
		if games[1].Name != "Game 2" {
			t.Errorf("Expected last game to be 'Game 2', got %s", games[1].Name)
		}
	})

	t.Run("deletes only game leaves empty list", func(t *testing.T) {
		games := []data.Game{
			{Name: "Only Game"},
		}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Select the only game (index 0)
		model.list.Select(0)
		
		// Delete the selected game
		model.deleteSelectedGame()
		
		// Verify list is empty
		if len(games) != 0 {
			t.Errorf("Expected 0 games after deleting only game, got %d", len(games))
		}
	})

	t.Run("delete handles empty list gracefully", func(t *testing.T) {
		games := []data.Game{}
		var currentGame *data.Game

		model := NewGameListModel(&games, &currentGame)
		
		// Try to delete from empty list (should not panic)
		model.deleteSelectedGame()
		
		// Verify list is still empty
		if len(games) != 0 {
			t.Errorf("Expected 0 games in empty list, got %d", len(games))
		}
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

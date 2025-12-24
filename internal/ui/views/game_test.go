package views

import (
	"initiative/internal/data"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestGameModel(t *testing.T) {
	t.Run("creates model with games", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Game 1"}, {Name: "Game 2"}},
		}

		model := NewGameModel(testData)

		if model.View() == "" {
			t.Error("GameModel view should return non-empty string")
		}
	})

	t.Run("handles window size updates", func(t *testing.T) {
		testData := &data.Data{Games: []data.Game{}}
		model := NewGameModel(testData)

		_, cmd := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		_ = cmd // Should not panic
	})

	t.Run("deletes selected game", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Game 1"}, {Name: "Game 2"}, {Name: "Game 3"}},
		}
		model := NewGameModel(testData)
		
		model.list.Select(1) // Select "Game 2"
		
		// Simulate the delegate generating a delete command and then handle it
		gameToDelete := &testData.Games[1]
		deleteMsg := deleteGameMsg{game: gameToDelete}
		updatedModel, _ := model.Update(deleteMsg)
		model = updatedModel.(GameModel)
		
		if len(testData.Games) != 2 {
			t.Errorf("Expected 2 games after deletion, got %d", len(testData.Games))
		}
		if testData.Games[0].Name != "Game 1" || testData.Games[1].Name != "Game 3" {
			t.Error("Wrong game was deleted")
		}
	})

	t.Run("starts new game editing", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Existing Game"}},
		}
		model := NewGameModel(testData)
		
		// Simulate the startNewGameMsg command being sent
		updatedModel, _ := model.Update(startNewGameMsg{})
		model = updatedModel.(GameModel)
		
		if len(testData.Games) != 2 {
			t.Errorf("Expected 2 games after starting new game, got %d", len(testData.Games))
		}
		if !model.isEditing() {
			t.Error("Model should be in editing mode")
		}
	})

	t.Run("starts editing selected game", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Original Name"}},
		}
		model := NewGameModel(testData)
		
		// Simulate the startEditingMsg command being sent for index 0
		updatedModel, _ := model.Update(startEditingMsg{index: 0})
		model = updatedModel.(GameModel)
		
		if !model.isEditing() {
			t.Error("Model should be in editing mode")
		}
		
		// Check the text input value from the editing item
		if selectedItem := model.list.SelectedItem(); selectedItem != nil {
			if gameItem, ok := selectedItem.(gameListItem); ok && gameItem.textInput != nil {
				if gameItem.textInput.Value() != "Original Name" {
					t.Errorf("Expected text input to contain 'Original Name', got '%s'", gameItem.textInput.Value())
				}
			} else {
				t.Error("Expected editing item to have a text input")
			}
		} else {
			t.Error("Expected a selected item")
		}
	})

	t.Run("finishes editing with new name", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Old Name"}},
		}
		model := NewGameModel(testData)
		
		// Start editing
		updatedModel, _ := model.Update(startEditingMsg{index: 0})
		model = updatedModel.(GameModel)
		
		// Set the text input value on the editing item
		if selectedItem := model.list.SelectedItem(); selectedItem != nil {
			if gameItem, ok := selectedItem.(gameListItem); ok && gameItem.textInput != nil {
				gameItem.textInput.SetValue("New Name")
			}
		}
		
		// Finish editing by sending enter key (which should trigger save)
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = updatedModel.(GameModel)
		
		if model.isEditing() {
			t.Error("Model should not be in editing mode after finishing")
		}
		if testData.Games[0].Name != "New Name" {
			t.Errorf("Expected game name 'New Name', got '%s'", testData.Games[0].Name)
		}
	})

	t.Run("cancels new game editing", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Existing Game"}},
		}
		model := NewGameModel(testData)
		
		// Start new game
		updatedModel, _ := model.Update(startNewGameMsg{})
		model = updatedModel.(GameModel)
		
		// Cancel editing
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		model = updatedModel.(GameModel)
		
		if len(testData.Games) != 1 {
			t.Errorf("Expected 1 game after canceling new game, got %d", len(testData.Games))
		}
		if model.isEditing() {
			t.Error("Model should not be in editing mode after canceling")
		}
	})

	t.Run("switches to game view when game selected", func(t *testing.T) {
		testData := &data.Data{
			Games: []data.Game{{Name: "Test Game"}},
		}
		model := NewGameModel(testData)
		
		// Simulate selecting a game
		model.currentGame = &testData.Games[0]
		model.state = gameView

		view := model.View()
		if !strings.Contains(view, "Test Game") {
			t.Error("Game view should contain the game name")
		}
	})

	t.Run("shows no game selected initially in game view", func(t *testing.T) {
		testData := &data.Data{Games: []data.Game{}}
		model := NewGameModel(testData)
		model.state = gameView

		if model.gameView() != "No game selected" {
			t.Error("Should show 'No game selected' when no game is set")
		}
	})

	t.Run("handles key presses without errors", func(t *testing.T) {
		testData := &data.Data{Games: []data.Game{}}
		model := NewGameModel(testData)

		updatedModel, cmd := model.Update(tea.KeyMsg{})

		_ = updatedModel // Should not panic
		_ = cmd          // Should not panic
	})
}

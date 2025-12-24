package views

import (
	"initiative/internal/data"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func TestCharacterModel(t *testing.T) {
	t.Run("creates new character model", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{{Name: "Test Character"}},
		}

		model := NewCharacterModel(game)

		if model.game != game {
			t.Error("CharacterModel should reference the provided game")
		}

		if len(model.list.Items()) != 1 {
			t.Errorf("Expected 1 character in list, got %d", len(model.list.Items()))
		}
	})

	t.Run("displays empty state when no characters", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{},
		}

		model := NewCharacterModel(game)
		view := model.View()

		if !strings.Contains(view, "No characters added yet") {
			t.Error("Should display empty state message when no characters")
		}

		if !strings.Contains(view, "Press 'a' to add a character") {
			t.Error("Should display help text for adding characters")
		}
	})

	t.Run("displays character list when characters exist", func(t *testing.T) {
		game := &data.Game{
			Name: "Test Game",
			Characters: []data.Character{
				{Name: "Character 1"},
				{Name: "Character 2"},
			},
		}

		model := NewCharacterModel(game)
		model.SetSize(80, 24)
		view := model.View()

		if strings.Contains(view, "No characters added yet") {
			t.Error("Should not display empty state when characters exist")
		}

		// The list should be rendered instead of empty state
		if len(view) == 0 {
			t.Error("View should not be empty when characters exist")
		}
	})

	t.Run("adds new character when 'a' key pressed", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{},
		}

		model := NewCharacterModel(game)
		model.SetSize(80, 24)

		// Simulate pressing 'a' key
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		updatedModel, _ := model.Update(keyMsg)

		if len(updatedModel.game.Characters) != 1 {
			t.Errorf("Expected 1 character after adding, got %d", len(updatedModel.game.Characters))
		}

		if updatedModel.game.Characters[0].Name != "New Character" {
			t.Errorf("Expected new character to be named 'New Character', got '%s'", updatedModel.game.Characters[0].Name)
		}

		// Check that the list was updated
		if len(updatedModel.list.Items()) != 1 {
			t.Errorf("Expected 1 item in list after adding, got %d", len(updatedModel.list.Items()))
		}
	})

	t.Run("adds multiple characters", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{{Name: "Existing Character"}},
		}

		model := NewCharacterModel(game)
		model.SetSize(80, 24)

		// Add first character
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		model, _ = model.Update(keyMsg)

		// Add second character
		model, _ = model.Update(keyMsg)

		if len(model.game.Characters) != 3 {
			t.Errorf("Expected 3 characters total, got %d", len(model.game.Characters))
		}

		if len(model.list.Items()) != 3 {
			t.Errorf("Expected 3 items in list, got %d", len(model.list.Items()))
		}
	})

	t.Run("updates window size", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{},
		}

		model := NewCharacterModel(game)

		// Set initial size
		model.SetSize(80, 24)
		if model.width != 80 || model.height != 24 {
			t.Errorf("Expected size 80x24, got %dx%d", model.width, model.height)
		}

		// Handle window size message
		windowMsg := tea.WindowSizeMsg{Width: 120, Height: 30}
		updatedModel, _ := model.Update(windowMsg)

		if updatedModel.width != 120 || updatedModel.height != 30 {
			t.Errorf("Expected size 120x30, got %dx%d", updatedModel.width, updatedModel.height)
		}
	})

	t.Run("converts characters to list items", func(t *testing.T) {
		characters := []data.Character{
			{Name: "Character 1"},
			{Name: "Character 2"},
			{Name: "Character 3"},
		}

		items := charactersToListItems(characters)

		if len(items) != 3 {
			t.Errorf("Expected 3 list items, got %d", len(items))
		}

		// Test first item
		charItem, ok := items[0].(characterListItem)
		if !ok {
			t.Error("List item should be of type characterListItem")
		}

		if charItem.Title() != "Character 1" {
			t.Errorf("Expected title 'Character 1', got '%s'", charItem.Title())
		}

		if charItem.FilterValue() != "Character 1" {
			t.Errorf("Expected filter value 'Character 1', got '%s'", charItem.FilterValue())
		}
	})

	t.Run("handles non-character key presses", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{{Name: "Test Character"}},
		}

		model := NewCharacterModel(game)
		initialCharCount := len(model.game.Characters)

		// Press a different key (not 'a')
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
		updatedModel, _ := model.Update(keyMsg)

		// Character count should remain the same
		if len(updatedModel.game.Characters) != initialCharCount {
			t.Error("Character count should not change for non-add key presses")
		}
	})

	t.Run("keymap configuration", func(t *testing.T) {
		keyMap := newCharacterListKeyMap()

		if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}, keyMap.add) {
			t.Error("'a' key should match add binding")
		}

		if !key.Matches(tea.KeyMsg{Type: tea.KeyEsc}, keyMap.back) {
			t.Error("Esc key should match back binding")
		}

		if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, keyMap.quit) {
			t.Error("'q' key should match quit binding")
		}
	})

	t.Run("delegate keymap configuration", func(t *testing.T) {
		keyMap := newCharacterDelegateKeyMap()

		if !key.Matches(tea.KeyMsg{Type: tea.KeyEnter}, keyMap.choose) {
			t.Error("Enter key should match choose binding")
		}

		if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, keyMap.delete) {
			t.Error("'d' key should match delete binding")
		}
	})
}
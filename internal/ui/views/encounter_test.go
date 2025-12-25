package views

import (
	"initiative/internal/data"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEncounterModel(t *testing.T) {
	t.Run("creates new encounter model", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{{Name: "Test Character"}},
		}

		model := NewEncounterModel(game)

		if model.game != game {
			t.Error("EncounterModel should reference the provided game")
		}

		if model.state != noEncounter {
			t.Error("EncounterModel should start with noEncounter state")
		}

		if model.encounter != nil {
			t.Error("EncounterModel should start with nil encounter")
		}
	})

	t.Run("displays empty state when no encounter", func(t *testing.T) {
		game := &data.Game{
			Name: "Test Game",
			Characters: []data.Character{
				{Name: "Character 1"},
			},
		}

		model := NewEncounterModel(game)
		view := model.View()

		if !strings.Contains(view, "No encounter started yet") {
			t.Error("Should display empty state message when no encounter")
		}
	})

	t.Run("displays no characters message when game has no characters", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{},
		}

		model := NewEncounterModel(game)
		view := model.View()

		if !strings.Contains(view, "No characters available") {
			t.Error("Should display no characters message")
		}

		if !strings.Contains(view, "Add characters first") {
			t.Error("Should prompt to add characters first")
		}
	})

	t.Run("handles nil game gracefully", func(t *testing.T) {
		model := EncounterModel{} // Zero value with nil game

		view := model.View()

		if !strings.Contains(view, "No game loaded") {
			t.Error("Should display no game loaded message when game is nil")
		}
	})

	t.Run("starts encounter creation when 'n' key pressed", func(t *testing.T) {
		game := &data.Game{
			Name: "Test Game",
			Characters: []data.Character{
				{Name: "Character 1"},
				{Name: "Character 2"},
			},
		}

		model := NewEncounterModel(game)

		// Simulate pressing 'n' key
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
		updatedModel, _ := model.Update(keyMsg)

		// For now, starting encounter just returns without doing anything
		if updatedModel.state != noEncounter {
			t.Error("Should remain in noEncounter state for now")
		}
	})

	t.Run("does not start encounter when no characters available", func(t *testing.T) {
		game := &data.Game{
			Name:       "Test Game",
			Characters: []data.Character{},
		}

		model := NewEncounterModel(game)

		// Simulate pressing 'n' key
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
		updatedModel, _ := model.Update(keyMsg)

		// Should still remain in noEncounter state
		if updatedModel.state != noEncounter {
			t.Error("Should remain in noEncounter state when no characters available")
		}
	})


	t.Run("stops encounter when 'x' key pressed", func(t *testing.T) {
		game := &data.Game{
			Name: "Test Game",
			Characters: []data.Character{
				{Name: "Character 1"},
			},
		}

		model := NewEncounterModel(game)
		
		// Set up a running encounter
		model.state = runningEncounter
		model.encounter = &Encounter{
			Name: "Test Encounter",
			InitiativeGroups: []InitiativeGroup{
				{Initiative: 15, Characters: []data.Character{{Name: "Character 1"}}},
			},
		}

		// Simulate pressing 'x' key
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
		updatedModel, _ := model.Update(keyMsg)

		if updatedModel.state != noEncounter {
			t.Error("Should transition to noEncounter state when 'x' is pressed")
		}

		if updatedModel.encounter != nil {
			t.Error("Should clear encounter when stopping")
		}

		// Should be in clean state after stopping
		if updatedModel.state != noEncounter {
			t.Error("Should be in noEncounter state after stopping")
		}
	})

	t.Run("displays running encounter correctly", func(t *testing.T) {
		game := &data.Game{
			Name: "Test Game",
			Characters: []data.Character{
				{Name: "Character 1"},
				{Name: "Character 2"},
			},
		}

		model := NewEncounterModel(game)
		
		// Set up a running encounter
		model.state = runningEncounter
		model.encounter = &Encounter{
			Name: "Test Encounter",
			InitiativeGroups: []InitiativeGroup{
				{Initiative: 18, Characters: []data.Character{{Name: "Character 1"}}},
				{Initiative: 12, Characters: []data.Character{{Name: "Character 2"}}},
			},
		}

		view := model.View()

		if !strings.Contains(view, "Test Encounter") {
			t.Error("Should display encounter name")
		}

		if !strings.Contains(view, "Initiative Order") {
			t.Error("Should display initiative order header")
		}

		if !strings.Contains(view, "Initiative 18") {
			t.Error("Should display first initiative group")
		}

		if !strings.Contains(view, "Initiative 12") {
			t.Error("Should display second initiative group")
		}

		if !strings.Contains(view, "Character 1") {
			t.Error("Should display character names")
		}

		if !strings.Contains(view, "Character 2") {
			t.Error("Should display character names")
		}
	})

	t.Run("handles window size updates", func(t *testing.T) {
		game := &data.Game{Name: "Test Game"}
		model := NewEncounterModel(game)

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


}
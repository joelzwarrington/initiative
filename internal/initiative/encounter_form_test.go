package initiative

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelzwarrington/initiative/dnd"
)

func TestEncounterFormInitialState(t *testing.T) {
	characters := map[string]*dnd.Character{
		"1": dnd.NewCharacter("Alice"),
		"2": dnd.NewCharacter("Bob"),
	}

	form := newEncounterForm(characters, nil, 80, 40)

	// Should start with one form (summary)
	if len(*form.forms) != 1 {
		t.Errorf("expected 1 form, got %d", len(*form.forms))
	}

	// Should have esc help enabled
	helpKeys := form.getHelpKeys()
	hasEscHelp := false
	for _, binding := range helpKeys {
		if binding.Help().Key == "esc" && binding.Help().Desc == "cancel" {
			hasEscHelp = true
			break
		}
	}

	if !hasEscHelp {
		t.Error("expected esc help to be present")
	}
}

func TestEncounterFormNavigatesToInitiativeList(t *testing.T) {
	characters := map[string]*dnd.Character{
		"1": dnd.NewCharacter("Alice"),
		"2": dnd.NewCharacter("Bob"),
	}

	form := newEncounterForm(characters, nil, 80, 40)

	// Set summary value and submit
	currentForm := form.getCurrentForm()
	if field := currentForm.GetFocusedField(); field != nil {
		// Type a summary
		for _, char := range "Test Encounter" {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
			form.Update(keyMsg)
		}
	}

	// Simulate form submission by sending nextStep message
	form.Update(encounterFormNextStepMsg{})

	// Should now have 2 forms (summary + initiative list)
	if len(*form.forms) != 2 {
		t.Errorf("expected 2 forms after next step, got %d", len(*form.forms))
	}

	// Should have initiative list field
	if form.initiativeListField == nil {
		t.Error("expected initiative list field to be created")
	}
}

func TestEncounterFormPreviousStepFromInitiativeList(t *testing.T) {
	characters := map[string]*dnd.Character{
		"1": dnd.NewCharacter("Alice"),
	}

	form := newEncounterForm(characters, nil, 80, 40)

	// Go to initiative list
	form.Update(encounterFormNextStepMsg{})

	if len(*form.forms) != 2 {
		t.Fatalf("expected 2 forms, got %d", len(*form.forms))
	}

	// Go back
	form.Update(encounterFormPreviousStepMsg{})

	// Should be back to 1 form
	if len(*form.forms) != 1 {
		t.Errorf("expected 1 form after going back, got %d", len(*form.forms))
	}

	// Initiative list field should be cleared
	if form.initiativeListField != nil {
		t.Error("expected initiative list field to be cleared")
	}
}

func TestInitiativeListFieldPrePopulatesCharacters(t *testing.T) {
	characters := map[string]*dnd.Character{
		"uuid-1": dnd.NewCharacter("Alice"),
		"uuid-2": dnd.NewCharacter("Bob"),
		"uuid-3": dnd.NewCharacter("Charlie"),
	}

	field := NewInitiativeListField(characters, nil)

	entries, ok := field.GetValue().([]dnd.InitiativeEntry)
	if !ok {
		t.Fatal("expected GetValue to return []dnd.InitiativeEntry")
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	// All entries should be characters
	for _, entry := range entries {
		if entry.CreatureType != "character" {
			t.Errorf("expected creature type 'character', got %q", entry.CreatureType)
		}
		if entry.Quantity != 1 {
			t.Errorf("expected quantity 1 for character, got %d", entry.Quantity)
		}
	}
}

func TestInitiativeListFieldValidation(t *testing.T) {
	characters := map[string]*dnd.Character{
		"1": dnd.NewCharacter("Alice"),
	}

	field := NewInitiativeListField(characters, nil)

	// Should have error initially (no initiative set)
	if err := field.Error(); err == nil {
		t.Error("expected validation error when no initiative set")
	}
}

func TestInitiativeListFieldCharactersSortedAlphabetically(t *testing.T) {
	characters := map[string]*dnd.Character{
		"uuid-3": dnd.NewCharacter("Zara"),
		"uuid-1": dnd.NewCharacter("Alice"),
		"uuid-2": dnd.NewCharacter("Bob"),
	}

	field := NewInitiativeListField(characters, nil)

	entries, ok := field.GetValue().([]dnd.InitiativeEntry)
	if !ok {
		t.Fatal("expected GetValue to return []dnd.InitiativeEntry")
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// Verify alphabetical order
	expectedOrder := []string{"Alice", "Bob", "Zara"}
	for i, entry := range entries {
		if entry.Name != expectedOrder[i] {
			t.Errorf("entry %d: expected name %q, got %q", i, expectedOrder[i], entry.Name)
		}
	}
}

func TestInitiativeListFieldAddMonster(t *testing.T) {
	field := NewInitiativeListField(nil, nil)

	// Initially empty
	entries, _ := field.GetValue().([]dnd.InitiativeEntry)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries initially, got %d", len(entries))
	}

	// Add a monster
	statBlock := dnd.StatBlock{HitPoints: dnd.HitPoints{Fixed: 15}}
	field.AddMonster(dnd.InitiativeEntry{
		CreatureType: "monster",
		CreatureID:   "srd:Goblin",
		Name:         "Goblin",
		Initiative:   12,
		Quantity:     3,
		StatBlock:    &statBlock,
	})

	entries, _ = field.GetValue().([]dnd.InitiativeEntry)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after adding monster, got %d", len(entries))
	}

	entry := entries[0]
	if entry.CreatureType != "monster" {
		t.Errorf("expected creature type 'monster', got %q", entry.CreatureType)
	}
	if entry.Name != "Goblin" {
		t.Errorf("expected name 'Goblin', got %q", entry.Name)
	}
	if entry.Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", entry.Quantity)
	}
	if entry.Initiative != 12 {
		t.Errorf("expected initiative 12, got %d", entry.Initiative)
	}
}

func TestInitiativeListFieldUpdateMonster(t *testing.T) {
	field := NewInitiativeListField(nil, nil)

	// Add a monster
	statBlock := dnd.StatBlock{HitPoints: dnd.HitPoints{Fixed: 15}}
	field.AddMonster(dnd.InitiativeEntry{
		CreatureType: "monster",
		CreatureID:   "srd:Goblin",
		Name:         "Goblin",
		Quantity:     2,
		StatBlock:    &statBlock,
	})

	// Update the monster
	field.UpdateMonster(0, dnd.InitiativeEntry{
		CreatureType: "monster",
		CreatureID:   "srd:Goblin",
		Name:         "Sneaky Goblin",
		Initiative:   18,
		Quantity:     4,
		StatBlock:    &statBlock,
	})

	entries, _ := field.GetValue().([]dnd.InitiativeEntry)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Name != "Sneaky Goblin" {
		t.Errorf("expected name 'Sneaky Goblin', got %q", entry.Name)
	}
	if entry.Quantity != 4 {
		t.Errorf("expected quantity 4, got %d", entry.Quantity)
	}
	if entry.Initiative != 18 {
		t.Errorf("expected initiative 18, got %d", entry.Initiative)
	}
}

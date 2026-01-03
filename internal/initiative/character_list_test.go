package initiative

import (
	"strings"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
)

func TestCharacterListView(t *testing.T) {
	tests := []struct {
		name       string
		characters map[string]*dnd.Character
		expected   string
	}{
		{
			name: "displays sorted character names",
			characters: map[string]*dnd.Character{
				"1": dnd.NewCharacter("Zara"),
				"2": dnd.NewCharacter("alice"),
				"3": dnd.NewCharacter("Bob"),
			},
			expected: `              
  3 characters
              
  > 1. alice  
    2. Bob    
    3. Zara   
              
              `,
		},
		{
			name: "handles single character",
			characters: map[string]*dnd.Character{
				"solo": dnd.NewCharacter("Hero"),
			},
			expected: `             
  1 character
             
  > 1. Hero  
             
             
             
             `,
		},
		{
			name:       "empty list shows no items",
			characters: map[string]*dnd.Character{},
			expected: `               
  No characters
               
No characters. 
               
               
               
               `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := newCharacterList(tt.characters, 15, 8)
			view := cl.View()

			if view != tt.expected {
				t.Errorf("view mismatch:\ngot:\n%q\nwant:\n%q", view, tt.expected)
			}
		})
	}
}

func TestCharacterListUpdate_AddedMessage(t *testing.T) {
	characters := map[string]*dnd.Character{
		"existing": dnd.NewCharacter("Bob"),
	}

	cl := newCharacterList(characters, 80, 12)

	// Add a new character
	newChar := dnd.NewCharacter("Alice")
	msg := characterAddedMsg{uuid: "new", character: newChar}

	updatedModel, cmd := cl.Update(msg)
	cl = updatedModel.(*characterList)

	// Verify the character was added to the underlying map
	if cl.characters["new"].Name() != "Alice" {
		t.Error("character should be added to the characters map")
	}

	// Verify view shows sorted order (Alice before Bob)
	view := cl.View()
	alicePos := strings.Index(view, "Alice")
	bobPos := strings.Index(view, "Bob")

	if alicePos == -1 || bobPos == -1 {
		t.Fatalf("both Alice and Bob should appear in view: %s", view)
	}

	if alicePos > bobPos {
		t.Error("Alice should appear before Bob in sorted view")
	}

	// Should return a status message command
	if cmd == nil {
		t.Error("expected a command for status message")
	}
}

func TestCharacterListUpdate_UpdatedMessage(t *testing.T) {
	characters := map[string]*dnd.Character{
		"char1": dnd.NewCharacter("Bob"),
		"char2": dnd.NewCharacter("Charlie"),
	}

	cl := newCharacterList(characters, 80, 12)

	// Update character name to change sort order
	updatedChar := dnd.NewCharacter("Alice") // This should move to the front
	msg := characterUpdatedMsg{uuid: "char1", character: updatedChar}

	updatedModel, _ := cl.Update(msg)
	cl = updatedModel.(*characterList)

	// Verify the character was updated in the underlying map
	if cl.characters["char1"].Name() != "Alice" {
		t.Error("character should be updated in the characters map")
	}

	// Verify view shows new sorted order (Alice before Charlie)
	view := cl.View()
	if !strings.Contains(view, "1. Alice") {
		t.Error("Alice should now be first in the sorted list")
	}
	if !strings.Contains(view, "2. Charlie") {
		t.Error("Charlie should now be second in the sorted list")
	}
}

func TestCharacterListUpdate_DeletedMessage(t *testing.T) {
	characters := map[string]*dnd.Character{
		"char1": dnd.NewCharacter("Alice"),
		"char2": dnd.NewCharacter("Bob"),
		"char3": dnd.NewCharacter("Charlie"),
	}

	cl := newCharacterList(characters, 80, 12)

	// Delete middle character
	msg := deleteCharacterMsg{uuid: "char2"}

	updatedModel, cmd := cl.Update(msg)
	cl = updatedModel.(*characterList)

	// Verify the character was removed from the underlying map
	if _, exists := cl.characters["char2"]; exists {
		t.Error("character should be removed from the characters map")
	}

	// Verify view shows correct item count and remaining characters
	view := cl.View()
	if !strings.Contains(view, "1. Alice") || !strings.Contains(view, "2. Charlie") {
		t.Errorf("remaining characters should be renumbered correctly: %s", view)
	}

	// Verify only 2 characters are shown (status message about Bob leaving is expected)
	if !strings.Contains(view, "2 characters") {
		t.Error("should show 2 characters after deletion")
	}

	// Should return a status message command
	if cmd == nil {
		t.Error("expected a command for status message")
	}
}

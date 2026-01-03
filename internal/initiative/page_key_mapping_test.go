package initiative

import (
	"testing"

	"github.com/termkit/skeleton"
)

// Test that pages correctly expose key mappings for all their states
func TestPage_KeyMappingConsistency(t *testing.T) {
	tests := []struct {
		name     string
		pageType string
		setup    func() interface{}
	}{
		{
			name:     "encounter page",
			pageType: "encounter",
			setup: func() interface{} {
				s := skeleton.NewSkeleton()
				return newEncounterPage(s, nil, nil)
			},
		},
		{
			name:     "character page",
			pageType: "character",
			setup: func() interface{} {
				s := skeleton.NewSkeleton()
				return newCharacterPage(s, nil)
			},
		},
		{
			name:     "sources page",
			pageType: "sources",
			setup: func() interface{} {
				s := skeleton.NewSkeleton()
				return newSourcesPage(s, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.setup()

			// Test that pages expose key methods - using reflection would be overkill
			// Just verify specific page types work
			switch p := page.(type) {
			case *encounterPage:
				if len(p.ShortHelp()) == 0 {
					t.Errorf("%s page should have short help keys", tt.pageType)
				}
				if len(p.FullHelp()) == 0 {
					t.Errorf("%s page should have full help keys", tt.pageType)
				}
			case *characterPage:
				if len(p.ShortHelp()) == 0 {
					t.Errorf("%s page should have short help keys", tt.pageType)
				}
				if len(p.FullHelp()) == 0 {
					t.Errorf("%s page should have full help keys", tt.pageType)
				}
			case *sourcesPage:
				if len(p.ShortHelp()) == 0 {
					t.Errorf("%s page should have short help keys", tt.pageType)
				}
				if len(p.FullHelp()) == 0 {
					t.Errorf("%s page should have full help keys", tt.pageType)
				}
			}
		})
	}
}

func TestEncounterPage_KeyMappingsChangeByState(t *testing.T) {
	s := skeleton.NewSkeleton()
	page := newEncounterPage(s, nil, nil)

	// Test empty state keys
	emptyKeys := page.ShortHelp()
	hasNewKey := false
	for _, key := range emptyKeys {
		if key.Help().Key == "n" {
			hasNewKey = true
			break
		}
	}
	if !hasNewKey {
		t.Error("Expected 'new encounter' key in empty state")
	}

	// Test form state keys
	page.encounterForm = newEncounterForm(nil, nil, 80, 24)
	formKeys := page.ShortHelp()
	hasEscKey := false
	for _, key := range formKeys {
		if key.Help().Key == "esc" {
			hasEscKey = true
			break
		}
	}
	if !hasEscKey {
		t.Error("Expected 'escape' key in form state")
	}

	// Verify keys are different between states
	if len(emptyKeys) == len(formKeys) {
		same := true
		for i, key := range emptyKeys {
			if i < len(formKeys) && key.Help().Key != formKeys[i].Help().Key {
				same = false
				break
			}
		}
		if same {
			t.Error("Expected different key mappings between empty and form states")
		}
	}
}

func TestCharacterPage_FormKeysExposedCorrectly(t *testing.T) {
	s := skeleton.NewSkeleton()
	page := newCharacterPage(s, nil)

	// Test list state - should have new character key
	listKeys := page.ShortHelp()
	hasNewCharKey := false
	for _, key := range listKeys {
		if key.Help().Key == "n" && key.Help().Desc == "new character" {
			hasNewCharKey = true
			break
		}
	}
	if !hasNewCharKey {
		t.Error("Expected 'new character' key in list state")
	}

	// Test form state - should expose form keys
	page.characterForm = newCharacterForm("", nil, 80, 24)
	formKeys := page.ShortHelp()

	if len(formKeys) == 0 {
		t.Error("Expected form keys when character form is active")
	}

	// Should have escape key from form
	hasEscKey := false
	for _, key := range formKeys {
		if key.Help().Key == "esc" {
			hasEscKey = true
			break
		}
	}
	if !hasEscKey {
		t.Error("Expected form to expose escape key")
	}
}

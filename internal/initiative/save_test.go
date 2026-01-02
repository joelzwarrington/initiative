package initiative

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
)

func TestLoadGame(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		includeSRD  bool
		wantErr     bool
		wantChars   int
		setupFile   bool
	}{
		{
			name:       "load from non-existent file",
			includeSRD: false,
			wantErr:    false,
			wantChars:  0,
			setupFile:  false,
		},
		{
			name: "load valid game with characters",
			fileContent: `characters:
  char1:
    name: "Fighter"
  char2:
    name: "Wizard"`,
			includeSRD: false,
			wantErr:    false,
			wantChars:  2,
			setupFile:  true,
		},
		{
			name:        "load corrupted YAML",
			fileContent: "invalid: yaml: content: [",
			includeSRD:  false,
			wantErr:     true,
			wantChars:   0,
			setupFile:   true,
		},
		{
			name:       "load with SRD inclusion",
			includeSRD: true,
			wantErr:    false,
			wantChars:  0,
			setupFile:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "test_game.yaml")

			if tt.setupFile {
				err := os.WriteFile(filepath, []byte(tt.fileContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			game, err := LoadGame(filepath, []string{})

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if game == nil {
				t.Error("Expected game but got nil")
				return
			}

			if len(game.Characters) != tt.wantChars {
				t.Errorf("Expected %d characters, got %d", tt.wantChars, len(game.Characters))
			}

			if game.filepath != filepath {
				t.Errorf("Expected filepath %s, got %s", filepath, game.filepath)
			}
		})
	}
}

func TestGameSave(t *testing.T) {
	tests := []struct {
		name       string
		characters map[string]dnd.Character
		wantErr    bool
	}{
		{
			name:       "save empty game",
			characters: map[string]dnd.Character{},
			wantErr:    false,
		},
		{
			name: "save game with characters",
			characters: map[string]dnd.Character{
				"char1": dnd.NewCharacter("Fighter"),
				"char2": dnd.NewCharacter("Wizard"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "test_game.yaml")

			game := &Game{
				filepath:   filepath,
				Characters: tt.characters,
				sources:    make(map[string]*dnd.Source),
			}

			err := game.Save()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Error("Expected file to be created")
			}
		})
	}
}

func TestGameSaveInvalidPath(t *testing.T) {
	game := &Game{
		filepath:   "/invalid/path/that/does/not/exist/game.yaml",
		Characters: make(map[string]dnd.Character),
		sources:    make(map[string]*dnd.Source),
	}

	err := game.Save()
	if err == nil {
		t.Error("Expected error for invalid path but got none")
	}
}

func TestGameRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		characters map[string]dnd.Character
	}{
		{
			name:       "round trip empty game",
			characters: map[string]dnd.Character{},
		},
		{
			name: "round trip with characters",
			characters: map[string]dnd.Character{
				"fighter1": dnd.NewCharacter("Aragorn"),
				"wizard1":  dnd.NewCharacter("Gandalf"),
				"rogue1":   dnd.NewCharacter("Legolas"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filepath := filepath.Join(tmpDir, "roundtrip_game.yaml")

			originalGame := &Game{
				filepath:   filepath,
				Characters: tt.characters,
				sources:    make(map[string]*dnd.Source),
			}

			err := originalGame.Save()
			if err != nil {
				t.Fatalf("Failed to save game: %v", err)
			}

			loadedGame, err := LoadGame(filepath, []string{})
			if err != nil {
				t.Fatalf("Failed to load game: %v", err)
			}

			if len(loadedGame.Characters) != len(originalGame.Characters) {
				t.Errorf("Character count mismatch: expected %d, got %d",
					len(originalGame.Characters), len(loadedGame.Characters))
			}

			for uuid, originalChar := range originalGame.Characters {
				loadedChar, exists := loadedGame.Characters[uuid]
				if !exists {
					t.Errorf("Character %s not found in loaded game", uuid)
					continue
				}
				if loadedChar.Name() != originalChar.Name() {
					t.Errorf("Character name mismatch for %s: expected %s, got %s",
						uuid, originalChar.Name(), loadedChar.Name())
				}
			}
		})
	}
}

package dnd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSourceFromFS(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantError   bool
		errorString string
	}{
		{
			name:        "non-existent file",
			path:        "nonexistent.yaml",
			wantError:   true,
			errorString: "failed to read source file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadSourceFromFS(tt.path)

			if tt.wantError {
				if err == nil {
					t.Errorf("LoadSourceFromFS() error = nil, wantError %v", tt.wantError)
				}
				if !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("LoadSourceFromFS() error = %v, want to contain %v", err, tt.errorString)
				}
			} else {
				if err != nil {
					t.Errorf("LoadSourceFromFS() error = %v, wantError %v", err, tt.wantError)
				}
			}
		})
	}
}

func TestLoadSourceFromFile(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("non-existent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "nonexistent.yaml")
		_, err := LoadSourceFromFile(filePath)
		if err == nil {
			t.Errorf("LoadSourceFromFile() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("LoadSourceFromFile() error = %v, want to contain 'not found'", err)
		}
	})

	t.Run("valid yaml file", func(t *testing.T) {
		validYaml := `meta:
  name: "Test Source"
  key: "test"
  version: "1.0"
monsters:
  - name: "Test Monster"`

		filePath := filepath.Join(tempDir, "valid.yaml")
		err := os.WriteFile(filePath, []byte(validYaml), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		source, err := LoadSourceFromFile(filePath)
		if err != nil {
			t.Errorf("LoadSourceFromFile() error = %v, want nil", err)
		}
		if source == nil {
			t.Errorf("LoadSourceFromFile() returned nil source")
		}
		if source.Meta.Name != "Test Source" {
			t.Errorf("Expected meta name 'Test Source', got %v", source.Meta.Name)
		}
		if len(source.Monsters) != 1 {
			t.Errorf("Expected 1 monster, got %v", len(source.Monsters))
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		invalidYaml := "invalid: yaml: content: [unclosed"
		filePath := filepath.Join(tempDir, "invalid.yaml")
		err := os.WriteFile(filePath, []byte(invalidYaml), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = LoadSourceFromFile(filePath)
		if err == nil {
			t.Errorf("LoadSourceFromFile() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "failed to unmarshal") {
			t.Errorf("LoadSourceFromFile() error = %v, want to contain 'failed to unmarshal'", err)
		}
	})
}

func TestLoadSource(t *testing.T) {
	t.Run("known embedded source key", func(t *testing.T) {
		source, err := LoadSource("srd")
		if err != nil {
			t.Errorf("LoadSource() error = %v, want nil", err)
		}
		if source == nil {
			t.Errorf("LoadSource() returned nil source")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadSource("/path/to/nonexistent.yaml")
		if err == nil {
			t.Errorf("LoadSource() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("LoadSource() error = %v, want to contain 'not found'", err)
		}
	})
}

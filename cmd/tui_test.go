package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildTUIConfig(t *testing.T) {
	tests := []struct {
		name        string
		gameFile    string
		sources     []string
		expectSRD   bool
		expectError bool
	}{
		{
			name:      "defaults to srd source",
			gameFile:  "/tmp/test.yaml",
			sources:   []string{},
			expectSRD: true,
		},
		{
			name:      "uses specified sources",
			gameFile:  "/tmp/test.yaml",
			sources:   []string{"custom", "homebrew"},
			expectSRD: false,
		},
		{
			name:      "creates default game file path",
			gameFile:  "",
			sources:   []string{"srd"},
			expectSRD: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up command with flags
			cmd := &cobra.Command{}
			cmd.Flags().String("game", tt.gameFile, "")
			cmd.Flags().StringSlice("source", tt.sources, "")

			cfg, err := buildTUIConfig(cmd)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check sources
			if tt.expectSRD && !containsSource(cfg.Sources, "srd") {
				t.Error("expected SRD source but not found")
			}

			if !tt.expectSRD && containsSource(cfg.Sources, "srd") && len(tt.sources) > 0 {
				t.Error("did not expect SRD source but found it")
			}

			// Check game file
			if tt.gameFile != "" && cfg.GameFile != tt.gameFile {
				t.Errorf("expected game file %s, got %s", tt.gameFile, cfg.GameFile)
			}

			if tt.gameFile == "" && !strings.Contains(cfg.GameFile, "game.yaml") {
				t.Error("expected default game file path to contain 'game.yaml'")
			}
		})
	}
}

func TestTUIConfigDefaults(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("game", "", "")
	cmd.Flags().StringSlice("source", []string{}, "")

	cfg, err := buildTUIConfig(cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should default to SRD
	if !containsSource(cfg.Sources, "srd") {
		t.Error("expected default sources to include SRD")
	}

	// Should create a default game file in config dir
	if !strings.Contains(cfg.GameFile, ".config/initiative/game.yaml") {
		t.Errorf("expected default config path, got: %s", cfg.GameFile)
	}
}

func TestTUIConfigWithCustomPath(t *testing.T) {
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom.yaml")

	cmd := &cobra.Command{}
	cmd.Flags().String("game", customPath, "")
	cmd.Flags().StringSlice("source", []string{"custom-source"}, "")

	cfg, err := buildTUIConfig(cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.GameFile != customPath {
		t.Errorf("expected custom path %s, got %s", customPath, cfg.GameFile)
	}

	if !containsSource(cfg.Sources, "custom-source") {
		t.Error("expected custom source")
	}
}

// Test that configuration directory creation works
func TestTUIConfigDirectoryCreation(t *testing.T) {
	// Test with a temporary directory to ensure we can create dirs
	tempDir := t.TempDir()

	// Temporarily change HOME to our temp dir
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cmd := &cobra.Command{}
	cmd.Flags().String("game", "", "")
	cmd.Flags().StringSlice("source", []string{}, "")

	cfg, err := buildTUIConfig(cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the config directory was created
	expectedConfigDir := filepath.Join(tempDir, ".config", "initiative")
	if _, err := os.Stat(expectedConfigDir); os.IsNotExist(err) {
		t.Error("expected config directory to be created")
	}

	// Check that the game file path is correct
	expectedGameFile := filepath.Join(expectedConfigDir, "game.yaml")
	if cfg.GameFile != expectedGameFile {
		t.Errorf("expected game file %s, got %s", expectedGameFile, cfg.GameFile)
	}
}

// Helper function is defined in root_test.go

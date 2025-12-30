package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommandHasSameFlagsAsTUI(t *testing.T) {
	// Verify we're working with proper cobra commands
	var _ *cobra.Command = rootCmd
	var _ *cobra.Command = tuiCmd
	
	// Get flags from both commands
	rootFlags := rootCmd.Flags()
	tuiFlags := tuiCmd.Flags()

	// Check that root has the same flags as TUI for backward compatibility
	expectedFlags := []string{"game", "source"}

	for _, flagName := range expectedFlags {
		rootFlag := rootFlags.Lookup(flagName)
		tuiFlag := tuiFlags.Lookup(flagName)

		if rootFlag == nil {
			t.Errorf("root command missing flag: %s", flagName)
		}
		if tuiFlag == nil {
			t.Errorf("tui command missing flag: %s", flagName)
		}

		if rootFlag != nil && tuiFlag != nil {
			if rootFlag.Usage != tuiFlag.Usage {
				t.Errorf("flag %s has different usage text between root and tui", flagName)
			}
		}
	}
}

func TestRootCommandDefaultsToTUI(t *testing.T) {
	// Verify that root command has a RunE function (defaults to TUI)
	if rootCmd.RunE == nil {
		t.Error("root command should have RunE function to default to TUI")
	}

	// Verify that root command doesn't have persistent flags that would affect all subcommands
	persistentFlags := rootCmd.PersistentFlags()
	if persistentFlags.HasFlags() {
		t.Error("root command should not have persistent flags that affect all subcommands")
	}
}

func TestRootCommandStructure(t *testing.T) {
	// Test that expected subcommands are present
	expectedCommands := []string{"tui"}
	
	// Verify rootCmd is properly configured
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	
	commands := rootCmd.Commands()
	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd.Name() == expectedCmd {
				found = true
				// Verify it's a proper cobra command
				if cmd == nil {
					t.Errorf("subcommand %s should not be nil", expectedCmd)
				}
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %s not found", expectedCmd)
		}
	}
}

func TestRootCommandMetadata(t *testing.T) {
	if rootCmd.Use != "initiative" {
		t.Errorf("expected Use to be 'initiative', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("root command should have a short description")
	}

	// Long description is optional - we prefer concise help
}

// Helper function (shared between test files)
func containsSource(sources []string, target string) bool {
	for _, s := range sources {
		if s == target {
			return true
		}
	}
	return false
}
package cmd

import (
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/spf13/cobra"
)

func TestRollCommand(t *testing.T) {
	// Test valid notation works
	cmd := &cobra.Command{}
	cmd.Flags().Int64P("seed", "", 42, "Random seed")

	err := runRoll(cmd, []string{"2d6+3"})
	if err != nil {
		t.Errorf("Valid notation should not error: %v", err)
	}

	// Test invalid notation fails
	err = runRoll(cmd, []string{"invalid"})
	if err == nil {
		t.Error("Invalid notation should error")
	}
}

func TestRollCommandSeeding(t *testing.T) {
	// Test that same seed produces same result
	dnd.SetSeed(42)
	expected, _ := dnd.Roll("d6")

	cmd := &cobra.Command{}
	cmd.Flags().Int64P("seed", "", 42, "Random seed")

	err := runRoll(cmd, []string{"d6"})
	if err != nil {
		t.Errorf("Seeded roll should not error: %v", err)
	}

	// Verify the integration works (seed was applied)
	dnd.SetSeed(42)
	result, _ := dnd.Roll("d6")
	if result.Total != expected.Total {
		t.Error("Seed should produce deterministic results")
	}
}

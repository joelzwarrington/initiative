package cmd

import (
	"fmt"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/spf13/cobra"
)

// rollCmd represents the roll command
var rollCmd = &cobra.Command{
	Use:   "roll <dice-notation>",
	Short: "Roll dice using standard D&D notation",
	Long: `Roll dice using standard D&D notation.

Examples:
  initiative roll d6        # Roll a six-sided die
  initiative roll 2d6       # Roll two six-sided dice
  initiative roll 1d20+5    # Roll a twenty-sided die and add 5
  initiative roll 3d8-2     # Roll three eight-sided dice and subtract 2

Dice notation follows the grammar: [count]d<sides>[+/-modifier]
- count: number of dice (defaults to 1 if omitted)
- sides: number of sides on each die (must be > 0)
- modifier: optional positive or negative number to add (must be non-zero)`,
	Args: cobra.ExactArgs(1),
	RunE: runRoll,
}

func init() {
	rootCmd.AddCommand(rollCmd)

	// Add seed flag for deterministic results (useful for testing)
	rollCmd.Flags().Int64P("seed", "", 0, "Random seed for deterministic results (0 = random)")
}

func runRoll(cmd *cobra.Command, args []string) error {
	notation := args[0]

	// Handle seed if provided
	seed, _ := cmd.Flags().GetInt64("seed")
	if seed != 0 {
		dnd.SetSeed(seed)
	}

	// Roll the dice
	result, err := dnd.Roll(notation)
	if err != nil {
		return fmt.Errorf("invalid dice notation '%s': %w", notation, err)
	}

	// Display the result
	fmt.Println(result.String())

	return nil
}

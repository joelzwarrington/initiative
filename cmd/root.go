package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "initiative",
	Short: "CLI & TUI to run D&D games",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to TUI command when no subcommand is specified
		return runTUI(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Disable automatic error and usage display
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	err := rootCmd.Execute()
	if err != nil {
		red := "\033[31m"
		reset := "\033[0m"

		fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", red, err.Error(), reset)
		os.Exit(1)
	}
}

func init() {
	// Add the same flags for tui subcommand.
	rootCmd.Flags().StringP("game", "g", "", "Path to game file (default: ~/.config/initiative/game.yaml)")
	rootCmd.Flags().StringSliceP("source", "s", []string{}, "Source to use, can be specified multiple times. Use 'srd' to include the system reference document, or specify no source to include automatically (default: srd)")
}

package main

import (
	"fmt"
	"initiative/internal/initiative"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var saveFilename string
var withSRD bool

var rootCmd = &cobra.Command{
	Use:   "initiative",
	Short: "A CLI tool for managing tabletop RPG initiative tracking",
	Run: func(cmd *cobra.Command, args []string) {
		if saveFilename == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Error getting home directory: %v\n", err)
				os.Exit(1)
			}
			configDir := filepath.Join(homeDir, ".config", "initiative")
			if err := os.MkdirAll(configDir, 0755); err != nil {
				fmt.Printf("Error creating config directory: %v\n", err)
				os.Exit(1)
			}
			saveFilename = filepath.Join(configDir, "game.yaml")
		}

		noSrd, _ := cmd.Flags().GetBool("no-srd")
		if noSrd {
			withSRD = false
		}

		p := initiative.NewProgram(saveFilename, withSRD)

		if _, err := p.Run(); err != nil {
			panic(err)
		}

		// Save on normal exit too
		initiative.SaveGame()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&saveFilename, "game", "g", "", "Path to game file (default: ~/.config/initiative/game.yaml)")
	rootCmd.Flags().BoolVar(&withSRD, "srd", true, "Include System Reference Document monsters")
	rootCmd.Flags().BoolVar(&withSRD, "no-srd", false, "Exclude System Reference Document monsters")
	rootCmd.MarkFlagsMutuallyExclusive("srd", "no-srd")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

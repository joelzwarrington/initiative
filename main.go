package main

import (
	"fmt"
	"initiative/internal/data"
	"initiative/internal/ui"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dataFile string

var rootCmd = &cobra.Command{
	Use:   "initiative",
	Short: "A CLI tool for managing tabletop RPG initiative tracking",
	Run: func(cmd *cobra.Command, args []string) {
		// Set default data file if not provided
		if dataFile == "" {
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
			dataFile = filepath.Join(configDir, "data.yaml")
		}

		// Load data from file
		appData, err := data.Load(dataFile)
		if err != nil {
			fmt.Printf("Error loading data: %v\n", err)
			os.Exit(1)
		}

		p := ui.NewProgram(appData)

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&dataFile, "data", "d", "", "path to data file (default: ~/.config/initiative/data.yaml)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

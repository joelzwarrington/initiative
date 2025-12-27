package main

import (
	"fmt"
	"initiative/internal/nui"
	"os"

	"github.com/spf13/cobra"
)

var dataFile string

var rootCmd = &cobra.Command{
	Use:   "initiative",
	Short: "A CLI tool for managing tabletop RPG initiative tracking",
	Run: func(cmd *cobra.Command, args []string) {
		p := nui.NewProgram()

		if _, err := p.Run(); err != nil {
			panic(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

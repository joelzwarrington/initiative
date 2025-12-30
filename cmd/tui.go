package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joelzwarrington/initiative/internal/initiative"
	"github.com/spf13/cobra"
)

// TUIConfig holds the configuration for the TUI command
type TUIConfig struct {
	GameFile string
	Sources  []string
}

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the TUI interface for managing D&D encounters",
	RunE:  runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)

	// TUI-specific flags
	tuiCmd.Flags().StringP("game", "g", "", "Path to game file (default: ~/.config/initiative/game.yaml)")
	tuiCmd.Flags().StringSliceP("source", "s", []string{}, "Source to use, can be specified multiple times. Use 'srd' to include the system reference document, or specify no source to include automatically (default: srd)")
}

func runTUI(cmd *cobra.Command, args []string) error {
	cfg, err := buildTUIConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to build configuration: %w", err)
	}

	return runTUIWithConfig(cfg)
}

func buildTUIConfig(cmd *cobra.Command) (*TUIConfig, error) {
	// Read from command flags (inherited from root persistent flags)
	configGameFile, _ := cmd.Flags().GetString("game")
	configSources, _ := cmd.Flags().GetStringSlice("source")

	// Set default game file if none specified
	if configGameFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(homeDir, ".config", "initiative")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		configGameFile = filepath.Join(configDir, "game.yaml")
	}

	// Default to srd if no sources specified
	if len(configSources) == 0 {
		configSources = []string{"srd"}
	}

	return &TUIConfig{
		GameFile: configGameFile,
		Sources:  configSources,
	}, nil
}

func runTUIWithConfig(cfg *TUIConfig) error {
	p := initiative.NewProgram(cfg.GameFile, cfg.Sources)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	initiative.SaveGame()
	return nil
}

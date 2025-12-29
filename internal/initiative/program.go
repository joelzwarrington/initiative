package initiative

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

var gameInstance *Game

func NewProgram(filepath string, withSRD bool) *tea.Program {
	game, err := LoadGame(filepath, withSRD)
	if err != nil {
		panic(err)
	}

	p := &game.Characters

	s := skeleton.NewSkeleton()

	s.SetPagePosition(lipgloss.Left)

	s.KeyMap.SwitchTabRight = key.NewBinding(
		key.WithKeys("tab"))

	// To switch previous page
	s.KeyMap.SwitchTabLeft = key.NewBinding(
		key.WithKeys("shift+tab"))

	s.LockTabs().SetWrapTabs(true)

	addPage(s, newEncounterPage(s, p, game.sources))
	addPage(s, newCharacterPage(s, p))

	// Store game reference for saving on exit
	gameInstance = game

	// Set up signal handler for saving on exit
	setupSignalHandler()

	prog := tea.NewProgram(s)
	return prog
}

type page interface {
	tea.Model

	Key() string
	Title() string
}

func addPage(s *skeleton.Skeleton, p page) {
	s.AddPage(p.Key(), p.Title(), p)
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		SaveGame()
		os.Exit(0)
	}()
}

func SaveGame() {
	if gameInstance != nil {
		if err := gameInstance.Save(); err != nil {
			// We can't really handle errors here since we're exiting
			// Maybe log to stderr in the future
		}
	}
}

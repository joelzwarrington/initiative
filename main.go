package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type app struct{}

func newApp() app {
	return app{}
}

func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a app) View() string {
	return "Hello World!"
}

func main() {
	app := newApp()

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

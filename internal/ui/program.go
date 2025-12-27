package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/termkit/skeleton"
)

func NewProgram() *tea.Program {
	party := map[string]Character{
		uuid.New().String(): Character{name: "Lorem"},
		uuid.New().String(): Character{name: "Ipsum"},
	}
	p := &party

	s := skeleton.NewSkeleton()

	s.SetPagePosition(lipgloss.Left)

	s.KeyMap.SwitchTabRight = key.NewBinding(
		key.WithKeys("tab"))

	// To switch previous page
	s.KeyMap.SwitchTabLeft = key.NewBinding(
		key.WithKeys("shift+tab"))

	s.LockTabs().SetWrapTabs(true)

	s.AddPage("encounter", "Encounter", newEncounter(s, p))
	s.AddPage("party", "Party", newParty(s, p))

	return tea.NewProgram(s)
}

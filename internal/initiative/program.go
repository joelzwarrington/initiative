package initiative

import (
	"initiative/dnd"
	"initiative/internal/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

func NewProgram() *tea.Program {
	party := map[string]dnd.Character{}
	p := &party

	// Load SRD sources at program initialization
	sources := make(map[string]*dnd.Source)
	srd, err := dnd.GetSystemReferenceSource()
	if err == nil && srd != nil {
		sources[srd.Meta.Key] = srd
	}

	s := skeleton.NewSkeleton()

	s.SetPagePosition(lipgloss.Left)

	s.KeyMap.SwitchTabRight = key.NewBinding(
		key.WithKeys("tab"))

	// To switch previous page
	s.KeyMap.SwitchTabLeft = key.NewBinding(
		key.WithKeys("shift+tab"))

	s.LockTabs().SetWrapTabs(true)

	addPage(s, newEncounterPage(s, p, sources))
	s.AddPage("party", "Party", ui.NewParty(s, p))

	return tea.NewProgram(s)
}

type page interface {
	tea.Model

	Key() string
	Title() string
}

func addPage(s *skeleton.Skeleton, p page) {
	s.AddPage(p.Key(), p.Title(), p)
}

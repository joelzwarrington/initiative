package initiative

import (
	"initiative/dnd"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

func NewProgram() *tea.Program {
	characters := map[string]dnd.Character{}
	p := &characters

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
	addPage(s, newCharacterPage(s, p))

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

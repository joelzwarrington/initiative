package nui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/termkit/skeleton"
)

var _ tea.Model = (*encounter)(nil)

type encounter struct {
	skeleton *skeleton.Skeleton
	party    *map[string]Character
}

func newEncounter(skeleton *skeleton.Skeleton, party *map[string]Character) *encounter {
	return &encounter{
		skeleton: skeleton,
		party:    party,
	}
}

func (e encounter) Init() tea.Cmd {
	return nil
}

func (e encounter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e, nil
}

func (e encounter) View() string {
	return ""
}

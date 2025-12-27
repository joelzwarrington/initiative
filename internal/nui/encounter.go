package nui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/termkit/skeleton"
)

var _ tea.Model = (*encounter)(nil)

type encounterView int

const (
	encounterForm encounterView = iota
	encounterDetail
)

type encounter struct {
	Encounter

	skeleton *skeleton.Skeleton
	party    *map[string]Character

	view encounterView
	form *huh.Form
}

func newEncounter(skeleton *skeleton.Skeleton, party *map[string]Character) *encounter {
	var characterOptions []huh.Option[string]

	if party != nil {
		for uuid, character := range *party {
			characterOptions = append(characterOptions,
				huh.NewOption(character.Name(), uuid).Selected(true),
			)
		}
	}

	return &encounter{
		skeleton: skeleton,
		party:    party,

		view: encounterForm,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("summary").
					Title("Summary"),
				huh.NewMultiSelect[string]().
					Key("characters").
					Title("Characters").
					Options(characterOptions...),
			),
		),
	}
}

func (e encounter) Init() tea.Cmd {
	return e.form.Init()
}

func (e encounter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch e.view {
	case encounterForm:
		{
			form, cmd := e.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				e.form = f
			}

			if e.form.State == huh.StateCompleted {
				e.Summary = e.form.GetString("summary")
				e.StartedAt = time.Now()

				// Create initiative groups for each selected character
				selectedCharacterUUIDs := e.form.Get("characters").([]string)
				e.IniativeGroups = []IniativeGroup{}
				
				if e.party != nil {
					for _, uuid := range selectedCharacterUUIDs {
						if character, exists := (*e.party)[uuid]; exists {
							group := IniativeGroup{
								Iniative:  0, // Initiative will be set later
								Creatures: []Creature{character},
							}
							e.IniativeGroups = append(e.IniativeGroups, group)
						}
					}
				}

				e.form = nil
				e.view = encounterDetail
			}

			return e, cmd
		}
	case encounterDetail:
		{
			return e, nil
		}
	}
	return e, nil
}

func (e encounter) View() string {
	switch e.view {
	case encounterForm:
		{
			e.form.WithHeight(e.skeleton.GetContentHeight()).WithWidth(e.skeleton.GetContentWidth())
			return e.form.View()
		}
	case encounterDetail:
		{
			var content strings.Builder
			
			// Encounter summary header
			headerStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)
			content.WriteString(headerStyle.Render(fmt.Sprintf("Encounter: %s", e.Summary)))
			content.WriteString("\n\n")
			
			// Initiative order title
			titleStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("39")).
				MarginBottom(1)
			content.WriteString(titleStyle.Render("Initiative Order:"))
			content.WriteString("\n")
			
			if len(e.IniativeGroups) == 0 {
				noGroupsStyle := lipgloss.NewStyle().
					Italic(true).
					Foreground(lipgloss.Color("240"))
				content.WriteString(noGroupsStyle.Render("  No participants"))
			} else {
				// Create tree with styled nodes
				initiativeTree := tree.New().
					EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
					ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("252")))
				
				// Add each initiative group to the tree
				for _, group := range e.IniativeGroups {
					// Initiative value (0 if not set)
					initiativeText := "Initiative: TBD"
					if group.Iniative > 0 {
						initiativeText = fmt.Sprintf("Initiative: %d", group.Iniative)
					}
					
					// Style the initiative group
					initiativeStyle := lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("214"))
					styledInitiative := initiativeStyle.Render(initiativeText)
					
					// Create a child tree for creatures in this group
					if len(group.Creatures) == 1 {
						// Single creature - add directly
						initiativeTree.Child(styledInitiative + " - " + group.Creatures[0].Name())
					} else {
						// Multiple creatures - create subtree
						groupTree := tree.New().
							Root(styledInitiative).
							EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
							ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("252")))
						
						for _, creature := range group.Creatures {
							groupTree.Child(creature.Name())
						}
						
						initiativeTree.Child(groupTree)
					}
				}
				
				content.WriteString(initiativeTree.String())
			}
			
			return content.String()
		}
	}

	return ""
}

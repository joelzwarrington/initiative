package initiative

import (
	"fmt"
	"initiative/dnd"
	"initiative/internal/components"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

type encounterPage struct {
	s      *skeleton.Skeleton
	width  int
	height int
	keys   EncounterPageKeyMap

	characters *map[string]dnd.Character
	sources    map[string]*dnd.Source

	encounter *dnd.Encounter

	emptyState        *components.EmptyState
	encounterForm     *encounterForm
	encounterDelegate *encounterDelegate
	help              help.Model
}

func newEncounterPage(s *skeleton.Skeleton, characters *map[string]dnd.Character, sources map[string]*dnd.Source) *encounterPage {
	keys := defaultEncounterPageKeyMap()

	page := &encounterPage{
		s: s, width: s.GetContentWidth(), height: s.GetContentHeight(),
		keys: keys,

		characters: characters,
		sources:    sources,
	}

	page.emptyState = components.NewEmptyState(
		"No encounter",
		page,
		s.GetContentWidth(), s.GetContentHeight(),
	)

	page.encounterDelegate = newEncounterDelegate(s.GetContentWidth(), s.GetContentHeight())
	page.help = help.New()

	return page
}

func (p *encounterPage) Init() tea.Cmd {
	return nil
}

func (p *encounterPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = p.s.GetContentWidth()
		p.height = p.s.GetContentHeight()
		if p.encounterForm != nil {
			p.encounterForm.SetSize(p.width, p.height)
		}
		if p.encounterDelegate != nil {
			p.encounterDelegate.SetSize(p.width, p.height)
		}
		if p.emptyState != nil {
			p.emptyState.SetSize(p.width, p.height)
		}

		return p, nil
	case tea.KeyMsg:
		if key.Matches(msg, p.keys.NewEncounter) && p.isEmptyState() {
			return p, tea.Sequence(append(cmds, p.beginNewEncounterForm())...)
		}
		if key.Matches(msg, p.keys.NextTurn) && p.encounter != nil {
			p.encounter.AdvanceTurn()
			p.updateRoundWidget()
			return p, nil
		}
		if key.Matches(msg, p.keys.PrevTurn) && p.encounter != nil {
			p.encounter.PreviousTurn()
			p.updateRoundWidget()
			return p, nil
		}
	case encounterFormCancelledMsg:
		return p, tea.Sequence(append(cmds, p.cancelNewEncounterForm())...)
	case encounterFormSubmittedMsg:
		return p, tea.Sequence(append(cmds, p.addNewEncounter(msg))...)
	}

	if p.isAddingNewEncounter() {
		form, cmd := p.encounterForm.Update(msg)
		cmds = append(cmds, cmd)
		if form, ok := form.(*encounterForm); ok {
			p.encounterForm = form
		}
	}

	if p.encounter != nil && p.encounterDelegate != nil {
		cmds = append(cmds, p.encounterDelegate.Update(msg, p.encounter))
	}

	return p, tea.Batch(cmds...)
}

func (p *encounterPage) View() string {
	switch true {
	case p.encounter != nil:
		// Calculate available height for list (subtract help height)
		p.help.Width = p.s.GetContentWidth()
		helpStyle := lipgloss.NewStyle().Padding(0, 1)
		helpView := helpStyle.Render(p.help.View(p))
		helpHeight := lipgloss.Height(helpView)
		
		listHeight := p.s.GetContentHeight() - helpHeight
		p.encounterDelegate.SetSize(p.s.GetContentWidth(), listHeight)
		
		var buf strings.Builder
		p.encounterDelegate.Render(&buf, p.encounter)
		listView := buf.String()
		
		return lipgloss.JoinVertical(lipgloss.Left, listView, helpView)
	case p.isAddingNewEncounter():
		return p.encounterForm.View()
	case p.isEmptyState():
		return p.emptyState.View()
	default:
		// This shouldn't be possible, but we return to satisfy the types.
		return "No view"
	}
}

func (p *encounterPage) Key() string {
	return "encounter"
}

func (p *encounterPage) Title() string {
	return "Encounter"
}

func (p *encounterPage) FullHelp() [][]key.Binding {
	if p.encounter != nil {
		return [][]key.Binding{{
			p.keys.PrevTurn,
			p.keys.NextTurn,
		}}
	}
	return [][]key.Binding{{
		p.keys.NewEncounter,
	}}
}

func (p *encounterPage) ShortHelp() []key.Binding {
	if p.encounter != nil {
		return []key.Binding{
			p.keys.PrevTurn,
			p.keys.NextTurn,
		}
	}
	return []key.Binding{
		p.keys.NewEncounter,
	}
}

type EncounterPageKeyMap struct {
	NewEncounter key.Binding
	NextTurn     key.Binding
	PrevTurn     key.Binding
}

func defaultEncounterPageKeyMap() EncounterPageKeyMap {
	return EncounterPageKeyMap{
		NewEncounter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new encounter"),
		),
		NextTurn: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next turn"),
		),
		PrevTurn: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "previous turn"),
		),
	}
}

func (p *encounterPage) isAddingNewEncounter() bool {
	return p.encounterForm != nil
}

func (p *encounterPage) isEmptyState() bool {
	return !p.isAddingNewEncounter() && p.encounter == nil
}

func (p *encounterPage) beginNewEncounterForm() tea.Cmd {
	p.s.LockTabs()
	p.encounterForm = newEncounterForm(p.characters, p.sources, p.s.GetContentWidth(), p.s.GetContentHeight())
	return p.encounterForm.Init()
}

func (p *encounterPage) cancelNewEncounterForm() tea.Cmd {
	p.encounterForm = nil
	p.s.UnlockTabs()
	return nil
}

func (p *encounterPage) updateRoundWidget() {
	if p.encounter != nil {
		p.s.UpdateWidgetValue("round", fmt.Sprintf("Round: %d", p.encounter.Round))
	}
}

func (p *encounterPage) addNewEncounter(submission encounterFormSubmittedMsg) tea.Cmd {
	p.encounterForm = nil
	p.encounter = &submission.encounter
	
	// Add round tracking widget and unlock tabs
	p.s.AddWidget("round", fmt.Sprintf("Round: %d", p.encounter.Round))
	p.s.UnlockTabs()
	
	return nil
}

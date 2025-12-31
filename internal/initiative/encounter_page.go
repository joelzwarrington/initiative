package initiative

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
	"github.com/joelzwarrington/initiative/internal/components"
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
	hitPointForm      *hitPointForm
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
		if p.hitPointForm != nil {
			p.hitPointForm.SetSize(p.width, p.height)
		}
		if p.encounterDelegate != nil {
			p.encounterDelegate.SetSize(p.width, p.height)
		}
		if p.emptyState != nil {
			p.emptyState.SetSize(p.width, p.height)
		}

		return p, nil
	case tea.KeyMsg:
		// Don't process page-level keys when forms are active
		if p.isAddingNewEncounter() || p.isAdjustingHitPoints() {
			break
		}

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
		if key.Matches(msg, p.keys.EndEncounter) && p.encounter != nil {
			return p, tea.Sequence(append(cmds, p.endEncounter())...)
		}
	case encounterFormCancelledMsg:
		return p, tea.Sequence(append(cmds, p.cancelNewEncounterForm())...)
	case encounterFormSubmittedMsg:
		return p, tea.Sequence(append(cmds, p.addNewEncounter(msg))...)
	case adjustHitPointsMsg:
		return p, tea.Sequence(append(cmds, p.beginHitPointForm(msg))...)
	case hitPointFormCancelledMsg:
		return p, tea.Sequence(append(cmds, p.cancelHitPointForm())...)
	case hitPointFormSubmittedMsg:
		return p, tea.Sequence(append(cmds, p.processHitPointAdjustment(msg))...)
	}

	if p.isAddingNewEncounter() {
		form, cmd := p.encounterForm.Update(msg)
		cmds = append(cmds, cmd)
		if form, ok := form.(*encounterForm); ok {
			p.encounterForm = form
		}
	}

	if p.isAdjustingHitPoints() {
		form, cmd := p.hitPointForm.Update(msg)
		cmds = append(cmds, cmd)
		if form, ok := form.(*hitPointForm); ok {
			p.hitPointForm = form
		}
	}

	if p.encounter != nil && p.encounterDelegate != nil && !p.isAdjustingHitPoints() {
		cmds = append(cmds, p.encounterDelegate.Update(msg, p.encounter))
	}

	return p, tea.Batch(cmds...)
}

func (p *encounterPage) View() string {
	switch true {
	case p.isAdjustingHitPoints():
		return p.hitPointForm.View()
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
	if p.encounter != nil {
		summary := p.encounter.Summary
		if len(summary) > 18 {
			summary = summary[:15] + "..."
		}
		return "Encounter > " + summary
	}
	return "Encounter"
}

func (p *encounterPage) FullHelp() [][]key.Binding {
	if p.encounter != nil {
		keys := [][]key.Binding{{
			p.keys.PrevTurn,
			p.keys.NextTurn,
			p.keys.EndEncounter,
		}}
		// Add list navigation keys
		if p.encounterDelegate != nil {
			listKeys := p.encounterDelegate.FullHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	}
	return [][]key.Binding{{
		p.keys.NewEncounter,
	}}
}

func (p *encounterPage) ShortHelp() []key.Binding {
	if p.encounter != nil {
		keys := []key.Binding{
			p.keys.PrevTurn,
			p.keys.NextTurn,
			p.keys.EndEncounter,
		}
		// Add list navigation keys
		if p.encounterDelegate != nil {
			listKeys := p.encounterDelegate.ShortHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	}
	return []key.Binding{
		p.keys.NewEncounter,
	}
}

type EncounterPageKeyMap struct {
	NewEncounter key.Binding
	NextTurn     key.Binding
	PrevTurn     key.Binding
	EndEncounter key.Binding
}

func defaultEncounterPageKeyMap() EncounterPageKeyMap {
	return EncounterPageKeyMap{
		NewEncounter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new"),
		),
		NextTurn: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next turn"),
		),
		PrevTurn: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "previous turn"),
		),
		EndEncounter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "end"),
		),
	}
}

func (p *encounterPage) isAddingNewEncounter() bool {
	return p.encounterForm != nil
}

func (p *encounterPage) isAdjustingHitPoints() bool {
	return p.hitPointForm != nil
}

func (p *encounterPage) isEmptyState() bool {
	return !p.isAddingNewEncounter() && !p.isAdjustingHitPoints() && p.encounter == nil
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

	// Update page title with encounter summary
	summary := p.encounter.Summary
	if len(summary) > 18 {
		summary = summary[:15] + "..."
	}
	p.s.UpdatePageTitle("encounter", "Encounter > "+summary)

	// Add round tracking widget and unlock tabs
	p.s.AddWidget("round", fmt.Sprintf("Round: %d", p.encounter.Round))
	p.s.UnlockTabs()

	return nil
}

func (p *encounterPage) endEncounter() tea.Cmd {
	p.encounter = nil

	// Reset page title to default
	p.s.UpdatePageTitle("encounter", "Encounter")

	// Clear round tracking widget
	p.s.DeleteWidget("round")

	return nil
}

func (p *encounterPage) beginHitPointForm(msg adjustHitPointsMsg) tea.Cmd {
	if p.encounter == nil || len(p.encounter.InitiativeGroups) <= msg.groupIndex {
		return nil
	}

	group := p.encounter.InitiativeGroups[msg.groupIndex]
	if len(group.Creatures) <= msg.creatureIndex {
		return nil
	}

	creature := group.Creatures[msg.creatureIndex]
	p.s.LockTabs()
	p.hitPointForm = newHitPointForm(msg.groupIndex, msg.creatureIndex, creature.GetName(), p.s.GetContentWidth(), p.s.GetContentHeight(), msg.isDamage)
	return p.hitPointForm.Init()
}

func (p *encounterPage) cancelHitPointForm() tea.Cmd {
	p.hitPointForm = nil
	p.s.UnlockTabs()
	return nil
}

func (p *encounterPage) processHitPointAdjustment(msg hitPointFormSubmittedMsg) tea.Cmd {
	p.hitPointForm = nil
	p.s.UnlockTabs()

	if p.encounter == nil || len(p.encounter.InitiativeGroups) <= msg.groupIndex {
		return nil
	}

	group := &p.encounter.InitiativeGroups[msg.groupIndex]
	if len(group.Creatures) <= msg.creatureIndex {
		return nil
	}

	creature := group.Creatures[msg.creatureIndex]

	// Only process hit point adjustments for monsters
	if monster, ok := creature.(dnd.Monster); ok {
		if msg.adjustmentType == HitPointDamage {
			monster.HitPoints -= msg.adjustment
			if monster.HitPoints < 0 {
				monster.HitPoints = 0
			}
		} else {
			monster.HitPoints += msg.adjustment
			if monster.HitPoints > monster.MaximumHitPoints {
				monster.HitPoints = monster.MaximumHitPoints
			}
		}
		group.Creatures[msg.creatureIndex] = monster

		// Create status message with bold components and consistent coloring
		var statusMessage string
		if msg.adjustmentType == HitPointDamage {
			redBold := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
			redNormal := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
			creatureName := redBold.Render(creature.GetName())
			amount := redBold.Render(fmt.Sprintf("%d", msg.adjustment))
			actionType := redNormal.Render("damage")
			took := redNormal.Render(" took ")
			statusMessage = creatureName + took + amount + redNormal.Render(" ") + actionType
		} else {
			greenBold := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
			greenNormal := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
			creatureName := greenBold.Render(creature.GetName())
			amount := greenBold.Render(fmt.Sprintf("%d", msg.adjustment))
			actionType := greenNormal.Render("hit points")
			healed := greenNormal.Render(" healed ")
			statusMessage = creatureName + healed + amount + greenNormal.Render(" ") + actionType
		}

		// Return a status message command for the encounter delegate's list
		return p.encounterDelegate.list.NewStatusMessage(statusMessage)
	}

	return nil
}

package initiative

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

	characters map[string]*dnd.Character
	sources    map[string]*dnd.Source

	encounter *dnd.Encounter

	emptyState        *components.EmptyState
	encounterForm     *encounterForm
	hitPointForm      *hitPointForm
	cancellationForm  *cancellationForm
	encounterDelegate *encounterDelegate
	help              help.Model
}

func newEncounterPage(s *skeleton.Skeleton, characters map[string]*dnd.Character, sources map[string]*dnd.Source) *encounterPage {
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
		if p.cancellationForm != nil {
			p.cancellationForm.SetSize(p.width, p.height)
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
		if p.isAddingNewEncounter() || p.isAdjustingHitPoints() || p.isCancellingEncounter() {
			break
		}

		if key.Matches(msg, p.keys.NewEncounter) && p.isEmptyState() {
			return p, tea.Sequence(append(cmds, p.beginNewEncounterForm())...)
		}
		if key.Matches(msg, p.keys.NextTurn) && p.encounter != nil {
			p.encounter.AdvanceTurn()
			p.updateRoundWidget()
			p.navigateToCurrentTurn()
			return p, nil
		}
		if key.Matches(msg, p.keys.PrevTurn) && p.encounter != nil {
			p.encounter.PreviousTurn()
			p.updateRoundWidget()
			p.navigateToCurrentTurn()
			return p, nil
		}
		if key.Matches(msg, p.keys.EndEncounter) && p.encounter != nil {
			return p, tea.Sequence(append(cmds, p.beginCancellationForm())...)
		}
		if key.Matches(msg, p.keys.ToggleHelp) {
			p.help.ShowAll = !p.help.ShowAll
			return p, nil
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
		return p, tea.Sequence(append(cmds, p.adjustHitPoints(msg))...)
	case cancellationConfirmedMsg:
		return p, tea.Sequence(append(cmds, p.endEncounter())...)
	case cancellationCancelledMsg:
		return p, tea.Sequence(append(cmds, p.cancelCancellationForm())...)
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

	if p.isCancellingEncounter() {
		form, cmd := p.cancellationForm.Update(msg)
		cmds = append(cmds, cmd)
		if form, ok := form.(*cancellationForm); ok {
			p.cancellationForm = form
		}
	}

	if p.encounter != nil && p.encounterDelegate != nil && !p.isAdjustingHitPoints() && !p.isCancellingEncounter() {
		cmds = append(cmds, p.encounterDelegate.Update(msg, p.encounter))
	}

	return p, tea.Batch(cmds...)
}

func (p *encounterPage) View() string {
	// Calculate help view for all cases
	p.help.Width = p.s.GetContentWidth()
	helpStyle := lipgloss.NewStyle().Padding(0, 1)
	helpView := helpStyle.Render(p.help.View(p))

	var content string
	switch true {
	case p.isCancellingEncounter():
		content = p.cancellationForm.View()
	case p.isAdjustingHitPoints():
		content = p.hitPointForm.View()
	case p.encounter != nil:
		// Calculate available height for list (subtract help height)
		helpHeight := lipgloss.Height(helpView)
		listHeight := p.s.GetContentHeight() - helpHeight

		p.encounterDelegate.SetSize(p.s.GetContentWidth(), listHeight)

		var buf strings.Builder
		p.encounterDelegate.Render(&buf, p.encounter)
		content = buf.String()
	case p.isAddingNewEncounter():
		content = p.encounterForm.View()
	case p.isEmptyState():
		// Calculate available height for empty state (subtract help height)
		helpHeight := lipgloss.Height(helpView)
		emptyStateHeight := p.s.GetContentHeight() - helpHeight
		p.emptyState.SetSize(p.s.GetContentWidth(), emptyStateHeight)
		content = p.emptyState.View()
	default:
		// This shouldn't be possible, but we return to satisfy the types.
		content = "No view"
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, helpView)
}

func (p *encounterPage) Key() string {
	return "encounter"
}

func (p *encounterPage) Title() string {
	title := icons.EncounterTab.Join("Encounter", nil)

	if p.encounter != nil {
		summary := p.encounter.Summary()
		if len(summary) > 18 {
			summary = summary[:15] + "..."
		}
		return fmt.Sprintf("%s > %s", title, summary)
	}

	return title
}

func (p *encounterPage) FullHelp() [][]key.Binding {
	// Update key states based on current state
	p.updateKeyStates()

	// Handle form states first
	if p.isAddingNewEncounter() && p.encounterForm != nil {
		formKeys := p.encounterForm.HelpKeys()
		// Group form keys by type if needed, or return as single group
		return [][]key.Binding{formKeys}
	}

	if p.isAdjustingHitPoints() && p.hitPointForm != nil {
		formKeys := p.hitPointForm.HelpKeys()
		return [][]key.Binding{formKeys}
	}

	if p.isCancellingEncounter() && p.cancellationForm != nil {
		formKeys := p.cancellationForm.HelpKeys()
		return [][]key.Binding{formKeys}
	}

	if p.isEmptyState() {
		// Empty state: only show new and help
		return [][]key.Binding{{
			p.keys.NewEncounter,
			p.keys.ToggleHelp,
		}}
	}

	// Organize keys into the requested columns:
	// 1. List controls (up, down, prev page, next page, filter)
	var listControls []key.Binding
	if p.encounterDelegate != nil {
		listControls = []key.Binding{
			p.encounterDelegate.list.KeyMap.CursorUp,
			p.encounterDelegate.list.KeyMap.CursorDown,
			p.encounterDelegate.list.KeyMap.PrevPage,
			p.encounterDelegate.list.KeyMap.NextPage,
			p.encounterDelegate.list.KeyMap.Filter,
		}
	}

	// 2. Turn controls (prev turn, next turn)
	turnControls := []key.Binding{
		p.keys.PrevTurn,
		p.keys.NextTurn,
	}

	// 3. Delegate keys (damage, heal)
	var delegateKeys []key.Binding
	if p.encounterDelegate != nil {
		delegateKeys = []key.Binding{
			p.encounterDelegate.creatureKeys.dealDamage,
			p.encounterDelegate.creatureKeys.heal,
		}
	}

	// 4. End/help
	endHelpKeys := []key.Binding{
		p.keys.EndEncounter,
		p.keys.ToggleHelp,
	}

	return [][]key.Binding{listControls, turnControls, delegateKeys, endHelpKeys}
}

func (p *encounterPage) ShortHelp() []key.Binding {
	// Update key states based on current state
	p.updateKeyStates()

	// Handle form states first - show essential form keys
	if p.isAddingNewEncounter() && p.encounterForm != nil {
		formKeys := p.encounterForm.HelpKeys()
		// Return a subset of essential form keys, or first few
		if len(formKeys) > 4 {
			return formKeys[:4]
		}
		return formKeys
	}

	if p.isAdjustingHitPoints() && p.hitPointForm != nil {
		return p.hitPointForm.HelpKeys()
	}

	if p.isCancellingEncounter() && p.cancellationForm != nil {
		return p.cancellationForm.HelpKeys()
	}

	// Return only essential keys in specified order: up, down, prev turn, next turn, damage, heal, help
	keys := []key.Binding{}

	// Add new encounter key when in empty state
	if p.isEmptyState() {
		keys = append(keys, p.keys.NewEncounter)
	} else {
		// When encounter exists, show essential navigation and action keys
		if p.encounterDelegate != nil {
			keys = append(keys,
				p.encounterDelegate.list.KeyMap.CursorUp,   // up
				p.encounterDelegate.list.KeyMap.CursorDown, // down
			)
		}
		keys = append(keys,
			p.keys.PrevTurn, // prev turn
			p.keys.NextTurn, // next turn
		)
		if p.encounterDelegate != nil {
			keys = append(keys,
				p.encounterDelegate.creatureKeys.dealDamage, // damage
				p.encounterDelegate.creatureKeys.heal,       // heal
			)
		}
	}

	// Always show help
	keys = append(keys, p.keys.ToggleHelp)

	return keys
}

type EncounterPageKeyMap struct {
	NewEncounter key.Binding
	NextTurn     key.Binding
	PrevTurn     key.Binding
	EndEncounter key.Binding
	ToggleHelp   key.Binding
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
			key.WithHelp("←", "prev turn"),
		),
		EndEncounter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "end"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

func (p *encounterPage) isAddingNewEncounter() bool {
	return p.encounterForm != nil
}

func (p *encounterPage) isAdjustingHitPoints() bool {
	return p.hitPointForm != nil
}

func (p *encounterPage) isCancellingEncounter() bool {
	return p.cancellationForm != nil
}

func (p *encounterPage) isEmptyState() bool {
	return !p.isAddingNewEncounter() && !p.isAdjustingHitPoints() && !p.isCancellingEncounter() && p.encounter == nil
}

func (p *encounterPage) updateKeyStates() {
	// Only manage keys that change based on state
	hasEncounter := p.encounter != nil
	isEmpty := p.isEmptyState()

	// Check if we're at the very beginning (round 1, first turn)
	canGoBack := hasEncounter && !(p.encounter.Round() == 1 && p.encounter.TurnIndex() == 0)

	// Check if list is currently filtering
	isFiltering := false
	if p.encounterDelegate != nil {
		isFiltering = p.encounterDelegate.list.FilterState() == list.Filtering
	}

	// These keys are contextually enabled/disabled
	p.keys.NewEncounter.SetEnabled(isEmpty)
	p.keys.PrevTurn.SetEnabled(canGoBack) // Disable when first turn of first round
	p.keys.NextTurn.SetEnabled(hasEncounter)
	p.keys.EndEncounter.SetEnabled(hasEncounter && !isFiltering) // Disable when filtering
	// Help key is always enabled
	p.keys.ToggleHelp.SetEnabled(true)

	// Update delegate key states if available
	if p.encounterDelegate != nil {
		p.encounterDelegate.updateKeyStates()

		// Disable creature actions when no encounter
		if !hasEncounter {
			p.encounterDelegate.creatureKeys.dealDamage.SetEnabled(false)
			p.encounterDelegate.creatureKeys.heal.SetEnabled(false)
		}
	}
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
		p.s.UpdateWidgetValue("round", fmt.Sprintf("Round: %d", p.encounter.Round()))
	}
}

func (p *encounterPage) addNewEncounter(submission encounterFormSubmittedMsg) tea.Cmd {
	p.encounterForm = nil
	p.encounter = submission.encounter

	p.s.UpdatePageTitle(p.Key(), p.Title())

	// Add round tracking widget and unlock tabs
	p.s.AddWidget("round", fmt.Sprintf("Round: %d", p.encounter.Round()))
	p.s.UnlockTabs()

	return nil
}

func (p *encounterPage) endEncounter() tea.Cmd {
	p.encounter = nil

	// Clean up cancellation form
	p.cancellationForm = nil
	p.s.UnlockTabs()
	p.s.UpdatePageTitle(p.Key(), p.Title())

	// Clear round tracking widget
	p.s.DeleteWidget("round")

	return nil
}

func (p *encounterPage) beginHitPointForm(msg adjustHitPointsMsg) tea.Cmd {
	if p.encounter == nil || len(p.encounter.InitiativeGroups()) <= msg.groupIndex {
		return nil
	}

	group := p.encounter.InitiativeGroups()[msg.groupIndex]
	if len(group.Creatures()) <= msg.creatureIndex {
		return nil
	}

	creature := group.Creatures()[msg.creatureIndex]
	p.s.LockTabs()
	p.hitPointForm = newHitPointForm(msg.groupIndex, msg.creatureIndex, (*creature).Name(), p.s.GetContentWidth(), p.s.GetContentHeight(), msg.isDamage)
	return p.hitPointForm.Init()
}

func (p *encounterPage) cancelHitPointForm() tea.Cmd {
	p.hitPointForm = nil
	p.s.UnlockTabs()
	return nil
}

func (p *encounterPage) adjustHitPoints(msg hitPointFormSubmittedMsg) tea.Cmd {
	p.hitPointForm = nil
	p.s.UnlockTabs()

	if p.encounter == nil || len(p.encounter.InitiativeGroups()) <= msg.groupIndex {
		return nil
	}

	groups := p.encounter.InitiativeGroups()
	if len(groups[msg.groupIndex].Creatures()) <= msg.creatureIndex {
		return nil
	}

	// Convert damage to negative, healing to positive
	var adjustment int
	if msg.adjustmentType == HitPointDamage {
		adjustment = -msg.adjustment
	} else {
		adjustment = msg.adjustment
	}

	// Update the creature with new hit points and get actual adjustment
	actualAdjustment := p.encounter.UpdateCreature(msg.groupIndex, msg.creatureIndex, adjustment)

	// Return a status message command for the encounter delegate's list
	return p.encounterDelegate.list.NewStatusMessage(
		getHitPointAdjustmentStatusMessage(*groups[msg.groupIndex].Creatures()[msg.creatureIndex], actualAdjustment),
	)
}

func (p *encounterPage) beginCancellationForm() tea.Cmd {
	p.s.LockTabs()
	p.cancellationForm = newCancellationForm(p.s.GetContentWidth(), p.s.GetContentHeight())
	return p.cancellationForm.Init()
}

func (p *encounterPage) cancelCancellationForm() tea.Cmd {
	p.cancellationForm = nil
	p.s.UnlockTabs()
	return nil
}

func (p *encounterPage) navigateToCurrentTurn() {
	if p.encounter == nil || p.encounterDelegate == nil {
		return
	}

	// Calculate which item in the list corresponds to the current turn
	currentTurnIndex := p.encounter.TurnIndex()
	itemIndex := 0

	// Count items up to the current turn group
	for groupIndex, group := range p.encounter.InitiativeGroups() {
		if groupIndex == currentTurnIndex {
			// We found the current turn group, itemIndex points to the first creature in this group
			break
		}
		// Add all creatures in this group to the count
		itemIndex += len(group.Creatures())
	}

	// Select the first creature in the current turn group
	// This will automatically handle pagination to show the selected item
	p.encounterDelegate.list.Select(itemIndex)
}

// getHitPointAdjustmentStatusMessage returns a styled message for hit point adjustments
func getHitPointAdjustmentStatusMessage(creature dnd.Creature, amount int) string {
	name := creature.Name()
	var messageStyle lipgloss.Style
	var message string

	if amount < 0 {
		messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		message = fmt.Sprintf("%s took %d damage", name, -amount)
	} else {
		messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
		message = fmt.Sprintf("%s healed %d hit points", name, amount)
	}

	return messageStyle.Render(message)
}

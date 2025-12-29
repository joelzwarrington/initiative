package ui

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"initiative/internal/components"
	"initiative/internal/dnd"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/termkit/skeleton"
)

var _ tea.Model = (*encounter)(nil)

type encounterView int

const (
	encounterPlaceholder encounterView = iota
	encounterCreateForm
	encounterDetail
)

type encounter struct {
	dnd.Encounter

	skeleton *skeleton.Skeleton
	party    *map[string]dnd.Character
	sources  map[string]*dnd.Source

	view                encounterView
	encounterCreateForm *encounterCreationForm
	list                list.Model
	help                help.Model
	placeholderKeys     encounterPlaceholderKeyMap
	detailKeys          encounterDetailKeyMap
	currentTurn         int // Index of current turn in initiative order
}

func newEncounter(skeleton *skeleton.Skeleton, party *map[string]dnd.Character, sources map[string]*dnd.Source) *encounter {
	// Create empty list for initiative groups
	delegate := &initiativeGroupItemDelegate{
		health: components.NewHealth(30), // 30 chars width for health bar
	}
	initiativeList := list.New([]list.Item{}, delegate, skeleton.GetContentWidth(), skeleton.GetContentHeight())
	initiativeList.SetStatusBarItemName("group", "groups")
	initiativeList.SetShowTitle(false)
	initiativeList.SetShowStatusBar(false)
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()

	return &encounter{
		skeleton: skeleton,
		party:    party,
		sources:  sources,

		view:            encounterPlaceholder,
		list:            initiativeList,
		help:            help.New(),
		placeholderKeys: newEncounterPlaceholderKeyMap(),
		detailKeys:      newEncounterDetailKeyMap(),
	}
}

func (e encounter) Init() tea.Cmd {
	return nil
}

func (e encounter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch e.view {
		case encounterPlaceholder:
			if key.Matches(msg, e.placeholderKeys.startEncounter) {
				return e, tea.Cmd(func() tea.Msg {
					return startEncounterCreateMsg{}
				})
			}
		case encounterDetail:
			if key.Matches(msg, e.detailKeys.back) {
				e.view = encounterPlaceholder
				e.InitiativeGroups = []dnd.InitiativeGroup{}
				e.Summary = ""
				e.StartedAt = time.Time{}
				e.EndedAt = time.Time{}
				e.Round = 0
				e.currentTurn = 0
				e.encounterCreateForm = nil

				// Remove widgets when stopping encounter
				e.skeleton.DeleteAllWidgets()

				return e, nil
			}
			if key.Matches(msg, e.detailKeys.previousTurn) {
				e.previousTurn()
				return e, nil
			}
			if key.Matches(msg, e.detailKeys.nextTurn) {
				e.nextTurn()
				return e, nil
			}
		}
	case startEncounterCreateMsg:
		e.encounterCreateForm = newEncounterCreateForm(e.skeleton, e.party, e.sources)
		e.view = encounterCreateForm
		return e, e.encounterCreateForm.Init()
	case createEncounterMsg:
		e.Summary = msg.summary
		e.StartedAt = time.Now()
		e.Round = 1 // Start at round 1
		e.InitiativeGroups = msg.initiativeGroups
		e.encounterCreateForm = nil

		// Sort initiative groups by initiative value (highest to lowest)
		sort.Slice(e.InitiativeGroups, func(i, j int) bool {
			return e.InitiativeGroups[i].Initiative > e.InitiativeGroups[j].Initiative
		})

		// Update the list with sorted initiative groups
		e.updateInitiativeList()
		e.list.SetItems(e.createInitiativeListItems())

		// Initialize current turn and add round tracking widget
		e.currentTurn = 0 // Start with first creature in initiative order
		e.skeleton.AddWidget("round", fmt.Sprintf("Round: %d", e.Round))

		e.skeleton.UnlockTabs()
		e.view = encounterDetail
		return e, nil
	case cancelEncounterCreationMsg:
		e.encounterCreateForm = nil
		e.skeleton.UnlockTabs()
		e.view = encounterPlaceholder
		return e, nil
	}

	switch e.view {
	case encounterCreateForm:
		{
			e.skeleton.LockTabs()

			if e.encounterCreateForm != nil {
				var cmd tea.Cmd
				e.encounterCreateForm, cmd = e.encounterCreateForm.Update(msg)
				return e, cmd
			}
		}
	case encounterDetail:
		{
			var cmd tea.Cmd
			e.list, cmd = e.list.Update(msg)
			return e, cmd
		}
	}
	return e, nil
}

func (e encounter) View() string {
	switch e.view {
	case encounterPlaceholder:
		{
			availHeight := e.skeleton.GetContentHeight()
			e.help.Width = e.skeleton.GetContentWidth()
			helpStyle := lipgloss.NewStyle().Padding(0, 2)
			helpView := helpStyle.Render(e.help.View(e.placeholderKeys))
			availHeight = availHeight - lipgloss.Height(helpView)

			// Create main content area
			placeholderStyle := lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("240")).
				Align(lipgloss.Center)
			content := placeholderStyle.Render("No encounter started...")
			contentArea := lipgloss.NewStyle().
				Height(availHeight).
				Width(e.skeleton.GetContentWidth()).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render(content)

			return lipgloss.JoinVertical(lipgloss.Left, contentArea, helpView)
		}
	case encounterCreateForm:
		{
			if e.encounterCreateForm != nil {
				return e.encounterCreateForm.View()
			}
			return ""
		}
	case encounterDetail:
		{
			helpStyle := lipgloss.NewStyle().Padding(0, 1)
			availHeight := e.skeleton.GetContentHeight()
			e.help.Width = e.skeleton.GetContentWidth()

			// Create header with encounter summary
			headerStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)
			header := headerStyle.Render(fmt.Sprintf("Encounter: %s", e.Summary))
			help := helpStyle.Render(e.help.View(e.detailKeys))

			listHeight := availHeight - lipgloss.Height(header) - lipgloss.Height(help)

			e.list.SetHeight(listHeight)
			e.list.SetWidth(e.skeleton.GetContentWidth())

			return lipgloss.JoinVertical(lipgloss.Left, header, e.list.View(), help)
		}
	}

	return ""
}

// Widget management
func (e *encounter) updateRoundWidget() {
	e.skeleton.UpdateWidgetValue("round", fmt.Sprintf("Round: %d", e.Round))
}

// Initiative list management
func (e *encounter) createInitiativeListItems() []list.Item {
	items := []list.Item{}
	for i, group := range e.InitiativeGroups {
		items = append(items, initiativeGroupItem{
			group:         group,
			isCurrentTurn: i == e.currentTurn,
		})
	}
	return items
}

func (e *encounter) updateInitiativeList() {
	e.list.SetItems(e.createInitiativeListItems())
}

// Turn navigation
func (e *encounter) nextTurn() {
	if len(e.InitiativeGroups) == 0 {
		return
	}

	e.currentTurn++
	if e.currentTurn >= len(e.InitiativeGroups) {
		// End of round - go to beginning and increment round
		e.currentTurn = 0
		e.Round++
		e.updateRoundWidget()
	}
	e.updateInitiativeList()
}

func (e *encounter) previousTurn() {
	if len(e.InitiativeGroups) == 0 {
		return
	}

	// Don't allow going back if we're at round 1, turn 0 (start of encounter)
	if e.Round == 1 && e.currentTurn == 0 {
		return
	}

	e.currentTurn--
	if e.currentTurn < 0 {
		// Beginning of round - go to end and decrement round
		e.currentTurn = len(e.InitiativeGroups) - 1
		e.Round--
		e.updateRoundWidget()
	}
	e.updateInitiativeList()
}

// Messages
type startEncounterCreateMsg struct{}
type cancelEncounterCreationMsg struct{}

// Key mappings
type encounterPlaceholderKeyMap struct {
	startEncounter key.Binding
}

func newEncounterPlaceholderKeyMap() encounterPlaceholderKeyMap {
	return encounterPlaceholderKeyMap{
		startEncounter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new encounter"),
		),
	}
}

func (k encounterPlaceholderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.startEncounter}
}

func (k encounterPlaceholderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.startEncounter},
	}
}

type encounterDetailKeyMap struct {
	back         key.Binding
	previousTurn key.Binding
	nextTurn     key.Binding
}

func newEncounterDetailKeyMap() encounterDetailKeyMap {
	return encounterDetailKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "stop encounter"),
		),
		previousTurn: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "previous turn"),
		),
		nextTurn: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next turn"),
		),
	}
}

func (k encounterDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.previousTurn, k.nextTurn, k.back}
}

func (k encounterDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.previousTurn, k.nextTurn},
		{k.back},
	}
}

// List item for initiative groups
var _ list.Item = (*initiativeGroupItem)(nil)

type initiativeGroupItem struct {
	group         dnd.InitiativeGroup
	isCurrentTurn bool
}

func (i initiativeGroupItem) FilterValue() string {
	if len(i.group.Creatures) > 0 {
		return i.group.Creatures[0].Name()
	}
	return fmt.Sprintf("Initiative: %d", i.group.Initiative)
}

// List delegate for initiative groups
type initiativeGroupItemDelegate struct {
	health components.Health
}

func (d initiativeGroupItemDelegate) Height() int {
	return 1 // Single line per item - trees will expand as needed
}
func (d initiativeGroupItemDelegate) Spacing() int { return 1 }
func (d *initiativeGroupItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d initiativeGroupItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(initiativeGroupItem)
	if !ok {
		return
	}

	// Styling
	initiativeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))
	creatureStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
	currentTurnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")) // Green for current turn
	selectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")) // Purple for selection
	separatorStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#DDDADA",
		Dark:  "#3C3C3C",
	})

	var content string

	// Check if this is the current turn by accessing encounter's currentTurn
	isCurrentTurn := i.isCurrentTurn

	// Single creature on same line
	if len(i.group.Creatures) == 1 {
		creature := i.group.Creatures[0]

		// Format: "14 • Wizard" with right-aligned initiative (2 chars)
		initiativeNum := initiativeStyle.Render(fmt.Sprintf("%2d", i.group.Initiative))
		separator := separatorStyle.Render(" • ")
		creatureName := creatureStyle.Render(creature.Name())

		// Check if creature is a monster and add health bar
		var healthBar string
		if monster, ok := creature.(dnd.Monster); ok {
			healthBar = " " + d.health.View(monster.HitPoints, monster.MaximumHitPoints)
		}

		if isCurrentTurn {
			// Show green arrow for current turn: "→ 14 • Wizard [health bar]"
			arrow := currentTurnStyle.Render("→")
			content = fmt.Sprintf("%s %s%s%s%s", arrow, initiativeNum, separator, creatureName, healthBar)
		} else {
			// Regular format: "  14 • Wizard [health bar]"
			content = fmt.Sprintf("  %s%s%s%s", initiativeNum, separator, creatureName, healthBar)
		}
	} else {
		// Multiple creatures using tree
		initiativeNum := fmt.Sprintf("%2d", i.group.Initiative) // Right-aligned 2 chars

		// Create tree with initiative as root
		initiativeTree := tree.New().
			Root(initiativeStyle.Render(initiativeNum)).
			EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(5)).
			ItemStyle(creatureStyle)

		// Add each creature as child
		for _, creature := range i.group.Creatures {
			// Check if creature is a monster and add health bar to name
			var creatureDisplay string
			if monster, ok := creature.(dnd.Monster); ok {
				healthBar := d.health.View(monster.HitPoints, monster.MaximumHitPoints)
				creatureDisplay = creature.Name() + " " + healthBar
			} else {
				creatureDisplay = creature.Name()
			}
			initiativeTree.Child(creatureDisplay)
		}

		treeContent := initiativeTree.String()

		if isCurrentTurn {
			// Add green arrow before tree
			arrow := currentTurnStyle.Render("→")
			lines := strings.Split(treeContent, "\n")
			if len(lines) > 0 {
				lines[0] = arrow + " " + lines[0]
				content = strings.Join(lines, "\n")
			} else {
				content = arrow + " " + treeContent
			}
		} else {
			// Add spacing before tree
			lines := strings.Split(treeContent, "\n")
			for i, line := range lines {
				lines[i] = "  " + line
			}
			content = strings.Join(lines, "\n")
		}
	}

	// Apply selection highlighting
	if index == m.Index() {
		content = selectionStyle.Render("> ") + content
	} else {
		content = "  " + content
	}

	fmt.Fprint(w, content)
}

// Encounter creation form model
type encounterCreationStep int

const (
	stepSummaryAndCharacters encounterCreationStep = iota
	stepMonsterSelection
	stepMonsterQuantities
	stepGatheringInitiative
	stepComplete
)

type encounterCreationForm struct {
	step     encounterCreationStep
	form     *huh.Form
	skeleton *skeleton.Skeleton
	party    *map[string]dnd.Character
	sources  map[string]*dnd.Source

	// Form data
	summary                string
	selectedCharacterUUIDs []string
	selectedMonsters       []string
	monsterQuantities      map[string]int
	currentInitiativeIndex int
	initiativeGroups       []dnd.InitiativeGroup
}

func newEncounterCreateForm(skeleton *skeleton.Skeleton, party *map[string]dnd.Character, sources map[string]*dnd.Source) *encounterCreationForm {
	return &encounterCreationForm{
		step:              stepSummaryAndCharacters,
		skeleton:          skeleton,
		party:             party,
		sources:           sources,
		monsterQuantities: make(map[string]int),
		initiativeGroups:  []dnd.InitiativeGroup{},
	}
}

func customFormTheme() *huh.Theme {
	theme := huh.ThemeCharm()

	// Modify the focused error message to remove leading space and set red foreground
	theme.Focused.ErrorMessage = lipgloss.NewStyle().SetString("*").Foreground(lipgloss.AdaptiveColor{Dark: "#ff5555"})

	return theme
}

func customFormKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()

	// Add ESC key to quit the form
	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit"),
	)

	// Ensure help is enabled for the quit binding
	keyMap.Quit.SetEnabled(true)

	return keyMap
}

func (f *encounterCreationForm) Init() tea.Cmd {
	f.createSummaryForm()
	return f.form.Init()
}

func (f *encounterCreationForm) createSummaryForm() {
	var characterOptions []huh.Option[string]

	if f.party != nil {
		for uuid, character := range *f.party {
			characterOptions = append(characterOptions,
				huh.NewOption(character.Name(), uuid).Selected(true),
			)
		}
	}

	f.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("summary").
				Title("Summary").
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Summary is required")
					}
					return nil
				}).Inline(true),
			huh.NewMultiSelect[string]().
				Key("characters").
				Title("Characters").
				Options(characterOptions...),
		),
	)
}

func (f *encounterCreationForm) createMonsterSelectionForm() {
	var monsterOptions []huh.Option[string]

	// Collect monsters from all available sources
	for _, source := range f.sources {
		if source != nil {
			for _, monster := range source.Monsters {
				monsterOptions = append(monsterOptions,
					huh.NewOption(monster.Name(), monster.Name()),
				)
			}
		}
	}

	// If no monsters available, skip to initiative
	if len(monsterOptions) == 0 {
		f.step = stepGatheringInitiative
		f.createInitiativeForm()
		return
	}

	f.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("monsters").
				Title("Select Monsters").
				Options(monsterOptions...),
		),
	)
}

func (f *encounterCreationForm) createMonsterQuantitiesForm() {
	if len(f.selectedMonsters) == 0 {
		f.step = stepGatheringInitiative
		f.createInitiativeForm()
		return
	}

	fields := []huh.Field{
		huh.NewNote().Title("Monster Quantities"),
	}

	for _, monsterName := range f.selectedMonsters {
		fields = append(fields,
			huh.NewInput().
				Key(fmt.Sprintf("quantity_%s", monsterName)).
				Title(monsterName).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Quantity is required")
					}
					value, err := strconv.Atoi(strings.TrimSpace(str))
					if err != nil || value <= 0 {
						return fmt.Errorf("Quantity must be a positive number")
					}
					return nil
				}),
		)
	}

	f.form = huh.NewForm(
		huh.NewGroup(fields...),
	)
}

func (f *encounterCreationForm) createInitiativeForm() {
	if len(f.selectedCharacterUUIDs) == 0 && len(f.selectedMonsters) == 0 {
		f.step = stepComplete
		return
	}

	// Create all initiative inputs
	fields := []huh.Field{
		huh.NewNote().Title("Initiative"),
	}

	// Add character initiative fields
	for _, uuid := range f.selectedCharacterUUIDs {
		var characterName string
		if f.party != nil {
			if character, exists := (*f.party)[uuid]; exists {
				characterName = character.Name()
			}
		}

		fields = append(fields,
			huh.NewInput().
				Key(fmt.Sprintf("initiative_%s", uuid)).
				Title(fmt.Sprintf("%s", characterName)).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Initiative is required")
					}
					value, err := strconv.Atoi(strings.TrimSpace(str))
					if err != nil || value <= 0 {
						return fmt.Errorf("Initiative must be a positive number")
					}
					return nil
				}),
		)
	}

	// Add monster initiative fields (one per monster type)
	for _, monsterName := range f.selectedMonsters {
		fields = append(fields,
			huh.NewInput().
				Key(fmt.Sprintf("monster_initiative_%s", monsterName)).
				Title(fmt.Sprintf("%s", monsterName)).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Initiative is required")
					}
					value, err := strconv.Atoi(strings.TrimSpace(str))
					if err != nil || value <= 0 {
						return fmt.Errorf("Initiative must be a positive number")
					}
					return nil
				}),
		)
	}

	// Create form with group containing all fields
	f.form = huh.NewForm(
		huh.NewGroup(fields...),
	)
}

func (f *encounterCreationForm) Update(msg tea.Msg) (*encounterCreationForm, tea.Cmd) {
	if f.form == nil {
		return f, nil
	}

	form, cmd := f.form.Update(msg)
	if newForm, ok := form.(*huh.Form); ok {
		f.form = newForm
	}

	// Handle form quit (ESC pressed)
	if f.form.State == huh.StateAborted {
		return f, tea.Cmd(func() tea.Msg {
			return cancelEncounterCreationMsg{}
		})
	}

	if f.form.State == huh.StateCompleted {
		switch f.step {
		case stepSummaryAndCharacters:
			f.summary = f.form.GetString("summary")
			f.selectedCharacterUUIDs = f.form.Get("characters").([]string)
			f.step = stepMonsterSelection
			f.createMonsterSelectionForm()
			return f, f.form.Init()

		case stepMonsterSelection:
			f.selectedMonsters = f.form.Get("monsters").([]string)
			f.step = stepMonsterQuantities
			f.createMonsterQuantitiesForm()
			if f.step == stepGatheringInitiative {
				// Skip quantities if no monsters selected
				f.createInitiativeForm()
				if f.step == stepComplete {
					return f, tea.Cmd(func() tea.Msg {
						return createEncounterMsg{
							summary:          f.summary,
							initiativeGroups: f.initiativeGroups,
						}
					})
				}
			}
			return f, f.form.Init()

		case stepMonsterQuantities:
			// Parse monster quantities
			for _, monsterName := range f.selectedMonsters {
				quantityKey := fmt.Sprintf("quantity_%s", monsterName)
				quantityStr := f.form.GetString(quantityKey)
				quantity, err := strconv.Atoi(strings.TrimSpace(quantityStr))
				if err != nil {
					quantity = 1 // Default to 1 if parsing fails
				}
				f.monsterQuantities[monsterName] = quantity
			}
			f.step = stepGatheringInitiative
			f.createInitiativeForm()
			if f.step == stepComplete {
				return f, tea.Cmd(func() tea.Msg {
					return createEncounterMsg{
						summary:          f.summary,
						initiativeGroups: f.initiativeGroups,
					}
				})
			}
			return f, f.form.Init()

		case stepGatheringInitiative:
			// Parse all initiative values
			for _, uuid := range f.selectedCharacterUUIDs {
				initiativeKey := fmt.Sprintf("initiative_%s", uuid)
				initiativeStr := f.form.GetString(initiativeKey)

				// Parse initiative value (validation already ensures it's a positive integer)
				initiativeValue, err := strconv.Atoi(strings.TrimSpace(initiativeStr))
				if err != nil {
					// This shouldn't happen due to validation, but default to 1
					initiativeValue = 1
				}

				// Get character and create initiative group
				if f.party != nil {
					if character, exists := (*f.party)[uuid]; exists {
						group := dnd.InitiativeGroup{
							Initiative: initiativeValue,
							Creatures:  []dnd.Creature{character},
						}
						f.initiativeGroups = append(f.initiativeGroups, group)
					}
				}
			}

			// Process monster initiatives
			// Create a map for quick monster lookup from all sources
			monsterMap := make(map[string]dnd.Monster)
			for _, source := range f.sources {
				if source != nil {
					for _, monster := range source.Monsters {
						monsterMap[monster.Name()] = monster
					}
				}
			}

			for _, monsterName := range f.selectedMonsters {
				initiativeKey := fmt.Sprintf("monster_initiative_%s", monsterName)
				initiativeStr := f.form.GetString(initiativeKey)

				// Parse initiative value
				initiativeValue, err := strconv.Atoi(strings.TrimSpace(initiativeStr))
				if err != nil {
					initiativeValue = 1 // Default to 1 if parsing fails
				}

				// Get quantity for this monster type
				quantity := f.monsterQuantities[monsterName]
				if quantity <= 0 {
					quantity = 1 // Default to 1 if not set
				}

				// Create monsters and add to initiative group
				if monster, exists := monsterMap[monsterName]; exists {
					var creatures []dnd.Creature
					for i := 0; i < quantity; i++ {
						// Each monster in the group is a copy of the template with initialized health
						monsterCopy := dnd.NewMonster(monster.MonsterName, monster.StatBlock)
						creatures = append(creatures, monsterCopy)
					}

					// Create initiative group for all monsters of this type
					group := dnd.InitiativeGroup{
						Initiative: initiativeValue,
						Creatures:  creatures,
					}
					f.initiativeGroups = append(f.initiativeGroups, group)
				}
			}

			// All initiatives processed, complete the form
			f.step = stepComplete
			return f, tea.Cmd(func() tea.Msg {
				return createEncounterMsg{
					summary:          f.summary,
					initiativeGroups: f.initiativeGroups,
				}
			})
		}
	}

	return f, cmd
}

func (f *encounterCreationForm) View() string {
	if f.form != nil {
		horizontalPadding := 4 // 2 padding on left and right = 4 total
		verticalPadding := 2   // 1 padding on top and bottom = 2 total
		formHeight := f.skeleton.GetContentHeight() - verticalPadding
		formWidth := f.skeleton.GetContentWidth() - horizontalPadding

		f.form.WithHeight(formHeight).
			WithWidth(formWidth).
			WithShowErrors(true).
			WithShowHelp(true).
			WithAccessible(true)

		// Apply padding: 1 vertical, 2 horizontal
		paddingStyle := lipgloss.NewStyle().Padding(1, 2)
		return paddingStyle.Render(f.form.View())
	}
	return ""
}

type createEncounterMsg struct {
	summary          string
	initiativeGroups []dnd.InitiativeGroup
}

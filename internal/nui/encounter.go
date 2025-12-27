package nui

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
	Encounter

	skeleton *skeleton.Skeleton
	party    *map[string]Character

	view                encounterView
	encounterCreateForm *encounterCreationForm
	list                list.Model
	help                help.Model
	placeholderKeys     encounterPlaceholderKeyMap
	detailKeys          encounterDetailKeyMap
}

func newEncounter(skeleton *skeleton.Skeleton, party *map[string]Character) *encounter {
	// Create empty list for initiative groups
	initiativeList := list.New([]list.Item{}, &initiativeGroupItemDelegate{}, skeleton.GetContentWidth(), skeleton.GetContentHeight())
	initiativeList.SetStatusBarItemName("group", "groups")
	initiativeList.SetShowTitle(false)
	initiativeList.SetShowStatusBar(false)
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()

	return &encounter{
		skeleton: skeleton,
		party:    party,

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
				e.IniativeGroups = []IniativeGroup{}
				e.Summary = ""
				e.StartedAt = time.Time{}
				e.EndedAt = time.Time{}
				e.encounterCreateForm = nil
				return e, nil
			}
		}
	case startEncounterCreateMsg:
		e.encounterCreateForm = newEncounterCreateForm(e.skeleton, e.party)
		e.view = encounterCreateForm
		return e, e.encounterCreateForm.Init()
	case createEncounterMsg:
		e.Summary = msg.summary
		e.StartedAt = time.Now()
		e.IniativeGroups = msg.initiativeGroups
		e.encounterCreateForm = nil

		// Update the list with initiative groups
		items := []list.Item{}
		for _, group := range e.IniativeGroups {
			items = append(items, initiativeGroupItem{group: group})
		}
		e.list.SetItems(items)
		e.view = encounterDetail
		return e, nil
	case cancelEncounterCreationMsg:
		e.encounterCreateForm = nil
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
			// Calculate available height for content
			helpStyle := lipgloss.NewStyle().PaddingBottom(1)
			helpView := helpStyle.Render(e.help.View(e.placeholderKeys))
			availHeight := e.skeleton.GetContentHeight() - lipgloss.Height(helpView)

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
			// Calculate available height for content
			helpStyle := lipgloss.NewStyle().PaddingBottom(1)
			helpView := helpStyle.Render(e.help.View(e.detailKeys))
			availHeight := e.skeleton.GetContentHeight() - lipgloss.Height(helpView)

			// Create header with encounter summary
			headerStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)
			header := headerStyle.Render(fmt.Sprintf("Encounter: %s", e.Summary))

			// Set list dimensions accounting for header and help
			headerHeight := lipgloss.Height(header) + 1 // +1 for margin
			listHeight := availHeight - headerHeight

			e.list.SetHeight(listHeight)
			e.list.SetWidth(e.skeleton.GetContentWidth())

			// Combine header and list
			content := lipgloss.JoinVertical(lipgloss.Left, header, e.list.View())

			return lipgloss.JoinVertical(lipgloss.Left, content, helpView)
		}
	}

	return ""
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
	back key.Binding
}

func newEncounterDetailKeyMap() encounterDetailKeyMap {
	return encounterDetailKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k encounterDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.back}
}

func (k encounterDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.back},
	}
}

// List item for initiative groups
var _ list.Item = (*initiativeGroupItem)(nil)

type initiativeGroupItem struct {
	group IniativeGroup
}

func (i initiativeGroupItem) FilterValue() string {
	if len(i.group.Creatures) > 0 {
		return i.group.Creatures[0].Name()
	}
	return fmt.Sprintf("Initiative: %d", i.group.Iniative)
}

// List delegate for initiative groups
type initiativeGroupItemDelegate struct{}

func (d initiativeGroupItemDelegate) Height() int  { return 2 }
func (d initiativeGroupItemDelegate) Spacing() int { return 1 }
func (d *initiativeGroupItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d initiativeGroupItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(initiativeGroupItem)
	if !ok {
		return
	}

	// Initiative value styling
	initiativeText := "Initiative: TBD"
	if i.group.Iniative > 0 {
		initiativeText = fmt.Sprintf("Initiative: %d", i.group.Iniative)
	}

	initiativeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))

	// Creatures list
	creatureNames := []string{}
	for _, creature := range i.group.Creatures {
		creatureNames = append(creatureNames, creature.Name())
	}
	creaturesText := strings.Join(creatureNames, ", ")

	creatureStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// Combine text
	content := initiativeStyle.Render(initiativeText) + "\n" +
		creatureStyle.Render("  "+creaturesText)

	// Apply selection styling
	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170")).
				Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(content))
}

// Encounter creation form model
type encounterCreationStep int

const (
	stepSummaryAndCharacters encounterCreationStep = iota
	stepGatheringInitiative
	stepComplete
)

type encounterCreationForm struct {
	step     encounterCreationStep
	form     *huh.Form
	skeleton *skeleton.Skeleton
	party    *map[string]Character

	// Form data
	summary                string
	selectedCharacterUUIDs []string
	currentInitiativeIndex int
	initiativeGroups       []IniativeGroup
}

func newEncounterCreateForm(skeleton *skeleton.Skeleton, party *map[string]Character) *encounterCreationForm {
	return &encounterCreationForm{
		step:             stepSummaryAndCharacters,
		skeleton:         skeleton,
		party:            party,
		initiativeGroups: []IniativeGroup{},
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

func (f *encounterCreationForm) createInitiativeForm() {
	if len(f.selectedCharacterUUIDs) == 0 {
		f.step = stepComplete
		return
	}

	// Create all initiative inputs
	fields := []huh.Field{
		huh.NewNote().Title("Initiative"),
	}

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
						group := IniativeGroup{
							Iniative:  initiativeValue,
							Creatures: []Creature{character},
						}
						f.initiativeGroups = append(f.initiativeGroups, group)
					}
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
		paddingSize := 2
		formHeight := f.skeleton.GetContentHeight() - paddingSize
		formWidth := f.skeleton.GetContentWidth() - paddingSize

		f.form.WithHeight(formHeight).
			WithWidth(formWidth).
			WithShowErrors(true).
			WithShowHelp(true).
			WithAccessible(true)

		// Apply padding around the entire form
		paddingStyle := lipgloss.NewStyle().Padding(1)
		return paddingStyle.Render(f.form.View())
	}
	return ""
}

type createEncounterMsg struct {
	summary          string
	initiativeGroups []IniativeGroup
}

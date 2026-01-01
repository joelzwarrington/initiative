package initiative

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
	"github.com/joelzwarrington/initiative/internal/form"
)

type creatureChoice struct {
	id       string
	creature dnd.Creature
}

type encounterFormStyles struct {
	container lipgloss.Style
	help      lipgloss.Style
}

func (f *encounterForm) isFilteringActive(form *huh.Form) bool {
	// Get the currently focused field
	focusedField := form.GetFocusedField()
	if focusedField == nil {
		return false
	}

	// Check if the focused field is a MultiSelect and is filtering
	if multiSelect, ok := focusedField.(*huh.MultiSelect[string]); ok {
		return multiSelect.GetFiltering()
	}

	// Check if the focused field is a Select and is filtering
	if selectField, ok := focusedField.(*huh.Select[string]); ok {
		return selectField.GetFiltering()
	}

	return false
}

func (f *encounterForm) getHelpKeys() []key.Binding {
	form := f.getCurrentForm()
	if form == nil {
		return []key.Binding{}
	}

	// Get form's key bindings
	formKeys := form.KeyBinds()

	// Add our custom ESC quit binding if it's enabled
	if f.keys.Quit.Enabled() {
		formKeys = append(formKeys, f.keys.Quit)
	}

	return formKeys
}

type encounterFormCancelledMsg struct{}
type encounterFormSubmittedMsg struct {
	encounter dnd.Encounter
}
type encounterFormPreviousStepMsg struct{}
type encounterFormNextStepMsg struct{}

type encounterForm struct {
	forms  *[]*huh.Form
	styles encounterFormStyles
	keys   *huh.KeyMap

	characters *map[string]dnd.Character
	sources    map[string]*dnd.Source

	width  int
	height int
	help   help.Model
}

func newEncounterForm(c *map[string]dnd.Character, s map[string]*dnd.Source, width int, height int) *encounterForm {
	// Create shared keymap
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Key("summary").Title("Summary").
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Summary is required")
					}
					return nil
				}),
		).Title("Add new encounter\n"),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("characters").
				Title("Characters").
				Options(getCharacterOptions(c)...).
				Validate(func(values []string) error {
					return nil // Allow empty character selection
				}),
		).Title("Add characters to encounter"),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("monsters").
				Title("Monsters").
				Options(getMonsterOptions(s)...).
				Validate(func(values []string) error {
					return nil // Allow empty monster selection
				}),
		).Title("Add monsters to encounter"),
	).WithShowErrors(true).WithShowHelp(false).WithAccessible(true).WithKeyMap(keyMap)

	form.CancelCmd = func() tea.Msg {
		return encounterFormPreviousStepMsg{}
	}
	form.SubmitCmd = func() tea.Msg {
		return encounterFormNextStepMsg{}
	}

	f := &encounterForm{
		width: width, height: height,
		forms: &[]*huh.Form{form},
		keys:  keyMap,
		styles: encounterFormStyles{
			container: lipgloss.NewStyle().Padding(1, 2, 0, 2), // 1 top, 2 horizontal, 0 bottom
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},

		characters: c,
		sources:    s,
		help:       help.New(),
	}

	// Apply sizing using the standard SetSize method
	f.SetSize(width, height)

	return f
}

func (f *encounterForm) Init() tea.Cmd {
	return f.getCurrentForm().Init()
}

func (f *encounterForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	currentForm := f.getCurrentForm()

	if currentForm == nil {
		return f, nil
	}

	switch msg.(type) {
	case encounterFormPreviousStepMsg:
		return f, tea.Sequence(append(cmds, f.prevStep())...)
	case encounterFormNextStepMsg:
		return f, tea.Sequence(append(cmds, f.nextStep())...)
	}

	// Check filtering state before form update
	wasFiltering := f.isFilteringActive(currentForm)

	form, cmd := currentForm.Update(msg)
	cmds = append(cmds, cmd)
	if form, ok := form.(*huh.Form); ok {
		(*f.forms)[len(*f.forms)-1] = form

		// Check filtering state after form update and update keymap if changed
		isFiltering := f.isFilteringActive(form)
		if wasFiltering != isFiltering {
			f.keys.Quit.SetEnabled(!isFiltering)
		}

		// Handle form quit (ESC pressed)
		if form.State == huh.StateAborted {
			return f, tea.Cmd(func() tea.Msg {
				return encounterFormCancelledMsg{}
			})
		}
	}

	return f, tea.Sequence(cmds...)
}

func (f *encounterForm) View() string {
	form := f.getCurrentForm()
	if form == nil {
		return ""
	}

	// Render form content
	formView := f.styles.container.Render(form.View())

	// Render custom help
	f.help.Width = f.width - f.styles.help.GetHorizontalPadding()
	helpKeys := f.getHelpKeys()
	helpView := f.styles.help.Render(f.help.ShortHelpView(helpKeys))

	// Join form and help vertically
	return lipgloss.JoinVertical(lipgloss.Left, formView, helpView)
}

func (f *encounterForm) getCurrentForm() *huh.Form {
	if len(*f.forms) < 1 {
		return nil
	}

	return (*f.forms)[len(*f.forms)-1]
}

func (f *encounterForm) getSummary() string {
	if len(*f.forms) == 0 {
		return ""
	}

	form := (*f.forms)[0]
	return form.GetString("summary")
}

func (f *encounterForm) getCharacters() []creatureChoice {
	if len(*f.forms) == 0 || f.characters == nil {
		return []creatureChoice{}
	}

	form := (*f.forms)[0]
	uuids, ok := form.Get("characters").([]string)
	if !ok {
		return []creatureChoice{}
	}

	var options []creatureChoice
	for _, uuid := range uuids {
		if character, exists := (*f.characters)[uuid]; exists {
			options = append(options, creatureChoice{
				id:       uuid,
				creature: character,
			})
		}
	}

	// Sort by name (case-insensitive)
	sort.Slice(options, func(i, j int) bool {
		return strings.ToLower(options[i].creature.GetName()) < strings.ToLower(options[j].creature.GetName())
	})

	return options
}

func (f *encounterForm) getMonsters() []creatureChoice {
	if len(*f.forms) == 0 {
		return []creatureChoice{}
	}

	form := (*f.forms)[0]
	values, ok := form.Get("monsters").([]string)
	if !ok {
		return []creatureChoice{}
	}

	var options []creatureChoice
	for _, value := range values {
		// value format is "sourceKey:monsterName"
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			continue
		}

		sourceKey := parts[0]
		monsterName := parts[1]

		if source, exists := f.sources[sourceKey]; exists {
			for _, monster := range source.Monsters {
				if monster.GetName() == monsterName {
					options = append(options, creatureChoice{
						id:       value,
						creature: monster,
					})
					break
				}
			}
		}
	}

	// Sort by name (case-insensitive)
	sort.Slice(options, func(i, j int) bool {
		return strings.ToLower(options[i].creature.GetName()) < strings.ToLower(options[j].creature.GetName())
	})

	return options
}

// getCharacterOptions outputs a list of sorted form options from a map of characters
func getCharacterOptions(characters *map[string]dnd.Character) []huh.Option[string] {
	var options []huh.Option[string]
	if characters != nil {
		type entry struct {
			uuid      string
			character dnd.Character
		}

		var entries []entry
		for uuid, character := range *characters {
			entries = append(entries, entry{uuid: uuid, character: character})
		}

		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].character.GetName()) < strings.ToLower(entries[j].character.GetName())
		})

		for _, entry := range entries {
			options = append(options, huh.NewOption(entry.character.GetName(), entry.uuid).Selected(true))
		}
	}
	return options
}

func getMonsterOptions(sources map[string]*dnd.Source) []huh.Option[string] {
	var options []huh.Option[string]
	for sourceKey, source := range sources {
		for _, monster := range source.Monsters {
			value := sourceKey + ":" + monster.GetName()
			options = append(options, huh.NewOption(monster.GetName(), value))
		}
	}
	return options
}

func (f *encounterForm) SetSize(width int, height int) {
	f.width = width
	f.height = height

	// Extract padding from style configuration
	container := f.styles.container
	horizontalPadding := container.GetHorizontalPadding()
	verticalPadding := container.GetVerticalPadding()

	// Calculate help height by rendering a sample help view
	f.help.Width = width - f.styles.help.GetHorizontalPadding()
	helpKeys := f.getHelpKeys()
	sampleHelp := f.styles.help.Render(f.help.ShortHelpView(helpKeys))
	helpHeight := lipgloss.Height(sampleHelp)

	// Account for padding and actual help text height
	formWidth := width - horizontalPadding
	formHeight := height - verticalPadding - helpHeight

	for _, form := range *f.forms {
		if form != nil {
			form.WithHeight(formHeight).WithWidth(formWidth)
		}
	}
}

func (f *encounterForm) prevStep() tea.Cmd {
	// when more than 1 form, remove the last form
	if len(*f.forms) > 1 {
		*f.forms = (*f.forms)[:len(*f.forms)-1]
		return nil
	}

	// when there's just 1 form, send a cancelled message
	return func() tea.Msg {
		return encounterFormCancelledMsg{}
	}
}

func (f *encounterForm) nextStep() tea.Cmd {
	// Since we have a single form with multiple groups, we need to get data from the current form
	// and create the initiative collection form
	currentForm := f.getCurrentForm()
	if currentForm == nil {
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}

	// Get data from current form
	characters := f.getCharacters()
	monsters := f.getMonsters()

	// If no characters or monsters selected, return an error or cancel
	if len(characters) == 0 && len(monsters) == 0 {
		// Could show an error here instead, but for now let's allow it
		return f.submit()
	}

	switch len(*f.forms) {
	case 1:
		var groups []*huh.Group
		var fields []huh.Field

		// Add character initiative fields
		for _, characterOption := range characters {
			key := fmt.Sprintf("initiative_%s", characterOption.id)
			title := fmt.Sprintf("%s's initiative", characterOption.creature.GetName())

			fields = append(
				fields,
				huh.NewInput().Key(key).Title(title).
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

		// Only create character group if there are character fields
		if len(fields) > 0 {
			groups = append(groups, huh.NewGroup(fields...))
		}

		for _, monsterOption := range monsters {
			name := monsterOption.creature.GetName()
			title := fmt.Sprintf("%s's quantity and initiative\n", name)

			groups = append(
				groups,
				huh.NewGroup(
					huh.NewInput().
						Key(fmt.Sprintf("name_%s", monsterOption.id)).
						Title("Name").
						Description("The name can be customized to represent the monster this statblock is attached to.").
						Value(&name),
					huh.NewInput().
						Key(fmt.Sprintf("quantity_%s", monsterOption.id)).
						Title("Quantity").
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
					huh.NewInput().
						Key(fmt.Sprintf("initiative_%s", monsterOption.id)).
						Title("Initiative").
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
				).Title(title),
			)
		}

		// If no groups were created, skip to submit
		if len(groups) == 0 {
			return f.submit()
		}

		form := huh.NewForm(groups...).WithShowErrors(true).WithShowHelp(false).WithAccessible(true).WithKeyMap(f.keys)

		form.CancelCmd = func() tea.Msg {
			return encounterFormPreviousStepMsg{}
		}

		form.SubmitCmd = func() tea.Msg {
			return encounterFormNextStepMsg{}
		}

		*f.forms = append(
			*f.forms,
			form,
		)

		// Apply sizing using the standard SetSize method
		f.SetSize(f.width, f.height)

		return form.Init()
	case 2:
		return f.submit()
	default:
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}
}

func (f *encounterForm) submit() tea.Cmd {
	if len(*f.forms) < 2 {
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}

	firstForm := (*f.forms)[0]
	secondForm := (*f.forms)[1]

	// Get values from first form
	summary := firstForm.GetString("summary")
	characterUUIDs, _ := firstForm.Get("characters").([]string)
	monsterValues, _ := firstForm.Get("monsters").([]string)

	var initiativeGroups []dnd.InitiativeGroup

	// Process character initiatives
	for _, uuid := range characterUUIDs {
		initiativeKey := fmt.Sprintf("initiative_%s", uuid)
		initiative := form.GetIntWithDefault(secondForm, initiativeKey, 0)

		if character, exists := (*f.characters)[uuid]; exists {
			group := dnd.InitiativeGroup{
				Initiative: initiative,
				Creatures:  []dnd.Creature{character},
			}
			initiativeGroups = append(initiativeGroups, group)
		}
	}

	// Process monster initiatives
	for _, value := range monsterValues {
		// value format is "sourceKey:monsterName"
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			continue
		}

		sourceKey := parts[0]
		monsterName := parts[1]

		// Get custom name, quantity, and initiative from second form
		name := secondForm.GetString(fmt.Sprintf("name_%s", value))
		quantity := form.GetIntWithDefault(secondForm, fmt.Sprintf("quantity_%s", value), 1)
		initiative := form.GetIntWithDefault(secondForm, fmt.Sprintf("initiative_%s", value), 0)

		// Find the monster in sources
		if source, exists := f.sources[sourceKey]; exists {
			for _, monster := range source.Monsters {
				if monster.GetName() == monsterName {
					// Create multiple monsters for this group
					var creatures []dnd.Creature
					for i := 0; i < quantity; i++ {
						monsterName := name
						if monsterName == "" {
							monsterName = monster.GetName()
						}

						newMonster := dnd.NewMonster(monsterName, monster.StatBlock)
						creatures = append(creatures, newMonster)
					}

					group := dnd.InitiativeGroup{
						Initiative: initiative,
						Creatures:  creatures,
					}
					initiativeGroups = append(initiativeGroups, group)
					break
				}
			}
		}
	}

	// Sort initiative groups by initiative value in descending order (highest first)
	sort.Slice(initiativeGroups, func(i, j int) bool {
		return initiativeGroups[i].Initiative > initiativeGroups[j].Initiative
	})

	return func() tea.Msg {
		return encounterFormSubmittedMsg{
			encounter: dnd.Encounter{
				Summary:          summary,
				StartedAt:        time.Now(),
				Round:            1,
				InitiativeGroups: initiativeGroups,
			},
		}
	}
}

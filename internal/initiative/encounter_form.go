package initiative

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

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

func (f *encounterForm) isInitiativeEditing() bool {
	return f.initiativeListField != nil && f.initiativeListField.delegate.IsEditing()
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

// HelpKeys returns the key bindings for use by the parent page
func (f *encounterForm) HelpKeys() []key.Binding {
	// When showing cancel confirmation, use its help keys
	if f.confirmCancelForm != nil {
		return f.confirmCancelForm.KeyBinds()
	}

	// When showing monster add form, use its help keys
	if f.monsterAddForm != nil {
		return f.monsterAddForm.HelpKeys()
	}

	// When showing initiative list, use its key bindings
	if f.initiativeListField != nil && len(*f.forms) > 1 {
		return f.initiativeListField.KeyBinds()
	}

	return f.getHelpKeys()
}

type encounterFormCancelledMsg struct{}
type encounterFormSubmittedMsg struct {
	encounter *dnd.Encounter
}
type encounterFormPreviousStepMsg struct{}
type encounterFormNextStepMsg struct{}

type encounterForm struct {
	forms  *[]*huh.Form
	styles encounterFormStyles
	keys   *huh.KeyMap

	characters map[string]*dnd.Character
	sources    map[string]*dnd.Source

	// Initiative list components
	initiativeListField *InitiativeListField
	monsterAddForm      *monsterAddForm
	editingMonsterIndex int // Index of monster being edited (-1 if adding new)

	// Cancel confirmation
	confirmCancelForm *huh.Form
	confirmCancel     bool

	width  int
	height int
	help   help.Model
}

func newEncounterForm(c map[string]*dnd.Character, s map[string]*dnd.Source, width int, height int) *encounterForm {
	// Create shared keymap
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	// Step 1: Summary input
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
	).WithTheme(currentTheme.form).WithShowErrors(true).WithShowHelp(false).WithKeyMap(keyMap)

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
			container: lipgloss.NewStyle().Padding(1, 2, 0, 2),
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},

		characters:          c,
		sources:             s,
		editingMonsterIndex: -1,
		help:                help.New(),
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

	// Handle cancel confirmation form if active
	if f.confirmCancelForm != nil {
		form, cmd := f.confirmCancelForm.Update(msg)
		f.confirmCancelForm = form.(*huh.Form)

		if f.confirmCancelForm.State == huh.StateCompleted {
			if f.confirmCancel {
				// User confirmed cancellation
				return f, tea.Cmd(func() tea.Msg {
					return encounterFormCancelledMsg{}
				})
			}
			// User chose not to cancel - dismiss the confirmation
			f.confirmCancelForm = nil
			f.confirmCancel = false

			// Recreate the current form to reset its aborted state while preserving field data
			if len(*f.forms) == 2 && f.initiativeListField != nil {
				form := huh.NewForm(
					huh.NewGroup(f.initiativeListField),
				).WithTheme(currentTheme.form).WithShowErrors(false).WithShowHelp(false).WithKeyMap(f.keys)

				form.CancelCmd = func() tea.Msg {
					return encounterFormPreviousStepMsg{}
				}
				form.SubmitCmd = func() tea.Msg {
					return encounterFormNextStepMsg{}
				}

				(*f.forms)[1] = form
				f.SetSize(f.width, f.height)
				return f, form.Init()
			}

			return f, nil
		}

		return f, cmd
	}

	// Handle monster add/edit form if active
	if f.monsterAddForm != nil {
		switch msg := msg.(type) {
		case monsterAddFormCancelledMsg:
			f.monsterAddForm = nil
			f.editingMonsterIndex = -1
			return f, nil
		case monsterAddFormSubmittedMsg:
			f.monsterAddForm = nil
			var cmd tea.Cmd
			if f.initiativeListField != nil {
				if f.editingMonsterIndex >= 0 {
					// Update existing monster
					f.initiativeListField.UpdateMonster(f.editingMonsterIndex, msg.entry)
				} else {
					// Add new monster
					cmd = f.initiativeListField.AddMonster(msg.entry)
				}
			}
			f.editingMonsterIndex = -1
			return f, cmd
		default:
			form, cmd := f.monsterAddForm.Update(msg)
			if form, ok := form.(*monsterAddForm); ok {
				f.monsterAddForm = form
			}
			return f, cmd
		}
	}

	// Handle request to add monster
	if _, ok := msg.(requestAddMonsterMsg); ok {
		f.editingMonsterIndex = -1
		f.monsterAddForm = newMonsterAddForm(f.sources, f.width, f.height)
		return f, f.monsterAddForm.Init()
	}

	// Handle request to edit monster
	if editMsg, ok := msg.(requestEditMonsterMsg); ok {
		f.editingMonsterIndex = editMsg.index
		f.monsterAddForm = newMonsterEditForm(f.sources, editMsg.entry, f.width, f.height)
		return f, f.monsterAddForm.Init()
	}

	switch msg.(type) {
	case encounterFormPreviousStepMsg:
		return f, tea.Sequence(append(cmds, f.prevStep())...)
	case encounterFormNextStepMsg:
		return f, tea.Sequence(append(cmds, f.nextStep())...)
	}

	// Check states that should disable Quit key before form update
	wasQuitDisabled := f.isFilteringActive(currentForm) || f.isInitiativeEditing()

	// Disable Quit key before form processes the message if needed
	if wasQuitDisabled {
		f.keys.Quit.SetEnabled(false)
	}

	form, cmd := currentForm.Update(msg)
	cmds = append(cmds, cmd)
	if form, ok := form.(*huh.Form); ok {
		(*f.forms)[len(*f.forms)-1] = form

		// Update Quit key state after form update
		isQuitDisabled := f.isFilteringActive(form) || f.isInitiativeEditing()
		if wasQuitDisabled != isQuitDisabled {
			f.keys.Quit.SetEnabled(!isQuitDisabled)
		}

		// Handle form quit (ESC pressed) - show confirmation
		if form.State == huh.StateAborted {
			return f, f.showCancelConfirmation()
		}
	}

	return f, tea.Sequence(cmds...)
}

func (f *encounterForm) View() string {
	// Show cancel confirmation if active
	if f.confirmCancelForm != nil {
		return f.styles.container.Render(f.confirmCancelForm.View())
	}

	// Show monster add form if active
	if f.monsterAddForm != nil {
		return f.monsterAddForm.View()
	}

	form := f.getCurrentForm()
	if form == nil {
		return ""
	}

	// Render form content only - no help component
	return f.styles.container.Render(form.View())
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

	// Update initiative list field size if it exists
	if f.initiativeListField != nil {
		f.initiativeListField.WithWidth(formWidth).WithHeight(formHeight)
	}

	// Update monster add form size if it exists
	if f.monsterAddForm != nil {
		f.monsterAddForm.SetSize(width, height)
	}

	// Update cancel confirmation form size if it exists
	if f.confirmCancelForm != nil {
		f.confirmCancelForm.WithWidth(formWidth).WithHeight(formHeight)
	}
}

func (f *encounterForm) prevStep() tea.Cmd {
	// when more than 1 form, remove the last form
	if len(*f.forms) > 1 {
		*f.forms = (*f.forms)[:len(*f.forms)-1]
		// Clear initiative list field when going back
		if len(*f.forms) == 1 {
			f.initiativeListField = nil
		}
		return nil
	}

	// when there's just 1 form, show cancel confirmation
	return f.showCancelConfirmation()
}

func (f *encounterForm) showCancelConfirmation() tea.Cmd {
	f.confirmCancel = false
	f.confirmCancelForm = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to cancel?").
				Affirmative("Yes").
				Negative("No").
				Value(&f.confirmCancel),
		),
	).WithTheme(currentTheme.form).WithWidth(f.width).WithHeight(f.height)

	return f.confirmCancelForm.Init()
}

func (f *encounterForm) nextStep() tea.Cmd {
	currentForm := f.getCurrentForm()
	if currentForm == nil {
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}

	switch len(*f.forms) {
	case 1:
		// Create initiative list field with all characters pre-populated
		f.initiativeListField = NewInitiativeListField(f.characters, f.sources)

		// Create form with the initiative list field (no title - list shows its own status)
		// ShowErrors is false because the field shows its own status line
		form := huh.NewForm(
			huh.NewGroup(f.initiativeListField),
		).WithTheme(currentTheme.form).WithShowErrors(false).WithShowHelp(false).WithKeyMap(f.keys)

		form.CancelCmd = func() tea.Msg {
			return encounterFormPreviousStepMsg{}
		}

		form.SubmitCmd = func() tea.Msg {
			return encounterFormNextStepMsg{}
		}

		*f.forms = append(*f.forms, form)

		// Apply sizing
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
	if len(*f.forms) < 2 || f.initiativeListField == nil {
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}

	// Validate all entries have initiative
	if err := f.initiativeListField.Error(); err != nil {
		// Don't submit if validation fails
		return nil
	}

	firstForm := (*f.forms)[0]
	summary := firstForm.GetString("summary")

	// Get entries from initiative list field
	entries, ok := f.initiativeListField.GetValue().([]dnd.InitiativeEntry)
	if !ok {
		return func() tea.Msg {
			return encounterFormCancelledMsg{}
		}
	}

	var initiativeGroups []*dnd.InitiativeGroup

	// Convert entries to initiative groups
	for _, entry := range entries {
		if entry.CreatureType == "character" {
			// Find the character
			if character, exists := f.characters[entry.CreatureID]; exists {
				var creature dnd.Creature = character
				group := dnd.NewInitiativeGroup(entry.Initiative, []*dnd.Creature{&creature})
				initiativeGroups = append(initiativeGroups, group)
			}
		} else {
			// Monster entry
			if entry.StatBlock != nil {
				var creatures []*dnd.Creature
				for i := 0; i < entry.Quantity; i++ {
					monster := dnd.NewMonster(entry.Name, *entry.StatBlock)
					var creature dnd.Creature = monster
					creatures = append(creatures, &creature)
				}
				group := dnd.NewInitiativeGroup(entry.Initiative, creatures)
				initiativeGroups = append(initiativeGroups, group)
			}
		}
	}

	// Sort initiative groups by initiative value in descending order (highest first)
	sort.Slice(initiativeGroups, func(i, j int) bool {
		return initiativeGroups[i].Initiative() > initiativeGroups[j].Initiative()
	})

	return func() tea.Msg {
		encounter := dnd.NewEncounter(summary, time.Now(), initiativeGroups)
		return encounterFormSubmittedMsg{
			encounter: encounter,
		}
	}
}

package initiative

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

type monsterAddFormStyles struct {
	container lipgloss.Style
	help      lipgloss.Style
}

type monsterAddFormCancelledMsg struct{}

type monsterAddFormSubmittedMsg struct {
	entry dnd.InitiativeEntry
}

type monsterAddForm struct {
	form   *huh.Form
	styles monsterAddFormStyles

	sources map[string]*dnd.Source
	width   int
	height  int
	help    help.Model
}

func newMonsterAddForm(sources map[string]*dnd.Source, width int, height int) *monsterAddForm {
	var monsterSelection string
	var customName string
	var quantity string

	// Build monster options from all sources
	options := buildMonsterOptions(sources)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("monster").
				Title("Monster").
				Options(options...).
				Value(&monsterSelection).
				Height(8). // Limit height of dropdown
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Please select a monster")
					}
					return nil
				}),

			huh.NewInput().
				Key("quantity").
				Title("Quantity").
				Description("e.g. 3, 1d4, 2d6").
				Placeholder("1").
				Value(&quantity).
				Validate(func(str string) error {
					str = strings.TrimSpace(str)
					if str == "" {
						return nil // Will default to 1
					}
					_, err := parseAdjustment(str)
					return err
				}),

			huh.NewInput().
				Key("customName").
				Title("Custom Name").
				Description("Optional - leave empty to use monster name").
				Value(&customName),
		).Title("Add Monster\n"),
	).WithTheme(currentTheme.form).WithShowErrors(true).WithShowHelp(false).WithKeyMap(customMonsterAddFormKeyMap())

	f := &monsterAddForm{
		form:    form,
		sources: sources,
		styles: monsterAddFormStyles{
			container: lipgloss.NewStyle(),
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},
		width:  width,
		height: height,
		help:   help.New(),
	}

	f.SetSize(width, height)
	return f
}

func newMonsterEditForm(sources map[string]*dnd.Source, entry dnd.InitiativeEntry, width int, height int) *monsterAddForm {
	monsterSelection := entry.CreatureID
	quantity := fmt.Sprintf("%d", entry.Quantity)

	// Only set custom name if it differs from the monster's default name
	var customName string
	parts := strings.SplitN(entry.CreatureID, ":", 2)
	if len(parts) == 2 {
		sourceKey, monsterName := parts[0], parts[1]
		if source, exists := sources[sourceKey]; exists {
			for _, monster := range source.Monsters {
				if monster.Name() == monsterName && entry.Name != monster.Name() {
					customName = entry.Name
					break
				}
			}
		}
	}

	// Build monster options from all sources
	options := buildMonsterOptions(sources)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("monster").
				Title("Monster").
				Options(options...).
				Value(&monsterSelection).
				Height(8). // Limit height of dropdown
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("Please select a monster")
					}
					return nil
				}),

			huh.NewInput().
				Key("quantity").
				Title("Quantity").
				Description("e.g. 3, 1d4, 2d6").
				Placeholder("1").
				Value(&quantity).
				Validate(func(str string) error {
					str = strings.TrimSpace(str)
					if str == "" {
						return nil // Will default to 1
					}
					_, err := parseAdjustment(str)
					return err
				}),

			huh.NewInput().
				Key("customName").
				Title("Custom Name").
				Description("Optional - leave empty to use monster name").
				Value(&customName),
		).Title("Edit Monster\n"),
	).WithTheme(currentTheme.form).WithShowErrors(true).WithShowHelp(false).WithKeyMap(customMonsterAddFormKeyMap())

	f := &monsterAddForm{
		form:    form,
		sources: sources,
		styles: monsterAddFormStyles{
			container: lipgloss.NewStyle(),
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},
		width:  width,
		height: height,
		help:   help.New(),
	}

	f.SetSize(width, height)
	return f
}

func buildMonsterOptions(sources map[string]*dnd.Source) []huh.Option[string] {
	type monsterOption struct {
		label string
		value string
	}

	var options []monsterOption
	for sourceKey, source := range sources {
		for _, monster := range source.Monsters {
			value := sourceKey + ":" + monster.Name()
			options = append(options, monsterOption{
				label: monster.Name(),
				value: value,
			})
		}
	}

	// Sort by label
	sort.Slice(options, func(i, j int) bool {
		return strings.ToLower(options[i].label) < strings.ToLower(options[j].label)
	})

	result := make([]huh.Option[string], len(options))
	for i, opt := range options {
		result[i] = huh.NewOption(opt.label, opt.value)
	}

	return result
}

func customMonsterAddFormKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()
	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	return keyMap
}

func (f *monsterAddForm) getHelpKeys() []key.Binding {
	if f.form == nil {
		return []key.Binding{}
	}

	// Get form's context-sensitive keys (up/down navigation, filtering)
	formKeys := f.form.KeyBinds()

	// Add standard form navigation keys
	enterKey := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "next"),
	)
	escKey := key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	filterKey := key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	)

	// Check if focused field is a Select (supports filtering)
	focusedField := f.form.GetFocusedField()
	_, isSelect := focusedField.(*huh.Select[string])

	keys := append(formKeys, enterKey, escKey)
	if isSelect {
		keys = append(keys, filterKey)
	}

	return keys
}

// HelpKeys returns the key bindings for use by the parent page
func (f *monsterAddForm) HelpKeys() []key.Binding {
	return f.getHelpKeys()
}

func (f *monsterAddForm) Init() tea.Cmd {
	return f.form.Init()
}

func (f *monsterAddForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		f.form = form
	}

	if f.form.State == huh.StateAborted {
		return f, tea.Cmd(func() tea.Msg {
			return monsterAddFormCancelledMsg{}
		})
	}

	if f.form.State == huh.StateCompleted {
		return f, f.submit()
	}

	return f, cmd
}

func (f *monsterAddForm) submit() tea.Cmd {
	monsterSelection := f.form.GetString("monster")
	customName := strings.TrimSpace(f.form.GetString("customName"))
	quantityStr := strings.TrimSpace(f.form.GetString("quantity"))

	quantity := 1
	if quantityStr != "" {
		if val, err := parseAdjustment(quantityStr); err == nil {
			quantity = val
		}
	}

	// Parse monster selection (format: "sourceKey:monsterName")
	parts := strings.SplitN(monsterSelection, ":", 2)
	if len(parts) != 2 {
		return tea.Cmd(func() tea.Msg {
			return monsterAddFormCancelledMsg{}
		})
	}

	sourceKey := parts[0]
	monsterName := parts[1]

	// Find the monster in sources
	var statBlock *dnd.StatBlock
	var displayName string

	if source, exists := f.sources[sourceKey]; exists {
		for _, monster := range source.Monsters {
			if monster.Name() == monsterName {
				sb := monster.StatBlock()
				statBlock = &sb
				displayName = monster.Name()
				break
			}
		}
	}

	if statBlock == nil {
		return tea.Cmd(func() tea.Msg {
			return monsterAddFormCancelledMsg{}
		})
	}

	// Use custom name if provided
	if customName != "" {
		displayName = customName
	}

	entry := dnd.InitiativeEntry{
		CreatureType: "monster",
		CreatureID:   monsterSelection,
		Name:         displayName,
		Quantity:     quantity,
		StatBlock:    statBlock,
	}

	return tea.Cmd(func() tea.Msg {
		return monsterAddFormSubmittedMsg{entry: entry}
	})
}

func (f *monsterAddForm) View() string {
	if f.form == nil {
		return ""
	}

	return f.styles.container.Render(f.form.View())
}

func (f *monsterAddForm) SetSize(width int, height int) {
	f.width = width
	f.height = height

	if f.form == nil {
		return
	}

	container := f.styles.container
	horizontalPadding := container.GetHorizontalPadding()
	verticalPadding := container.GetVerticalPadding()

	f.help.Width = width - f.styles.help.GetHorizontalPadding()
	helpKeys := f.getHelpKeys()
	sampleHelp := f.styles.help.Render(f.help.ShortHelpView(helpKeys))
	helpHeight := lipgloss.Height(sampleHelp)

	formWidth := width - horizontalPadding
	formHeight := height - verticalPadding - helpHeight

	f.form.WithHeight(formHeight).WithWidth(formWidth)
}

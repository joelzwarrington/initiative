package initiative

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

type characterFormStyles struct {
	container lipgloss.Style
	help      lipgloss.Style
}

type characterFormCancelledMsg struct{}
type characterFormSubmittedMsg struct {
	uuid      string
	character *dnd.Character
}

type characterForm struct {
	form   *huh.Form
	styles characterFormStyles

	uuid   string // empty for new character, populated for edit
	width  int
	height int
	help   help.Model
}

func newCharacterForm(uuid string, existingCharacter *dnd.Character, width int, height int) *characterForm {
	title := "Add character\n"
	var name string
	var maxHP string

	if uuid != "" && existingCharacter != nil {
		name = existingCharacter.Name()
		title = "Edit " + name + "\n"
		if existingCharacter.MaxHP() > 0 {
			maxHP = fmt.Sprintf("%d", existingCharacter.MaxHP())
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Value(&name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return fmt.Errorf("Name is required")
					}
					return nil
				}),
			huh.NewInput().
				Key("max_hp").
				Title("Maximum Hit Points").
				Description("Leave empty for 0 HP (no health bar)").
				Value(&maxHP).
				Validate(func(str string) error {
					trimmed := strings.TrimSpace(str)
					if trimmed == "" {
						return nil // Allow empty for 0 HP
					}
					value, err := strconv.Atoi(trimmed)
					if err != nil || value < 0 {
						return fmt.Errorf("Maximum Hit Points must be a non-negative number")
					}
					return nil
				}),
		).Title(title),
	).WithShowErrors(true).WithShowHelp(false).WithAccessible(true).WithKeyMap(customCharacterFormKeyMap())

	f := &characterForm{
		form: form,
		uuid: uuid,
		styles: characterFormStyles{
			container: lipgloss.NewStyle().Padding(1, 2, 0, 2), // 1 top, 2 horizontal, 0 bottom
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},
		width:  width,
		height: height,
		help:   help.New(),
	}

	f.SetSize(width, height)
	return f
}

func customCharacterFormKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()

	// Set ESC to quit the form
	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	return keyMap
}

func (f *characterForm) getHelpKeys() []key.Binding {
	if f.form == nil {
		return []key.Binding{}
	}

	// Get form's key bindings
	formKeys := f.form.KeyBinds()

	// Add our custom ESC binding
	escKey := key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	// Combine form keys with our custom keys
	allKeys := append(formKeys, escKey)
	return allKeys
}

func (f *characterForm) Init() tea.Cmd {
	return f.form.Init()
}

func (f *characterForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		f.form = form
	}

	// Handle form quit (ESC pressed)
	if f.form.State == huh.StateAborted {
		return f, tea.Cmd(func() tea.Msg {
			return characterFormCancelledMsg{}
		})
	}

	// Handle form completion
	if f.form.State == huh.StateCompleted {
		name := f.form.GetString("name")
		maxHPStr := strings.TrimSpace(f.form.GetString("max_hp"))

		var maxHP int
		if maxHPStr != "" {
			var err error
			maxHP, err = strconv.Atoi(maxHPStr)
			if err != nil || maxHP < 0 {
				maxHP = 0
			}
		}

		var character *dnd.Character
		if maxHP > 0 {
			character = dnd.NewCharacter(name).WithHealth(maxHP, maxHP)
		} else {
			character = dnd.NewCharacter(name)
		}

		return f, tea.Cmd(func() tea.Msg {
			return characterFormSubmittedMsg{
				uuid:      f.uuid,
				character: character,
			}
		})
	}

	return f, cmd
}

func (f *characterForm) View() string {
	if f.form == nil {
		return ""
	}

	// Render form content
	formView := f.styles.container.Render(f.form.View())

	// Render custom help
	f.help.Width = f.width - f.styles.help.GetHorizontalPadding()
	helpKeys := f.getHelpKeys()
	helpView := f.styles.help.Render(f.help.ShortHelpView(helpKeys))

	// Join form and help vertically
	return lipgloss.JoinVertical(lipgloss.Left, formView, helpView)
}

func (f *characterForm) SetSize(width int, height int) {
	f.width = width
	f.height = height

	if f.form == nil {
		return
	}

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

	f.form.WithHeight(formHeight).WithWidth(formWidth)
}

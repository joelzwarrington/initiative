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
)

type hitPointFormStyles struct {
	container lipgloss.Style
	help      lipgloss.Style
}

type HitPointAdjustmentType int

const (
	HitPointDamage HitPointAdjustmentType = iota
	HitPointHeal
)

func (h HitPointAdjustmentType) String() string {
	switch h {
	case HitPointDamage:
		return "damage"
	case HitPointHeal:
		return "heal"
	default:
		return "damage"
	}
}

type hitPointFormCancelledMsg struct{}

type hitPointFormSubmittedMsg struct {
	groupIndex     int
	creatureIndex  int
	adjustment     int
	adjustmentType HitPointAdjustmentType
}

type hitPointForm struct {
	form   *huh.Form
	styles hitPointFormStyles

	groupIndex    int
	creatureIndex int
	creatureName  string
	width         int
	height        int
	help          help.Model
}

func newHitPointForm(groupIndex, creatureIndex int, creatureName string, width int, height int, initialIsDamage bool) *hitPointForm {
	var adjustment string
	var isDamage bool = initialIsDamage

	title := fmt.Sprintf("Adjust %s's hit points\n", creatureName)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("adjustment").
				Title("HP adjustment").
				Value(&adjustment).
				Validate(func(str string) error {
					str = strings.TrimSpace(str)
					if str == "" {
						return fmt.Errorf("Amount is required")
					}
					val, err := strconv.Atoi(str)
					if err != nil {
						return fmt.Errorf("Amount must be a valid number")
					}
					if val <= 0 {
						return fmt.Errorf("Amount must be a positive number")
					}
					return nil
				}),

			huh.NewConfirm().
				Key("isDamage").
				Title("Type").
				Affirmative("Damage").
				Negative("Heal").
				Value(&isDamage),
		).Title(title),
	).WithShowErrors(true).WithShowHelp(false).WithAccessible(true).WithKeyMap(customHitPointFormKeyMap())

	f := &hitPointForm{
		form:          form,
		groupIndex:    groupIndex,
		creatureIndex: creatureIndex,
		creatureName:  creatureName,
		styles: hitPointFormStyles{
			container: lipgloss.NewStyle().Padding(1, 2, 0, 2),
			help:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 2),
		},
		width:  width,
		height: height,
		help:   help.New(),
	}

	f.SetSize(width, height)
	return f
}

func customHitPointFormKeyMap() *huh.KeyMap {
	keyMap := huh.NewDefaultKeyMap()

	keyMap.Quit = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	return keyMap
}

func (f *hitPointForm) getHelpKeys() []key.Binding {
	if f.form == nil {
		return []key.Binding{}
	}

	formKeys := f.form.KeyBinds()

	escKey := key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	allKeys := append(formKeys, escKey)
	return allKeys
}

// HelpKeys returns the key bindings for use by the parent page
func (f *hitPointForm) HelpKeys() []key.Binding {
	return f.getHelpKeys()
}

func (f *hitPointForm) Init() tea.Cmd {
	return f.form.Init()
}

func (f *hitPointForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		f.form = form
	}

	if f.form.State == huh.StateAborted {
		return f, tea.Cmd(func() tea.Msg {
			return hitPointFormCancelledMsg{}
		})
	}

	if f.form.State == huh.StateCompleted {
		adjustmentStr := f.form.GetString("adjustment")
		isDamage := f.form.GetBool("isDamage")

		adjustment, _ := strconv.Atoi(adjustmentStr)

		var adjustmentType HitPointAdjustmentType
		if isDamage {
			adjustmentType = HitPointDamage
		} else {
			adjustmentType = HitPointHeal
		}

		return f, tea.Cmd(func() tea.Msg {
			return hitPointFormSubmittedMsg{
				groupIndex:     f.groupIndex,
				creatureIndex:  f.creatureIndex,
				adjustment:     adjustment,
				adjustmentType: adjustmentType,
			}
		})
	}

	return f, cmd
}

func (f *hitPointForm) View() string {
	if f.form == nil {
		return ""
	}

	// Render form content only - no help component
	return f.styles.container.Render(f.form.View())
}

func (f *hitPointForm) SetSize(width int, height int) {
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

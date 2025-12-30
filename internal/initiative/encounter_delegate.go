package initiative

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
	"github.com/joelzwarrington/initiative/internal/components"
)

type encounterDelegate struct {
	list      list.Model
	healthBar components.HealthBar
}

func newEncounterDelegate(width, height int) *encounterDelegate {
	delegate := &creatureItemDelegate{
		healthBar: components.NewHealthBar(30),
	}

	initiativeList := list.New([]list.Item{}, delegate, width, height)
	initiativeList.SetStatusBarItemName("creature", "creatures")
	initiativeList.SetShowTitle(false)
	initiativeList.SetShowStatusBar(true)
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()

	return &encounterDelegate{
		list:      initiativeList,
		healthBar: components.NewHealthBar(30),
	}
}

func (d *encounterDelegate) Render(w io.Writer, encounter *dnd.Encounter) {
	if encounter == nil {
		return
	}

	items := []list.Item{}
	for groupIndex, group := range encounter.InitiativeGroups {
		isCurrentTurnGroup := groupIndex == encounter.GetTurnIndex()

		for creatureIndex, creature := range group.Creatures {
			items = append(items, creatureItem{
				creature:      creature,
				initiative:    group.Initiative,
				groupIndex:    groupIndex,
				creatureIndex: creatureIndex,
				totalInGroup:  len(group.Creatures),
				isCurrentTurn: isCurrentTurnGroup,
			})
		}
	}

	d.list.SetItems(items)
	fmt.Fprint(w, d.list.View())
}

func (d *encounterDelegate) Update(msg tea.Msg, encounter *dnd.Encounter) tea.Cmd {
	var cmd tea.Cmd
	d.list, cmd = d.list.Update(msg)
	return cmd
}

func (d *encounterDelegate) SetSize(width, height int) {
	d.list.SetWidth(width)
	d.list.SetHeight(height)
}

func (d *encounterDelegate) ShortHelp() []key.Binding {
	// Get list navigation keys
	listKeys := []key.Binding{
		d.list.KeyMap.CursorUp,
		d.list.KeyMap.CursorDown,
	}

	// Add list-specific keys
	if d.list.FilterState() != list.Filtering {
		listKeys = append(listKeys, d.list.KeyMap.Filter)
	}

	return listKeys
}

func (d *encounterDelegate) FullHelp() [][]key.Binding {
	// Get list navigation keys
	navKeys := []key.Binding{
		d.list.KeyMap.CursorUp,
		d.list.KeyMap.CursorDown,
		d.list.KeyMap.Filter,
	}

	return [][]key.Binding{navKeys}
}

// List item for individual creatures
type creatureItem struct {
	creature      dnd.Creature
	initiative    int
	groupIndex    int
	creatureIndex int
	totalInGroup  int
	isCurrentTurn bool
}

func (c creatureItem) FilterValue() string {
	return c.creature.GetName()
}

// List delegate for individual creatures
type creatureItemDelegate struct {
	healthBar components.HealthBar
}

func (d creatureItemDelegate) Height() int {
	return 3
}

func (d creatureItemDelegate) Spacing() int {
	return 0
}

func (d *creatureItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d creatureItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	creature, ok := listItem.(creatureItem)
	if !ok {
		return
	}

	isFirst := creature.creatureIndex == 0
	isLast := creature.creatureIndex == creature.totalInGroup-1

	turnAwareForeground := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	if creature.isCurrentTurn {
		turnAwareForeground = lipgloss.AdaptiveColor{Light: "#00FF00", Dark: "#46FF46"}
	}

	prefix := " "
	if m.Index() == index {
		prefix = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render(">")
	}

	if isFirst {
		initiative := lipgloss.NewStyle().Foreground(turnAwareForeground).Render(fmt.Sprintf("%2d", creature.initiative))
		separator := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}).Render("•")

		prefix = fmt.Sprintf("%s %s %s ", prefix, initiative, separator)
	} else {
		prefix = prefix + strings.Repeat(" ", 6)
	}

	health := ""
	if monster, ok := creature.creature.(dnd.Monster); ok {
		health = " " + d.healthBar.View(monster.HitPoints, monster.MaximumHitPoints)
	}

	top := 0
	if !isFirst {
		top = 1
	}

	content := lipgloss.NewStyle().
		Padding(top, 1, 0, 1).
		Render(prefix + creature.creature.GetName() + health)

	top = 0
	if isFirst && index > 0 {
		top = 1
	}

	border := lipgloss.NewStyle().
		Margin(top, 0, 0, 0).
		Border(lipgloss.RoundedBorder(), isFirst, true, isLast, true).
		BorderForeground(turnAwareForeground).
		Width(m.Width() - 2)

	fmt.Fprint(w, border.Render(content))
}

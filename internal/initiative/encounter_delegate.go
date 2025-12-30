package initiative

import (
	"fmt"
	"io"
	"strings"
	"time"

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
		keys:      newCreatureItemKeyMap(),
	}

	initiativeList := list.New([]list.Item{}, delegate, width, height)
	initiativeList.SetStatusBarItemName("creature", "creatures")
	initiativeList.SetShowTitle(false)
	initiativeList.SetShowStatusBar(true)
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()
	initiativeList.StatusMessageLifetime = time.Duration(3000) * time.Millisecond
	initiativeList.StatusMessageLocation = list.InStatusBar

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

	// Add creature action keys if items exist
	if len(d.list.Items()) > 0 {
		delegate := newCreatureItemKeyMap()
		listKeys = append(listKeys, delegate.dealDamage, delegate.heal)
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

	// Add creature action keys if items exist
	actionKeys := []key.Binding{}
	if len(d.list.Items()) > 0 {
		delegate := newCreatureItemKeyMap()
		actionKeys = append(actionKeys, delegate.dealDamage, delegate.heal)
	}

	if len(actionKeys) > 0 {
		return [][]key.Binding{navKeys, actionKeys}
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

// Messages
type adjustHitPointsMsg struct {
	groupIndex    int
	creatureIndex int
	isDamage      bool
}

// List delegate for individual creatures
type creatureItemDelegate struct {
	healthBar components.HealthBar
	keys      creatureItemKeyMap
}

func (d creatureItemDelegate) Height() int {
	return 4
}

func (d creatureItemDelegate) Spacing() int {
	return 0
}

func (d *creatureItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	item, ok := m.SelectedItem().(creatureItem)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.dealDamage):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return adjustHitPointsMsg{
						groupIndex:    item.groupIndex,
						creatureIndex: item.creatureIndex,
						isDamage:      true,
					}
				})
			}
		case key.Matches(msg, d.keys.heal):
			if ok {
				return tea.Cmd(func() tea.Msg {
					return adjustHitPointsMsg{
						groupIndex:    item.groupIndex,
						creatureIndex: item.creatureIndex,
						isDamage:      false,
					}
				})
			}
		}
	}

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

// Key mappings for creature items
type creatureItemKeyMap struct {
	dealDamage key.Binding
	heal       key.Binding
}

func newCreatureItemKeyMap() creatureItemKeyMap {
	return creatureItemKeyMap{
		dealDamage: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "deal damage"),
		),
		heal: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "heal"),
		),
	}
}

func (d creatureItemKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.dealDamage,
		d.heal,
	}
}

func (d creatureItemKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.dealDamage,
			d.heal,
		},
	}
}

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
	list         list.Model
	healthBar    components.HealthBar
	creatureKeys creatureItemKeyMap
}

func newEncounterDelegate(width, height int) *encounterDelegate {
	creatureKeys := newCreatureItemKeyMap()
	delegate := &creatureItemDelegate{
		healthBar: components.NewHealthBar(30),
		keys:      creatureKeys,
	}

	initiativeList := list.New([]list.Item{}, delegate, width, height)
	initiativeList.SetStatusBarItemName("creature", "creatures")
	initiativeList.Title = "Encounter"
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()
	initiativeList.StatusMessageLifetime = time.Duration(3000) * time.Millisecond
	initiativeList.StatusMessageLocation = list.InStatusBar
	initiativeList.Styles = currentTheme.list

	// Customize keymap to avoid conflicts with our creature actions
	keyMap := list.DefaultKeyMap()
	// Update up/down help text to just show arrows
	keyMap.CursorUp = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "up"),
	)
	keyMap.CursorDown = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "down"),
	)
	// Remove 'd' from NextPage to avoid conflict with dealDamage
	keyMap.NextPage = key.NewBinding(
		key.WithKeys("ctrl+right", "l", "pgdown", "f"),
		key.WithHelp("ctrl+→", "next page"),
	)
	// Remove 'h' from PrevPage to avoid conflict with heal
	keyMap.PrevPage = key.NewBinding(
		key.WithKeys("ctrl+left", "pgup", "b", "u"),
		key.WithHelp("ctrl+←", "prev page"),
	)
	initiativeList.KeyMap = keyMap

	return &encounterDelegate{
		list:         initiativeList,
		healthBar:    components.NewHealthBar(30),
		creatureKeys: creatureKeys,
	}
}

func (d *encounterDelegate) Render(w io.Writer, encounter *dnd.Encounter) {
	if encounter == nil {
		return
	}
	d.list.Title = encounter.Summary()

	items := []list.Item{}
	for groupIndex, group := range encounter.InitiativeGroups() {
		isCurrentTurnGroup := groupIndex == encounter.TurnIndex()

		for creatureIndex, creature := range group.Creatures() {
			items = append(items, creatureItem{
				creature:      *creature,
				initiative:    group.Initiative(),
				groupIndex:    groupIndex,
				creatureIndex: creatureIndex,
				totalInGroup:  len(group.Creatures()),
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
	// Update key states
	d.updateKeyStates()

	// Get list navigation keys (excluding filter for short help)
	keys := []key.Binding{
		d.list.KeyMap.CursorUp,
		d.list.KeyMap.CursorDown,
		d.list.KeyMap.PrevPage,
		d.list.KeyMap.NextPage,
		// Note: Filter excluded from short help to keep it concise
	}

	// Add creature action keys
	keys = append(keys, d.creatureKeys.dealDamage, d.creatureKeys.heal)

	return keys
}

func (d *encounterDelegate) FullHelp() [][]key.Binding {
	// Update key states
	d.updateKeyStates()

	// Organize keys into logical groups
	navigationKeys := []key.Binding{
		d.list.KeyMap.CursorUp,
		d.list.KeyMap.CursorDown,
		d.list.KeyMap.PrevPage,
		d.list.KeyMap.NextPage,
	}

	filterKeys := []key.Binding{
		d.list.KeyMap.Filter,
	}

	result := [][]key.Binding{navigationKeys, filterKeys}

	// Add creature action keys as separate column
	actionKeys := []key.Binding{
		d.creatureKeys.dealDamage,
		d.creatureKeys.heal,
	}
	result = append(result, actionKeys)

	return result
}

func (d *encounterDelegate) updateKeyStates() {
	hasItems := len(d.list.Items()) > 0
	isFiltering := d.list.FilterState() == list.Filtering

	// Navigation keys should only be enabled when there are items
	d.list.KeyMap.CursorUp.SetEnabled(hasItems)
	d.list.KeyMap.CursorDown.SetEnabled(hasItems)
	d.list.KeyMap.PrevPage.SetEnabled(hasItems)
	d.list.KeyMap.NextPage.SetEnabled(hasItems)

	// Filter key enabled when not currently filtering and has items
	d.list.KeyMap.Filter.SetEnabled(hasItems && !isFiltering)

	// Check if selected creature has health (max HP > 0)
	selectedHasHealth := false
	if hasItems && !isFiltering {
		if selectedItem := d.list.SelectedItem(); selectedItem != nil {
			if creatureItem, ok := selectedItem.(creatureItem); ok {
				selectedHasHealth = creatureItem.creature.MaxHP() > 0
			}
		}
	}

	// Creature actions only when items exist, not filtering, and creature has health
	creatureActionsEnabled := hasItems && !isFiltering && selectedHasHealth
	d.creatureKeys.dealDamage.SetEnabled(creatureActionsEnabled)
	d.creatureKeys.heal.SetEnabled(creatureActionsEnabled)
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
	return c.creature.Name()
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
	return 2
}

func (d creatureItemDelegate) Spacing() int {
	return 1
}

func (d *creatureItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	item, ok := m.SelectedItem().(creatureItem)

	// Only allow damage/heal actions if creature has health
	hasHealth := ok && item.creature.MaxHP() > 0

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.dealDamage):
			if hasHealth {
				return tea.Cmd(func() tea.Msg {
					return adjustHitPointsMsg{
						groupIndex:    item.groupIndex,
						creatureIndex: item.creatureIndex,
						isDamage:      true,
					}
				})
			}
		case key.Matches(msg, d.keys.heal):
			if hasHealth {
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
	item, ok := listItem.(creatureItem)
	if !ok {
		return
	}

	turnAwareForeground := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	if item.isCurrentTurn {
		turnAwareForeground = lipgloss.AdaptiveColor{Light: "#00FF00", Dark: "#46FF46"}
	}

	spacer := " "
	seperator := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}).Render(" • ")

	selection := " "
	if m.Index() == index {
		selection = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render(">")
	}

	initiative := lipgloss.NewStyle().Foreground(turnAwareForeground).Render(fmt.Sprintf("%2d •", item.initiative))

	prefix := lipgloss.JoinHorizontal(
		lipgloss.Left,
		spacer, selection, spacer, initiative, spacer,
	)

	ac := ""
	if item.creature.AC() > 0 {
		ac = seperator + icons.ArmorClass.Join(fmt.Sprintf("%d", item.creature.AC()), nil)
	}

	health := ""
	if item.creature.MaxHP() > 0 {
		health = d.healthBar.View(item.creature.HP(), item.creature.MaxHP())
	}

	content :=
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				prefix,
				item.creature.Name(),
				ac,
			),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				strings.Repeat(spacer, lipgloss.Width(prefix)),
				health,
			),
		)

	fmt.Fprint(w, content)
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
			key.WithHelp("d", "damage"),
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

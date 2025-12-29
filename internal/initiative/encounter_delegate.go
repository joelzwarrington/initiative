

package initiative

import (
	"fmt"
	"initiative/dnd"
	"initiative/internal/components"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type encounterDelegate struct {
	list      list.Model
	healthBar components.HealthBar
}

func newEncounterDelegate(width, height int) *encounterDelegate {
	delegate := &initiativeGroupItemDelegate{
		healthBar: components.NewHealthBar(30),
	}
	
	initiativeList := list.New([]list.Item{}, delegate, width, height)
	initiativeList.SetStatusBarItemName("initiative group", "initiative groups")
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
	for i, group := range encounter.InitiativeGroups {
		items = append(items, initiativeGroupItem{
			group:         group,
			isCurrentTurn: i == encounter.GetTurnIndex(),
		})
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

// List item for initiative groups
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
	healthBar components.HealthBar
}

func (d initiativeGroupItemDelegate) Height() int {
	return 1
}

func (d initiativeGroupItemDelegate) Spacing() int { 
	return 1 
}

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
			healthBar = " " + d.healthBar.View(monster.HitPoints, monster.MaximumHitPoints)
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
				healthBar := d.healthBar.View(monster.HitPoints, monster.MaximumHitPoints)
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
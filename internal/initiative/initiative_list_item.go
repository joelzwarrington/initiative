package initiative

import (
	"fmt"
	"io"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

// initiativeItem represents a creature in the initiative list
type initiativeItem struct {
	entry dnd.InitiativeEntry
}

// FilterValue implements list.Item interface
func (i initiativeItem) FilterValue() string {
	return i.entry.Name
}

// initiativeItemDelegate handles rendering and key handling for initiative items
type initiativeItemDelegate struct {
	textInput    textinput.Model
	editingIndex int // Index of item being edited (-1 if none)
	styles       initiativeItemStyles
}

type initiativeItemStyles struct {
	selected   lipgloss.Style
	normal     lipgloss.Style
	initiative lipgloss.Style
	noInit     lipgloss.Style
	quantity   lipgloss.Style
	inputStyle lipgloss.Style
}

func newInitiativeItemDelegate() *initiativeItemDelegate {
	ti := textinput.New()
	ti.Width = 2
	ti.CharLimit = 2
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"})

	return &initiativeItemDelegate{
		textInput:    ti,
		editingIndex: -1,
		styles: initiativeItemStyles{
			selected:   lipgloss.NewStyle().Foreground(currentTheme.colors.accent),
			normal:     lipgloss.NewStyle(),
			initiative: lipgloss.NewStyle().Foreground(currentTheme.colors.success),
			noInit:     lipgloss.NewStyle().Foreground(currentTheme.colors.error),
			quantity:   lipgloss.NewStyle().Foreground(currentTheme.colors.muted),
			inputStyle: lipgloss.NewStyle().Background(lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#3C3C3C"}),
		},
	}
}

func (d *initiativeItemDelegate) Height() int {
	return 1
}

func (d *initiativeItemDelegate) Spacing() int {
	return 0
}

func (d *initiativeItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	// Handle text input updates when editing
	if d.editingIndex >= 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			keyStr := msg.String()
			switch keyStr {
			case "enter":
				// Confirm edit
				return d.confirmEdit(m)
			case "esc":
				// Cancel edit
				return d.cancelEdit()
			case "backspace":
				// Allow backspace
				var cmd tea.Cmd
				d.textInput, cmd = d.textInput.Update(msg)
				return cmd
			default:
				// Only allow numeric input
				if len(keyStr) == 1 && keyStr[0] >= '0' && keyStr[0] <= '9' {
					var cmd tea.Cmd
					d.textInput, cmd = d.textInput.Update(msg)
					// Auto-submit after 2 digits
					if len(d.textInput.Value()) >= 2 {
						return tea.Batch(cmd, d.confirmEdit(m))
					}
					return cmd
				}
				// Ignore non-numeric input
				return nil
			}
		default:
			var cmd tea.Cmd
			d.textInput, cmd = d.textInput.Update(msg)
			return cmd
		}
	}

	// Handle key messages when not editing
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		// Check if it's a digit key to start editing
		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
			return d.startEdit(m, key)
		}
	}

	return nil
}

func (d *initiativeItemDelegate) startEdit(m *list.Model, initialDigit string) tea.Cmd {
	d.editingIndex = m.Index()
	d.textInput.SetValue(initialDigit)
	d.textInput.CursorEnd()
	return d.textInput.Focus()
}

func (d *initiativeItemDelegate) confirmEdit(m *list.Model) tea.Cmd {
	if d.editingIndex < 0 || d.editingIndex >= len(m.Items()) {
		d.editingIndex = -1
		d.textInput.Blur()
		return nil
	}

	// Parse the initiative value
	initValue, err := strconv.Atoi(d.textInput.Value())
	if err != nil || initValue <= 0 {
		// Invalid value, cancel edit
		d.editingIndex = -1
		d.textInput.Blur()
		return nil
	}

	// Update the item
	item, ok := m.Items()[d.editingIndex].(initiativeItem)
	if !ok {
		d.editingIndex = -1
		d.textInput.Blur()
		return nil
	}

	item.entry.Initiative = initValue
	cmd := m.SetItem(d.editingIndex, item)
	d.editingIndex = -1
	d.textInput.Blur()
	return cmd
}

func (d *initiativeItemDelegate) cancelEdit() tea.Cmd {
	d.editingIndex = -1
	d.textInput.Blur()
	return nil
}

// IsEditing returns true if currently editing an item
func (d *initiativeItemDelegate) IsEditing() bool {
	return d.editingIndex >= 0
}

func (d *initiativeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(initiativeItem)
	if !ok {
		return
	}

	// Selection indicator
	selection := "  "
	if m.Index() == index {
		selection = d.styles.selected.Render("> ")
	}

	// Initiative value or input
	// Format: [ 12 ] or [    ] or [input]
	var initDisplay string
	if d.editingIndex == index {
		// Show text input when editing - add leading space and right-pad to match non-editing format
		inputView := " " + d.textInput.View()
		for len(inputView) < 4 {
			inputView = inputView + " "
		}
		initDisplay = d.styles.inputStyle.Render("[" + inputView + "]")
	} else if item.entry.Initiative > 0 {
		initDisplay = d.styles.initiative.Render(fmt.Sprintf("[%3d ]", item.entry.Initiative))
	} else {
		initDisplay = d.styles.noInit.Render("[    ]")
	}

	// Quantity suffix for monsters with quantity > 1
	quantitySuffix := ""
	if item.entry.Quantity > 1 {
		quantitySuffix = d.styles.quantity.Render(fmt.Sprintf(" x%d", item.entry.Quantity))
	}

	// Name
	name := item.entry.Name + quantitySuffix

	// Calculate available width for name
	// Format: "> [ 12 ] Name x2"
	prefixWidth := lipgloss.Width(selection + initDisplay + " ")
	availableWidth := m.Width() - prefixWidth - 2 // 2 for padding

	// Ensure availableWidth is at least 1
	if availableWidth < 1 {
		availableWidth = 1
	}

	// Truncate name if needed
	nameWidth := lipgloss.Width(name)
	if nameWidth > availableWidth {
		if availableWidth > 3 {
			// Truncate name to available width minus 3 chars for "..."
			truncLen := availableWidth - 3
			if truncLen > len(name) {
				truncLen = len(name)
			}
			if truncLen < 0 {
				truncLen = 0
			}
			name = name[:truncLen] + "..."
		} else if availableWidth > 0 {
			// Very narrow - just truncate
			truncLen := availableWidth
			if truncLen > len(name) {
				truncLen = len(name)
			}
			name = name[:truncLen]
		} else {
			name = ""
		}
	}

	line := fmt.Sprintf("%s%s %s", selection, initDisplay, name)

	fmt.Fprint(w, line)
}

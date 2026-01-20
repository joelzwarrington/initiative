package initiative

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
)

// requestAddMonsterMsg is sent when the user presses 'a' to add a monster
type requestAddMonsterMsg struct{}

// requestEditMonsterMsg is sent when the user presses 'e' to edit a monster
type requestEditMonsterMsg struct {
	index int
	entry dnd.InitiativeEntry
}

// InitiativeListField is a custom huh.Field that displays a list of creatures
// with inline initiative editing
type InitiativeListField struct {
	key      string
	list     list.Model
	delegate *initiativeItemDelegate
	keys     initiativeListKeyMap
	focused  bool
	width    int
	height   int
	position huh.FieldPosition
	theme    *huh.Theme
	keyMap   *huh.KeyMap

	characters map[string]*dnd.Character
	sources    map[string]*dnd.Source

	styles initiativeListStyles
}

type initiativeListStyles struct {
	status lipgloss.Style
}

type initiativeListKeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Add         key.Binding
	Delete      key.Binding
	Edit        key.Binding
	EditEntry   key.Binding
	Roll        key.Binding
	Submit      key.Binding
	Cancel      key.Binding
	ConfirmEdit key.Binding
	CancelEdit  key.Binding
}

func defaultInitiativeListKeyMap() initiativeListKeyMap {
	return initiativeListKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓", "down"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "remove"),
		),
		Edit: key.NewBinding(
			key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
			key.WithHelp("0-9", "initiative"),
		),
		EditEntry: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Roll: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "roll"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "start"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		ConfirmEdit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		CancelEdit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

// NewInitiativeListField creates a new initiative list field
func NewInitiativeListField(characters map[string]*dnd.Character, sources map[string]*dnd.Source) *InitiativeListField {
	delegate := newInitiativeItemDelegate()
	keys := defaultInitiativeListKeyMap()

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetStatusBarItemName("creature", "creatures")
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	l.SetShowPagination(true)

	// Apply theme styles and remove padding from status bar
	l.Styles = currentTheme.list
	l.Styles.StatusBar = l.Styles.StatusBar.Padding(0)

	// Customize list keymaps
	l.KeyMap.CursorUp = keys.Up
	l.KeyMap.CursorDown = keys.Down

	f := &InitiativeListField{
		key:        "initiative_list",
		list:       l,
		delegate:   delegate,
		keys:       keys,
		characters: characters,
		sources:    sources,
		styles: initiativeListStyles{
			status: lipgloss.NewStyle().Foreground(currentTheme.colors.error),
		},
	}

	// Pre-populate with all characters
	f.addAllCharacters()
	f.updateKeyStates()

	return f
}

func (f *InitiativeListField) addAllCharacters() {
	// Collect characters for sorting
	type charEntry struct {
		uuid      string
		character *dnd.Character
	}
	chars := make([]charEntry, 0, len(f.characters))
	for uuid, character := range f.characters {
		chars = append(chars, charEntry{uuid, character})
	}

	// Sort by name for consistent ordering
	sort.Slice(chars, func(i, j int) bool {
		return chars[i].character.Name() < chars[j].character.Name()
	})

	items := make([]list.Item, 0, len(chars))
	for _, c := range chars {
		entry := dnd.InitiativeEntry{
			CreatureType: "character",
			CreatureID:   c.uuid,
			Name:         c.character.Name(),
			Initiative:   0,
			Quantity:     1,
			StatBlock:    nil,
		}
		items = append(items, initiativeItem{entry: entry})
	}
	f.list.SetItems(items)
}

// AddMonster adds a monster entry to the list and returns a status message command
func (f *InitiativeListField) AddMonster(entry dnd.InitiativeEntry) tea.Cmd {
	items := f.list.Items()
	items = append(items, initiativeItem{entry: entry})
	f.list.SetItems(items)
	// Move selection to the newly added item
	f.list.Select(len(items) - 1)
	f.updateKeyStates()

	// Return status message
	msg := fmt.Sprintf("Added %s", entry.Name)
	if entry.Quantity > 1 {
		msg = fmt.Sprintf("Added %d %s", entry.Quantity, entry.Name)
	}
	return f.list.NewStatusMessage(msg)
}

// UpdateMonster updates a monster entry at the given index
func (f *InitiativeListField) UpdateMonster(index int, entry dnd.InitiativeEntry) {
	items := f.list.Items()
	if index >= 0 && index < len(items) {
		items[index] = initiativeItem{entry: entry}
		f.list.SetItems(items)
	}
}

func (f *InitiativeListField) removeSelected() tea.Cmd {
	items := f.list.Items()
	idx := f.list.Index()
	if idx < 0 || idx >= len(items) {
		return nil
	}

	// Get the name before removing
	var name string
	if item, ok := items[idx].(initiativeItem); ok {
		name = item.entry.Name
	}

	// Remove the item
	newItems := make([]list.Item, 0, len(items)-1)
	newItems = append(newItems, items[:idx]...)
	newItems = append(newItems, items[idx+1:]...)
	f.list.SetItems(newItems)

	// Adjust selection if necessary
	if idx >= len(newItems) && len(newItems) > 0 {
		f.list.Select(len(newItems) - 1)
	}
	f.updateKeyStates()

	// Return status message
	if name != "" {
		return f.list.NewStatusMessage(fmt.Sprintf("Removed %s", name))
	}
	return nil
}

func (f *InitiativeListField) rollInitiative() tea.Cmd {
	dexMod, hasDex := f.getSelectedDexModifier()
	if !hasDex {
		return nil
	}

	items := f.list.Items()
	idx := f.list.Index()
	if idx < 0 || idx >= len(items) {
		return nil
	}

	item, ok := items[idx].(initiativeItem)
	if !ok {
		return nil
	}

	// Build dice notation: 1d20+modifier or 1d20-modifier
	var notation string
	if dexMod >= 0 {
		notation = fmt.Sprintf("1d20+%d", dexMod)
	} else {
		notation = fmt.Sprintf("1d20%d", dexMod) // negative already includes minus
	}

	roll, err := dnd.Roll(notation)
	if err != nil {
		return nil
	}

	// Update item with rolled initiative
	item.entry.Initiative = roll.Total
	f.list.SetItem(idx, item)

	// Move to next creature needing initiative
	f.delegate.selectNextNeedingInitiative(&f.list)
	f.updateKeyStates()

	// Return status message
	return f.list.NewStatusMessage(fmt.Sprintf("%s rolled %d (%s)", item.entry.Name, roll.Total, roll.Formula))
}

func (f *InitiativeListField) countNeedingInitiative() int {
	count := 0
	for _, item := range f.list.Items() {
		if initItem, ok := item.(initiativeItem); ok {
			if initItem.entry.Initiative <= 0 {
				count++
			}
		}
	}
	return count
}

// getSelectedDexModifier returns the DEX modifier for the selected creature
// Returns (modifier, true) if DEX is available, (0, false) otherwise
func (f *InitiativeListField) getSelectedDexModifier() (int, bool) {
	selectedItem := f.list.SelectedItem()
	if selectedItem == nil {
		return 0, false
	}

	item, ok := selectedItem.(initiativeItem)
	if !ok {
		return 0, false
	}

	// For monsters, get DEX from stat block
	if item.entry.StatBlock != nil {
		if item.entry.StatBlock.DEX.Score > 0 {
			return item.entry.StatBlock.DEX.Modifier, true
		}
	}

	// For characters, look up in characters map
	if item.entry.CreatureType == "character" {
		if char, exists := f.characters[item.entry.CreatureID]; exists {
			if char.DEX().Score > 0 {
				return char.DEX().Modifier, true
			}
		}
	}

	return 0, false
}

func (f *InitiativeListField) updateKeyStates() {
	hasItems := len(f.list.Items()) > 0
	isEditing := f.delegate.IsEditing()

	// Navigation and actions only available when not editing and has items
	f.keys.Up.SetEnabled(hasItems && !isEditing)
	f.keys.Down.SetEnabled(hasItems && !isEditing)
	f.keys.Delete.SetEnabled(hasItems && !isEditing)
	f.keys.Edit.SetEnabled(hasItems && !isEditing)
	f.keys.Add.SetEnabled(!isEditing)

	// Submit only available when not editing and all creatures have initiative
	allHaveInitiative := f.countNeedingInitiative() == 0
	f.keys.Submit.SetEnabled(!isEditing && allHaveInitiative)

	// Cancel only available when not editing
	f.keys.Cancel.SetEnabled(!isEditing)

	// Edit entry only available for monsters
	canEditEntry := false
	canRoll := false
	if hasItems && !isEditing {
		if selectedItem := f.list.SelectedItem(); selectedItem != nil {
			if item, ok := selectedItem.(initiativeItem); ok {
				canEditEntry = item.entry.CreatureType == "monster"
				// Can roll if creature has DEX
				_, hasDex := f.getSelectedDexModifier()
				canRoll = hasDex
			}
		}
	}
	f.keys.EditEntry.SetEnabled(canEditEntry)
	f.keys.Roll.SetEnabled(canRoll)

	// Confirm/cancel only available when editing
	f.keys.ConfirmEdit.SetEnabled(isEditing)
	f.keys.CancelEdit.SetEnabled(isEditing)
}

// Init implements huh.Field
func (f *InitiativeListField) Init() tea.Cmd {
	return nil
}

// Update implements huh.Field
func (f *InitiativeListField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Always forward non-key messages to list for internal state (status bar timeouts, etc.)
	if _, isKey := msg.(tea.KeyMsg); !isKey {
		var cmd tea.Cmd
		f.list, cmd = f.list.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// When editing, delegate key messages to the item delegate
	if f.delegate.IsEditing() {
		if keyMsg, isKey := msg.(tea.KeyMsg); isKey {
			cmd := f.delegate.Update(keyMsg, &f.list)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		f.updateKeyStates()
		return f, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.Add):
			// Request to add a monster
			return f, func() tea.Msg { return requestAddMonsterMsg{} }

		case key.Matches(msg, f.keys.EditEntry):
			// Request to edit selected monster
			if selectedItem := f.list.SelectedItem(); selectedItem != nil {
				if item, ok := selectedItem.(initiativeItem); ok {
					if item.entry.CreatureType == "monster" {
						return f, func() tea.Msg {
							return requestEditMonsterMsg{
								index: f.list.Index(),
								entry: item.entry,
							}
						}
					}
				}
			}
			return f, nil

		case key.Matches(msg, f.keys.Delete):
			// Remove selected item
			cmd := f.removeSelected()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return f, tea.Batch(cmds...)

		case key.Matches(msg, f.keys.Roll):
			// Roll initiative for selected creature
			cmd := f.rollInitiative()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return f, tea.Batch(cmds...)

		case key.Matches(msg, f.keys.Submit):
			// Submit form - return NextField command
			if f.keyMap != nil {
				return f, func() tea.Msg { return huh.NextField() }
			}
			return f, nil

		default:
			// Check if it's a digit key to start editing (only if there are items)
			if len(f.list.Items()) > 0 {
				keyStr := msg.String()
				if len(keyStr) == 1 && keyStr[0] >= '0' && keyStr[0] <= '9' {
					cmd := f.delegate.Update(msg, &f.list)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
					f.updateKeyStates()
					return f, tea.Batch(cmds...)
				}
			}
		}

		// Forward key messages to list for navigation
		var cmd tea.Cmd
		f.list, cmd = f.list.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	f.updateKeyStates()
	return f, tea.Batch(cmds...)
}

// View implements huh.Field
func (f *InitiativeListField) View() string {
	var sb strings.Builder

	// List view (no border)
	sb.WriteString(f.list.View())

	// Status line
	needInit := f.countNeedingInitiative()
	total := len(f.list.Items())
	if needInit > 0 {
		status := fmt.Sprintf("%d of %d need initiative", needInit, total)
		sb.WriteString("\n")
		sb.WriteString(f.styles.status.Render(status))
	}

	// Help text is rendered by the parent page via KeyBinds()

	return sb.String()
}

// Blur implements huh.Field
func (f *InitiativeListField) Blur() tea.Cmd {
	f.focused = false
	return nil
}

// Focus implements huh.Field
func (f *InitiativeListField) Focus() tea.Cmd {
	f.focused = true
	return nil
}

// Error implements huh.Field
func (f *InitiativeListField) Error() error {
	// Validate that all entries have initiative set
	for _, item := range f.list.Items() {
		if initItem, ok := item.(initiativeItem); ok {
			if initItem.entry.Initiative <= 0 {
				return fmt.Errorf("all creatures must have initiative set")
			}
		}
	}
	return nil
}

// Run implements huh.Field
func (f *InitiativeListField) Run() error {
	return nil
}

// RunAccessible implements huh.Field
func (f *InitiativeListField) RunAccessible(w io.Writer, r io.Reader) error {
	return nil
}

// Skip implements huh.Field
func (f *InitiativeListField) Skip() bool {
	return false
}

// Zoom implements huh.Field
func (f *InitiativeListField) Zoom() bool {
	return true // Take full height
}

// KeyBinds implements huh.Field
func (f *InitiativeListField) KeyBinds() []key.Binding {
	f.updateKeyStates()

	if f.delegate.IsEditing() {
		return []key.Binding{f.keys.ConfirmEdit, f.keys.CancelEdit}
	}

	var bindings []key.Binding
	if f.keys.Up.Enabled() {
		bindings = append(bindings, f.keys.Up, f.keys.Down)
	}
	if f.keys.Edit.Enabled() {
		bindings = append(bindings, f.keys.Edit)
	}
	if f.keys.EditEntry.Enabled() {
		bindings = append(bindings, f.keys.EditEntry)
	}
	if f.keys.Roll.Enabled() {
		bindings = append(bindings, f.keys.Roll)
	}
	if f.keys.Add.Enabled() {
		bindings = append(bindings, f.keys.Add)
	}
	if f.keys.Delete.Enabled() {
		bindings = append(bindings, f.keys.Delete)
	}
	if f.keys.Submit.Enabled() {
		bindings = append(bindings, f.keys.Submit)
	}
	if f.keys.Cancel.Enabled() {
		bindings = append(bindings, f.keys.Cancel)
	}

	return bindings
}

// WithTheme implements huh.Field
func (f *InitiativeListField) WithTheme(theme *huh.Theme) huh.Field {
	f.theme = theme
	return f
}

// WithAccessible implements huh.Field
func (f *InitiativeListField) WithAccessible(accessible bool) huh.Field {
	return f
}

// WithKeyMap implements huh.Field
func (f *InitiativeListField) WithKeyMap(keyMap *huh.KeyMap) huh.Field {
	f.keyMap = keyMap
	return f
}

// WithWidth implements huh.Field
func (f *InitiativeListField) WithWidth(width int) huh.Field {
	f.width = width
	f.list.SetWidth(width)
	return f
}

// WithHeight implements huh.Field
func (f *InitiativeListField) WithHeight(height int) huh.Field {
	f.height = height
	// Account for status and help lines
	listHeight := height - 4
	if listHeight < 3 {
		listHeight = 3
	}
	f.list.SetHeight(listHeight)
	return f
}

// WithPosition implements huh.Field
func (f *InitiativeListField) WithPosition(position huh.FieldPosition) huh.Field {
	f.position = position
	return f
}

// GetKey implements huh.Field
func (f *InitiativeListField) GetKey() string {
	return f.key
}

// GetValue implements huh.Field
func (f *InitiativeListField) GetValue() any {
	entries := make([]dnd.InitiativeEntry, 0, len(f.list.Items()))
	for _, item := range f.list.Items() {
		if initItem, ok := item.(initiativeItem); ok {
			entries = append(entries, initItem.entry)
		}
	}
	return entries
}

// Key sets the field key
func (f *InitiativeListField) Key(key string) *InitiativeListField {
	f.key = key
	return f
}

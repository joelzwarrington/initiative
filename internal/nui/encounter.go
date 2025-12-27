package nui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/termkit/skeleton"
)

var _ tea.Model = (*encounter)(nil)

type encounterView int

const (
	encounterPlaceholder encounterView = iota
	encounterForm
	encounterDetail
)

type encounter struct {
	Encounter

	skeleton *skeleton.Skeleton
	party    *map[string]Character

	view            encounterView
	form            *huh.Form
	list            list.Model
	help            help.Model
	placeholderKeys encounterPlaceholderKeyMap
	detailKeys      encounterDetailKeyMap
}

func newEncounter(skeleton *skeleton.Skeleton, party *map[string]Character) *encounter {
	// Create empty list for initiative groups
	initiativeList := list.New([]list.Item{}, &initiativeGroupItemDelegate{}, skeleton.GetContentWidth(), skeleton.GetContentHeight())
	initiativeList.SetStatusBarItemName("group", "groups")
	initiativeList.SetShowTitle(false)
	initiativeList.SetShowStatusBar(false)
	initiativeList.SetShowHelp(false)
	initiativeList.DisableQuitKeybindings()

	return &encounter{
		skeleton: skeleton,
		party:    party,

		view:            encounterPlaceholder,
		list:            initiativeList,
		help:            help.New(),
		placeholderKeys: newEncounterPlaceholderKeyMap(),
		detailKeys:      newEncounterDetailKeyMap(),
	}
}

func (e encounter) Init() tea.Cmd {
	if e.form != nil {
		return e.form.Init()
	}
	return nil
}

func (e encounter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch e.view {
		case encounterPlaceholder:
			if key.Matches(msg, e.placeholderKeys.startEncounter) {
				return e, tea.Cmd(func() tea.Msg {
					return startEncounterMsg{}
				})
			}
		case encounterDetail:
			if key.Matches(msg, e.detailKeys.back) {
				e.view = encounterPlaceholder
				e.IniativeGroups = []IniativeGroup{}
				e.Summary = ""
				e.StartedAt = time.Time{}
				e.EndedAt = time.Time{}
				return e, nil
			}
		}
	case startEncounterMsg:
		var characterOptions []huh.Option[string]

		if e.party != nil {
			for uuid, character := range *e.party {
				characterOptions = append(characterOptions,
					huh.NewOption(character.Name(), uuid).Selected(true),
				)
			}
		}

		e.form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("summary").
					Title("Summary"),
				huh.NewMultiSelect[string]().
					Key("characters").
					Title("Characters").
					Options(characterOptions...),
			),
		)
		e.view = encounterForm
		return e, e.form.Init()
	}

	switch e.view {
	case encounterForm:
		{
			form, cmd := e.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				e.form = f
			}

			if e.form.State == huh.StateCompleted {
				e.Summary = e.form.GetString("summary")
				e.StartedAt = time.Now()

				// Create initiative groups for each selected character
				selectedCharacterUUIDs := e.form.Get("characters").([]string)
				e.IniativeGroups = []IniativeGroup{}

				if e.party != nil {
					for _, uuid := range selectedCharacterUUIDs {
						if character, exists := (*e.party)[uuid]; exists {
							group := IniativeGroup{
								Iniative:  0, // Initiative will be set later
								Creatures: []Creature{character},
							}
							e.IniativeGroups = append(e.IniativeGroups, group)
						}
					}
				}

				// Update the list with initiative groups
				items := []list.Item{}
				for _, group := range e.IniativeGroups {
					items = append(items, initiativeGroupItem{group: group})
				}
				e.list.SetItems(items)

				e.form = nil
				e.view = encounterDetail
			}

			return e, cmd
		}
	case encounterDetail:
		{
			var cmd tea.Cmd
			e.list, cmd = e.list.Update(msg)
			return e, cmd
		}
	}
	return e, nil
}

func (e encounter) View() string {
	switch e.view {
	case encounterPlaceholder:
		{
			// Calculate available height for content
			helpStyle := lipgloss.NewStyle().PaddingBottom(1)
			helpView := helpStyle.Render(e.help.View(e.placeholderKeys))
			availHeight := e.skeleton.GetContentHeight() - lipgloss.Height(helpView)

			// Create main content area
			placeholderStyle := lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("240")).
				Align(lipgloss.Center)
			content := placeholderStyle.Render("No encounter started...")
			contentArea := lipgloss.NewStyle().
				Height(availHeight).
				Width(e.skeleton.GetContentWidth()).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				Render(content)

			return lipgloss.JoinVertical(lipgloss.Left, contentArea, helpView)
		}
	case encounterForm:
		{
			e.form.WithHeight(e.skeleton.GetContentHeight()).WithWidth(e.skeleton.GetContentWidth())
			return e.form.View()
		}
	case encounterDetail:
		{
			// Calculate available height for content
			helpStyle := lipgloss.NewStyle().PaddingBottom(1)
			helpView := helpStyle.Render(e.help.View(e.detailKeys))
			availHeight := e.skeleton.GetContentHeight() - lipgloss.Height(helpView)

			// Create header with encounter summary
			headerStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginBottom(1)
			header := headerStyle.Render(fmt.Sprintf("Encounter: %s", e.Summary))

			// Set list dimensions accounting for header and help
			headerHeight := lipgloss.Height(header) + 1 // +1 for margin
			listHeight := availHeight - headerHeight

			e.list.SetHeight(listHeight)
			e.list.SetWidth(e.skeleton.GetContentWidth())

			// Combine header and list
			content := lipgloss.JoinVertical(lipgloss.Left, header, e.list.View())

			return lipgloss.JoinVertical(lipgloss.Left, content, helpView)
		}
	}

	return ""
}

// Messages
type startEncounterMsg struct{}

// Key mappings
type encounterPlaceholderKeyMap struct {
	startEncounter key.Binding
}

func newEncounterPlaceholderKeyMap() encounterPlaceholderKeyMap {
	return encounterPlaceholderKeyMap{
		startEncounter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new encounter"),
		),
	}
}

func (k encounterPlaceholderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.startEncounter}
}

func (k encounterPlaceholderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.startEncounter},
	}
}

type encounterDetailKeyMap struct {
	back key.Binding
}

func newEncounterDetailKeyMap() encounterDetailKeyMap {
	return encounterDetailKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

func (k encounterDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.back}
}

func (k encounterDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.back},
	}
}

// List item for initiative groups
var _ list.Item = (*initiativeGroupItem)(nil)

type initiativeGroupItem struct {
	group IniativeGroup
}

func (i initiativeGroupItem) FilterValue() string {
	if len(i.group.Creatures) > 0 {
		return i.group.Creatures[0].Name()
	}
	return fmt.Sprintf("Initiative: %d", i.group.Iniative)
}

// List delegate for initiative groups
type initiativeGroupItemDelegate struct{}

func (d initiativeGroupItemDelegate) Height() int  { return 2 }
func (d initiativeGroupItemDelegate) Spacing() int { return 1 }
func (d *initiativeGroupItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d initiativeGroupItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(initiativeGroupItem)
	if !ok {
		return
	}

	// Initiative value styling
	initiativeText := "Initiative: TBD"
	if i.group.Iniative > 0 {
		initiativeText = fmt.Sprintf("Initiative: %d", i.group.Iniative)
	}

	initiativeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))

	// Creatures list
	creatureNames := []string{}
	for _, creature := range i.group.Creatures {
		creatureNames = append(creatureNames, creature.Name())
	}
	creaturesText := strings.Join(creatureNames, ", ")

	creatureStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// Combine text
	content := initiativeStyle.Render(initiativeText) + "\n" +
		creatureStyle.Render("  "+creaturesText)

	// Apply selection styling
	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170")).
				Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(content))
}

package initiative

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/joelzwarrington/initiative/dnd"
	"github.com/joelzwarrington/initiative/internal/components"
	"github.com/termkit/skeleton"
)

type characterPage struct {
	s      *skeleton.Skeleton
	width  int
	height int
	keys   CharacterPageKeyMap

	characters map[string]*dnd.Character

	// Current view state
	currentCharacter string // UUID of character being viewed

	// Sub-models
	emptyState    *components.EmptyState
	characterList *characterList
	characterForm *characterForm
	help          help.Model
}

func newCharacterPage(s *skeleton.Skeleton, characters map[string]*dnd.Character) *characterPage {
	keys := defaultCharacterPageKeyMap()

	page := &characterPage{
		s:      s,
		width:  s.GetContentWidth(),
		height: s.GetContentHeight(),
		keys:   keys,

		characters: characters,
	}

	page.emptyState = components.NewEmptyState(
		"No characters in party",
		page,
		s.GetContentWidth(), s.GetContentHeight(),
	)

	page.characterList = newCharacterList(characters, s.GetContentWidth(), s.GetContentHeight())
	page.help = help.New()

	return page
}

func (p *characterPage) Init() tea.Cmd {
	return nil
}

func (p *characterPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = p.s.GetContentWidth()
		p.height = p.s.GetContentHeight()
		if p.characterForm != nil {
			p.characterForm.SetSize(p.width, p.height)
		}
		if p.characterList != nil {
			p.characterList.SetSize(p.width, p.height)
		}
		if p.emptyState != nil {
			p.emptyState.SetSize(p.width, p.height)
		}
		return p, nil

	case tea.KeyMsg:
		// Handle global keys when not in form
		if !p.isAddingOrEditingCharacter() {
			switch {
			case key.Matches(msg, p.keys.NewCharacter):
				return p, tea.Sequence(append(cmds, p.beginNewCharacterForm())...)
			case key.Matches(msg, p.keys.Back) && p.isViewingCharacter():
				return p, tea.Sequence(append(cmds, p.backToList())...)
			}
		}

	case viewCharacterMsg:
		return p, tea.Sequence(append(cmds, p.viewCharacter(msg.uuid))...)

	case editCharacterMsg:
		return p, tea.Sequence(append(cmds, p.editCharacter(msg.uuid))...)

	case characterFormCancelledMsg:
		return p, tea.Sequence(append(cmds, p.cancelCharacterForm())...)

	case characterFormSubmittedMsg:
		return p, tea.Sequence(append(cmds, p.submitCharacterForm(msg))...)
	}

	// Update sub-models based on current state
	switch {
	case p.isAddingOrEditingCharacter():
		form, cmd := p.characterForm.Update(msg)
		cmds = append(cmds, cmd)
		if form, ok := form.(*characterForm); ok {
			p.characterForm = form
		}

	case p.isViewingList() || p.isViewingCharacter():
		list, cmd := p.characterList.Update(msg)
		cmds = append(cmds, cmd)
		if list, ok := list.(*characterList); ok {
			p.characterList = list
		}
	}

	return p, tea.Batch(cmds...)
}

func (p *characterPage) View() string {
	switch {
	case p.isAddingOrEditingCharacter():
		return p.characterForm.View()

	case p.isViewingCharacter():
		return p.renderCharacterDetail()

	case p.isViewingList():
		if p.hasCharacters() {
			return p.renderCharacterList()
		}
		return p.renderEmptyState()

	default:
		return "No view"
	}
}

func (p *characterPage) Key() string {
	return "party"
}

func (p *characterPage) Title() string {
	title := icons.PartyTab.Join("Party", nil)

	if p.isViewingCharacter() && p.currentCharacter != "" {
		if char, exists := p.characters[p.currentCharacter]; exists {
			name := char.Name()
			if len(name) > 18 {
				name = name[:15] + "..."
			}
			return fmt.Sprintf("%s > %s", title, name)
		}
	}
	return title
}

func (p *characterPage) ShortHelp() []key.Binding {
	switch {
	case p.isViewingCharacter():
		return []key.Binding{p.keys.Back}
	case p.isViewingList():
		keys := []key.Binding{p.keys.NewCharacter}
		if p.hasCharacters() && p.characterList != nil {
			// Add list navigation and item keys
			listKeys := p.characterList.ShortHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	default:
		return []key.Binding{}
	}
}

func (p *characterPage) FullHelp() [][]key.Binding {
	switch {
	case p.isViewingCharacter():
		return [][]key.Binding{{p.keys.Back}}
	case p.isViewingList():
		keys := [][]key.Binding{{p.keys.NewCharacter}}
		if p.hasCharacters() && p.characterList != nil {
			// Add list navigation and item keys
			listKeys := p.characterList.FullHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	default:
		return [][]key.Binding{}
	}
}

// State management helpers
func (p *characterPage) isAddingOrEditingCharacter() bool {
	return p.characterForm != nil
}

func (p *characterPage) isViewingCharacter() bool {
	return p.currentCharacter != ""
}

func (p *characterPage) isViewingList() bool {
	return !p.isAddingOrEditingCharacter() && !p.isViewingCharacter()
}

func (p *characterPage) hasCharacters() bool {
	return p.characters != nil && len(p.characters) > 0
}

// Navigation methods
func (p *characterPage) beginNewCharacterForm() tea.Cmd {
	p.s.LockTabs()
	p.characterForm = newCharacterForm("", nil, p.s.GetContentWidth(), p.s.GetContentHeight())
	return p.characterForm.Init()
}

func (p *characterPage) editCharacter(uuid string) tea.Cmd {
	var character *dnd.Character
	if p.characters != nil {
		if c, exists := p.characters[uuid]; exists {
			character = c
		}
	}

	p.s.LockTabs()
	p.characterForm = newCharacterForm(uuid, character, p.s.GetContentWidth(), p.s.GetContentHeight())
	return p.characterForm.Init()
}

func (p *characterPage) viewCharacter(uuid string) tea.Cmd {
	p.currentCharacter = uuid

	p.s.UpdatePageTitle(p.Key(), p.Title())
	return nil
}

func (p *characterPage) backToList() tea.Cmd {
	p.currentCharacter = ""
	p.s.UpdatePageTitle(p.Key(), p.Title())
	return nil
}

func (p *characterPage) cancelCharacterForm() tea.Cmd {
	p.characterForm = nil
	p.s.UnlockTabs()
	return nil
}

func (p *characterPage) submitCharacterForm(submission characterFormSubmittedMsg) tea.Cmd {
	p.characterForm = nil
	p.s.UnlockTabs()

	if submission.uuid != "" {
		// Editing existing character
		return tea.Cmd(func() tea.Msg {
			return characterUpdatedMsg{
				uuid:      submission.uuid,
				character: submission.character,
			}
		})
	} else {
		// Adding new character - generate new UUID
		uuid := uuid.New().String()
		return tea.Cmd(func() tea.Msg {
			return characterAddedMsg{
				uuid:      uuid,
				character: submission.character,
			}
		})
	}
}

// Render methods
func (p *characterPage) renderCharacterList() string {
	// Calculate available height for list (subtract help height)
	p.help.Width = p.s.GetContentWidth()
	helpStyle := lipgloss.NewStyle().Padding(0, 1)
	helpView := helpStyle.Render(p.help.View(p))
	helpHeight := lipgloss.Height(helpView)

	listHeight := p.s.GetContentHeight() - helpHeight
	p.characterList.SetSize(p.s.GetContentWidth(), listHeight)

	listView := p.characterList.View()
	return lipgloss.JoinVertical(lipgloss.Left, listView, helpView)
}

func (p *characterPage) renderEmptyState() string {
	// Let the empty state handle its own help display
	p.emptyState.SetSize(p.s.GetContentWidth(), p.s.GetContentHeight())
	return p.emptyState.View()
}

func (p *characterPage) renderCharacterDetail() string {
	var characterName string
	if p.characters != nil {
		if character, exists := p.characters[p.currentCharacter]; exists {
			characterName = character.Name()
		}
	}

	// Calculate available height for content
	p.help.Width = p.s.GetContentWidth()
	helpStyle := lipgloss.NewStyle().Padding(0, 1)
	helpView := helpStyle.Render(p.help.View(p))
	helpHeight := lipgloss.Height(helpView)
	availHeight := p.s.GetContentHeight() - helpHeight

	// Create main content area
	content := fmt.Sprintf("Viewing character: %s", characterName)
	contentArea := lipgloss.NewStyle().
		Height(availHeight).
		Width(p.s.GetContentWidth()).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top).
		Padding(1, 2, 0, 2).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, contentArea, helpView)
}

// Key mappings
type CharacterPageKeyMap struct {
	NewCharacter key.Binding
	Back         key.Binding
}

func defaultCharacterPageKeyMap() CharacterPageKeyMap {
	return CharacterPageKeyMap{
		NewCharacter: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new character"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

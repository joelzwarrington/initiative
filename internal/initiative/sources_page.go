package initiative

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joelzwarrington/initiative/dnd"
	"github.com/termkit/skeleton"
)

type sourcesPage struct {
	s      *skeleton.Skeleton
	width  int
	height int
	keys   SourcesPageKeyMap

	sources map[string]*dnd.Source

	// Current view state
	currentSource string // Key of source being viewed

	// Sub-models
	sourcesList *sourcesList
	help        help.Model
}

func newSourcesPage(s *skeleton.Skeleton, sources map[string]*dnd.Source) *sourcesPage {
	keys := defaultSourcesPageKeyMap()

	page := &sourcesPage{
		s:      s,
		width:  s.GetContentWidth(),
		height: s.GetContentHeight(),
		keys:   keys,

		sources: sources,
	}

	page.sourcesList = newSourcesList(sources, s.GetContentWidth(), s.GetContentHeight())
	page.help = help.New()
	page.help.Styles = currentTheme.help

	return page
}

func (p *sourcesPage) Init() tea.Cmd {
	return nil
}

func (p *sourcesPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = p.s.GetContentWidth()
		p.height = p.s.GetContentHeight()
		if p.sourcesList != nil {
			p.sourcesList.SetSize(p.width, p.height)
		}
		return p, nil

	case tea.KeyMsg:
		// Handle global keys when not viewing a source
		if !p.isViewingSource() {
			switch {
			case key.Matches(msg, p.keys.Back) && p.isViewingSource():
				return p, tea.Sequence(append(cmds, p.backToList())...)
			}
		} else {
			// When viewing a source, handle back key
			switch {
			case key.Matches(msg, p.keys.Back):
				return p, tea.Sequence(append(cmds, p.backToList())...)
			}
		}

	case viewSourceMsg:
		return p, tea.Sequence(append(cmds, p.viewSource(msg.key))...)
	}

	// Update sub-models based on current state
	if p.isViewingList() {
		list, cmd := p.sourcesList.Update(msg)
		cmds = append(cmds, cmd)
		if list, ok := list.(*sourcesList); ok {
			p.sourcesList = list
		}
	}

	return p, tea.Batch(cmds...)
}

func (p *sourcesPage) View() string {
	// Calculate help view for all cases
	p.help.Width = p.s.GetContentWidth()
	helpView := p.help.View(p)

	var content string
	switch {
	case p.isViewingSource():
		content = p.renderSourceDetailContent()

	case p.isViewingList():
		content = p.renderSourcesListContent()

	default:
		content = "No view"
	}

	return currentTheme.container.Render(lipgloss.JoinVertical(lipgloss.Left, content, helpView))
}

func (p *sourcesPage) Key() string {
	return "sources"
}

func (p *sourcesPage) Title() string {
	title := icons.SourcesTab.Join("Sources", nil)

	if p.isViewingSource() && p.currentSource != "" {
		if source, exists := p.sources[p.currentSource]; exists {
			name := source.Meta.Name
			if len(name) > 18 {
				name = name[:15] + "..."
			}
			return fmt.Sprintf("%s > %s", title, name)
		}
	}
	return title
}

func (p *sourcesPage) ShortHelp() []key.Binding {
	switch {
	case p.isViewingSource():
		return []key.Binding{p.keys.Back}
	case p.isViewingList():
		keys := []key.Binding{}
		if p.sourcesList != nil {
			// Add list navigation and item keys
			listKeys := p.sourcesList.ShortHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	default:
		return []key.Binding{}
	}
}

func (p *sourcesPage) FullHelp() [][]key.Binding {
	switch {
	case p.isViewingSource():
		return [][]key.Binding{{p.keys.Back}}
	case p.isViewingList():
		keys := [][]key.Binding{}
		if p.sourcesList != nil {
			// Add list navigation and item keys
			listKeys := p.sourcesList.FullHelp()
			keys = append(keys, listKeys...)
		}
		return keys
	default:
		return [][]key.Binding{}
	}
}

// State management helpers
func (p *sourcesPage) isViewingSource() bool {
	return p.currentSource != ""
}

func (p *sourcesPage) isViewingList() bool {
	return !p.isViewingSource()
}

// Navigation methods
func (p *sourcesPage) viewSource(key string) tea.Cmd {
	p.currentSource = key

	p.s.UpdatePageTitle(p.Key(), p.Title())
	return nil
}

func (p *sourcesPage) backToList() tea.Cmd {
	p.currentSource = ""
	p.s.UpdatePageTitle(p.Key(), p.Title())
	return nil
}

// Render methods
func (p *sourcesPage) renderSourcesListContent() string {
	// Calculate available height for list (subtract help height)
	p.help.Width = p.s.GetContentWidth()
	helpView := p.help.View(p)
	helpHeight := lipgloss.Height(helpView)

	listHeight := p.s.GetContentHeight() - helpHeight
	p.sourcesList.SetSize(p.s.GetContentWidth(), listHeight)

	return p.sourcesList.View()
}

func (p *sourcesPage) renderSourceDetailContent() string {
	var sourceName string
	if source, exists := p.sources[p.currentSource]; exists {
		sourceName = source.Meta.Name
	}

	// Calculate available height for content
	p.help.Width = p.s.GetContentWidth()
	helpView := p.help.View(p)
	helpHeight := lipgloss.Height(helpView)
	availHeight := p.s.GetContentHeight() - helpHeight

	// Create main content area
	content := fmt.Sprintf("Viewing source: %s", sourceName)
	return lipgloss.NewStyle().
		Height(availHeight).
		Width(p.s.GetContentWidth()).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top).
		Padding(1, 2, 0, 2).
		Render(content)
}

// Key mappings
type SourcesPageKeyMap struct {
	Back key.Binding
}

func defaultSourcesPageKeyMap() SourcesPageKeyMap {
	return SourcesPageKeyMap{
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
	}
}

// Messages
type viewSourceMsg struct {
	key string
}

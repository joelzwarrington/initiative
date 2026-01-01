package initiative

import (
	"strings"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/termkit/skeleton"
)

func TestSourcesPage_ViewingState(t *testing.T) {
	s := skeleton.NewSkeleton()
	sources := map[string]*dnd.Source{
		"test": {Meta: dnd.SourceMeta{Key: "test", Name: "Test Source", Version: "1.0"}},
	}
	page := newSourcesPage(s, sources)

	// Initially should be viewing list
	if !page.isViewingList() {
		t.Error("Expected to be viewing list initially")
	}
	if page.isViewingSource() {
		t.Error("Expected not to be viewing source initially")
	}

	// After viewing a source
	page.viewSource("test")
	if page.isViewingList() {
		t.Error("Expected not to be viewing list after viewing source")
	}
	if !page.isViewingSource() {
		t.Error("Expected to be viewing source after calling viewSource")
	}

	// After going back to list
	page.backToList()
	if !page.isViewingList() {
		t.Error("Expected to be viewing list after going back")
	}
	if page.isViewingSource() {
		t.Error("Expected not to be viewing source after going back")
	}
}

func TestSourcesPage_RenderSourcesList(t *testing.T) {
	tests := []struct {
		name     string
		sources  map[string]*dnd.Source
		expected []string // Contains patterns that should be in the output
	}{
		{
			name:    "empty sources list",
			sources: map[string]*dnd.Source{},
			expected: []string{
				"No sources", // The list should show some indication of no sources
			},
		},
		{
			name: "single source",
			sources: map[string]*dnd.Source{
				"srd": {Meta: dnd.SourceMeta{Key: "srd", Name: "System Reference Document", Version: "5.1"}},
			},
			expected: []string{
				"1. System Reference Document", // Should show the source name
				"1 source",                     // Status bar should show count
			},
		},
		{
			name: "multiple sources sorted",
			sources: map[string]*dnd.Source{
				"xgte": {Meta: dnd.SourceMeta{Key: "xgte", Name: "Xanathar's Guide to Everything", Version: "1.0"}},
				"phb":  {Meta: dnd.SourceMeta{Key: "phb", Name: "Player's Handbook", Version: "1.0"}},
				"mm":   {Meta: dnd.SourceMeta{Key: "mm", Name: "Monster Manual", Version: "1.0"}},
			},
			expected: []string{
				"1. Monster Manual",    // Should be sorted alphabetically
				"2. Player's Handbook", // by name
				"3. Xanathar's Guide to Everything",
				"3 sources", // Status bar should show count
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skeleton.NewSkeleton()
			page := newSourcesPage(s, tt.sources)

			// Set consistent dimensions directly on the page
			page.width = 80
			page.height = 24
			if page.sourcesList != nil {
				page.sourcesList.SetSize(80, 24)
			}

			// Ensure we're viewing the list
			page.currentSource = ""

			output := page.View()

			for _, pattern := range tt.expected {
				if !strings.Contains(output, pattern) {
					t.Errorf("Expected output to contain %q, but it didn't.\nFull output:\n%s", pattern, output)
				}
			}
		})
	}
}

func TestSourcesPage_RenderSourceDetail(t *testing.T) {
	tests := []struct {
		name           string
		sources        map[string]*dnd.Source
		selectedSource string
		expected       []string
	}{
		{
			name: "viewing source detail",
			sources: map[string]*dnd.Source{
				"srd": {Meta: dnd.SourceMeta{Key: "srd", Name: "System Reference Document", Version: "5.1"}},
			},
			selectedSource: "srd",
			expected: []string{
				"Viewing source: System Reference Document",
			},
		},
		{
			name: "viewing non-existent source",
			sources: map[string]*dnd.Source{
				"srd": {Meta: dnd.SourceMeta{Key: "srd", Name: "System Reference Document", Version: "5.1"}},
			},
			selectedSource: "nonexistent",
			expected: []string{
				"Viewing source:", // Should still show viewing template even with empty name
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skeleton.NewSkeleton()
			page := newSourcesPage(s, tt.sources)

			// Set consistent dimensions directly on the page
			page.width = 80
			page.height = 24

			// Set the current source to trigger detail view
			page.currentSource = tt.selectedSource

			output := page.View()

			for _, pattern := range tt.expected {
				if !strings.Contains(output, pattern) {
					t.Errorf("Expected output to contain %q, but it didn't.\nFull output:\n%s", pattern, output)
				}
			}
		})
	}
}

func TestSourcesPage_ExactOutput(t *testing.T) {
	tests := []struct {
		name     string
		sources  map[string]*dnd.Source
		viewMode string   // "list" or "detail"
		detail   string   // source key to view in detail mode
		contains []string // exact patterns that must be in output
	}{
		{
			name: "exact list output with multiple sources",
			sources: map[string]*dnd.Source{
				"mm":  {Meta: dnd.SourceMeta{Key: "mm", Name: "Monster Manual", Version: "1.0"}},
				"phb": {Meta: dnd.SourceMeta{Key: "phb", Name: "Player's Handbook", Version: "1.0"}},
			},
			viewMode: "list",
			contains: []string{
				"1. Monster Manual",
				"2. Player's Handbook",
				"filter", // Should have filter functionality
				"↑", "↓", // Should have navigation keys
				"enter", // Should have enter to view
			},
		},
		{
			name: "exact detail output",
			sources: map[string]*dnd.Source{
				"srd": {Meta: dnd.SourceMeta{Key: "srd", Name: "System Reference Document", Version: "5.1"}},
			},
			viewMode: "detail",
			detail:   "srd",
			contains: []string{
				"Viewing source: System Reference Document",
				"esc", // Should have escape key to go back
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skeleton.NewSkeleton()
			page := newSourcesPage(s, tt.sources)

			// Set consistent dimensions
			page.width = 80
			page.height = 24
			if page.sourcesList != nil {
				page.sourcesList.SetSize(80, 24)
			}

			// Set view mode
			if tt.viewMode == "detail" {
				page.currentSource = tt.detail
			} else {
				page.currentSource = ""
			}

			output := page.View()

			for _, pattern := range tt.contains {
				if !strings.Contains(output, pattern) {
					t.Errorf("Expected output to contain %q, but it didn't.\nFull output:\n%q", pattern, output)
				}
			}
		})
	}
}

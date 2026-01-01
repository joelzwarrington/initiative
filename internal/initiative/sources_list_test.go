package initiative

import (
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
)

func TestSourcesList_FilterValue(t *testing.T) {
	source := sourceItem{
		Source: dnd.Source{
			Meta: dnd.SourceMeta{
				Key:  "srd",
				Name: "System Reference Document",
			},
		},
	}

	expected := "srd System Reference Document"
	if source.FilterValue() != expected {
		t.Errorf("Expected FilterValue %q, got %q", expected, source.FilterValue())
	}
}

func TestSourcesList_ToItems(t *testing.T) {
	sources := map[string]*dnd.Source{
		"zebra": {Meta: dnd.SourceMeta{Key: "zebra", Name: "Zebra Source"}},
		"alpha": {Meta: dnd.SourceMeta{Key: "alpha", Name: "Alpha Source"}},
		"beta":  {Meta: dnd.SourceMeta{Key: "beta", Name: "Beta Source"}},
	}

	sl := newSourcesList(sources, 80, 24)
	items := sl.toItems()

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Check that items are sorted by name
	expectedOrder := []string{"Alpha Source", "Beta Source", "Zebra Source"}
	for i, item := range items {
		sourceItem := item.(sourceItem)
		if sourceItem.Meta.Name != expectedOrder[i] {
			t.Errorf("Expected item %d to be %q, got %q", i, expectedOrder[i], sourceItem.Meta.Name)
		}
	}
}

func TestSourcesList_EmptyList(t *testing.T) {
	sl := newSourcesList(map[string]*dnd.Source{}, 80, 24)
	items := sl.toItems()

	if len(items) != 0 {
		t.Errorf("Expected 0 items for empty sources list, got %d", len(items))
	}
}

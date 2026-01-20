package initiative

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joelzwarrington/initiative/dnd"
)

func TestSelectNextNeedingInitiative(t *testing.T) {
	tests := []struct {
		name          string
		items         []dnd.InitiativeEntry
		currentIndex  int
		expectedIndex int
	}{
		{
			name: "selects next item needing initiative",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 10},
				{Name: "B", Initiative: 0},
				{Name: "C", Initiative: 0},
			},
			currentIndex:  0,
			expectedIndex: 1,
		},
		{
			name: "skips items that have initiative",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 10},
				{Name: "B", Initiative: 15},
				{Name: "C", Initiative: 0},
			},
			currentIndex:  0,
			expectedIndex: 2,
		},
		{
			name: "wraps around to beginning",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 0},
				{Name: "B", Initiative: 10},
				{Name: "C", Initiative: 10},
			},
			currentIndex:  1,
			expectedIndex: 0,
		},
		{
			name: "stays on current when all have initiative",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 10},
				{Name: "B", Initiative: 15},
				{Name: "C", Initiative: 20},
			},
			currentIndex:  1,
			expectedIndex: 1,
		},
		{
			name: "handles single item needing initiative",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 10},
				{Name: "B", Initiative: 0},
				{Name: "C", Initiative: 10},
			},
			currentIndex:  0,
			expectedIndex: 1,
		},
		{
			name: "handles last item in list",
			items: []dnd.InitiativeEntry{
				{Name: "A", Initiative: 0},
				{Name: "B", Initiative: 10},
				{Name: "C", Initiative: 10},
			},
			currentIndex:  2,
			expectedIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := newInitiativeItemDelegate()

			listItems := make([]list.Item, len(tt.items))
			for i, entry := range tt.items {
				listItems[i] = initiativeItem{entry: entry}
			}

			l := list.New(listItems, delegate, 80, 24)
			l.Select(tt.currentIndex)

			delegate.selectNextNeedingInitiative(&l)

			if l.Index() != tt.expectedIndex {
				t.Errorf("expected index %d, got %d", tt.expectedIndex, l.Index())
			}
		})
	}
}

func TestSelectNextNeedingInitiativeEmptyList(t *testing.T) {
	delegate := newInitiativeItemDelegate()
	l := list.New([]list.Item{}, delegate, 80, 24)

	// Should not panic on empty list
	delegate.selectNextNeedingInitiative(&l)

	if l.Index() != 0 {
		t.Errorf("expected index 0 for empty list, got %d", l.Index())
	}
}

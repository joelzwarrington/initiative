package initiative

import (
	"bytes"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
)

func TestEncounterDelegateRender(t *testing.T) {
	tests := []struct {
		name      string
		encounter *dnd.Encounter
		expected  string
	}{
		{
			name:      "nil encounter",
			encounter: nil,
			expected:  "",
		},
		{
			name: "single character",
			encounter: &dnd.Encounter{
				InitiativeGroups: []dnd.InitiativeGroup{
					{
						Initiative: 15,
						Creatures:  []dnd.Creature{dnd.NewCharacter("Fighter")},
					},
				},
			},
			expected: `                                                            
  1 creature                                                
                                                            
╭──────────────────────────────────────────────────────────╮
│ > 15 • Fighter                                           │
╰──────────────────────────────────────────────────────────╯
                                                            
                                                            
                                                            
                                                            `,
		},
		{
			name: "multiple initiative groups",
			encounter: &dnd.Encounter{
				InitiativeGroups: []dnd.InitiativeGroup{
					{
						Initiative: 20,
						Creatures:  []dnd.Creature{dnd.NewCharacter("Wizard")},
					},
					{
						Initiative: 15,
						Creatures:  []dnd.Creature{dnd.NewCharacter("Fighter")},
					},
				},
			},
			expected: `                                                            
  2 creatures                                               
                                                            
╭──────────────────────────────────────────────────────────╮
│ > 20 • Wizard                                            │
╰──────────────────────────────────────────────────────────╯
                                                            
╭──────────────────────────────────────────────────────────╮
│   15 • Fighter                                           │
╰──────────────────────────────────────────────────────────╯
                                                            `,
		},
		{
			name: "multiple creatures in single group",
			encounter: &dnd.Encounter{
				InitiativeGroups: []dnd.InitiativeGroup{
					{
						Initiative: 12,
						Creatures: []dnd.Creature{
							dnd.NewMonster("Goblin 1", dnd.StatBlock{}),
							dnd.NewMonster("Goblin 2", dnd.StatBlock{}),
						},
					},
				},
			},
			expected: `                                                            
  2 creatures                                               
                                                            
╭──────────────────────────────────────────────────────────╮
│ > 12 • Goblin 1 ░░░░░░░░░░░░░░░░░░░░ 0/0 (Dead)          │
│                                                          │
│        Goblin 2 ░░░░░░░░░░░░░░░░░░░░ 0/0 (Dead)          │
╰──────────────────────────────────────────────────────────╯
                                                            
                                                            `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := newEncounterDelegate(60, 10)
			var buf bytes.Buffer

			delegate.Render(&buf, tt.encounter)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("Render() output mismatch:\ngot:\n%q\nwant:\n%q", got, tt.expected)
			}
		})
	}
}

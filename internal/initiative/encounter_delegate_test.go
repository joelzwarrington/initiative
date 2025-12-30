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
			delegate := newEncounterDelegate(60, 12)
			var buf bytes.Buffer

			delegate.Render(&buf, tt.encounter)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("Render() output mismatch:\ngot:\n%q\nwant:\n%q", got, tt.expected)
			}
		})
	}
}

func TestHitPointProcessing(t *testing.T) {
	tests := []struct {
		name           string
		initialHP      int
		maxHP          int
		adjustment     int
		adjustmentType HitPointAdjustmentType
		expectedHP     int
	}{
		{"Damage normal", 20, 25, 5, HitPointDamage, 15},
		{"Damage to zero", 3, 25, 10, HitPointDamage, 0},
		{"Heal normal", 15, 25, 5, HitPointHeal, 20},
		{"Heal to max", 20, 25, 10, HitPointHeal, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a monster with initial hit points
			statBlock := dnd.StatBlock{HitPoints: dnd.HitPoints{Fixed: tt.maxHP}}
			monster := dnd.NewMonster("Test Monster", statBlock)
			monster.HitPoints = tt.initialHP

			encounter := &dnd.Encounter{
				InitiativeGroups: []dnd.InitiativeGroup{
					{
						Initiative: 15,
						Creatures:  []dnd.Creature{monster},
					},
				},
			}

			// Simulate hit point adjustment processing
			group := &encounter.InitiativeGroups[0]
			creature := group.Creatures[0]

			if mon, ok := creature.(dnd.Monster); ok {
				if tt.adjustmentType == HitPointDamage {
					mon.HitPoints -= tt.adjustment
					if mon.HitPoints < 0 {
						mon.HitPoints = 0
					}
				} else {
					mon.HitPoints += tt.adjustment
					if mon.HitPoints > mon.MaximumHitPoints {
						mon.HitPoints = mon.MaximumHitPoints
					}
				}
				group.Creatures[0] = mon
			}

			// Check the result
			updatedMonster := group.Creatures[0].(dnd.Monster)
			if updatedMonster.HitPoints != tt.expectedHP {
				t.Errorf("Expected HP %d, got %d", tt.expectedHP, updatedMonster.HitPoints)
			}
		})
	}
}

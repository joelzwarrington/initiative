package initiative

import (
	"bytes"
	"testing"
	"time"

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
			encounter: func() *dnd.Encounter {
				char := dnd.NewCharacter("Fighter")
				var creature dnd.Creature = char
				return dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(15, []*dnd.Creature{&creature}),
				})
			}(),
			expected: `               
  1 creature   
               
 > 15 • Fighter
               
               
               
               
               
               
               
               `,
		},
		{
			name: "multiple initiative groups",
			encounter: func() *dnd.Encounter {
				wizard := dnd.NewCharacter("Wizard")
				fighter := dnd.NewCharacter("Fighter")
				var wizardCreature dnd.Creature = wizard
				var fighterCreature dnd.Creature = fighter
				e := dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(20, []*dnd.Creature{&wizardCreature}),
					dnd.NewInitiativeGroup(15, []*dnd.Creature{&fighterCreature}),
				})
				return e
			}(),
			expected: `               
  2 creatures  
               
 > 20 • Wizard 
               
               
   15 • Fighter
               
               
               
               
               `,
		},
		{
			name: "multiple creatures in single group",
			encounter: func() *dnd.Encounter {
				goblin1 := dnd.NewMonster("Goblin 1", dnd.StatBlock{HitPoints: dnd.HitPoints{Fixed: 7}})
				goblin2 := dnd.NewMonster("Goblin 2", dnd.StatBlock{HitPoints: dnd.HitPoints{Fixed: 7}})
				var creature1 dnd.Creature = goblin1
				var creature2 dnd.Creature = goblin2
				return dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(12, []*dnd.Creature{&creature1, &creature2}),
				})
			}(),
			expected: `                                          
  2 creatures                             
                                          
 > 12 • Goblin 1                          
        ████████████████████ 7/7 (Healthy)
                                          
   12 • Goblin 2                          
        ████████████████████ 7/7 (Healthy)
                                          
                                          
                                          
                                          `,
		},
		{
			name: "character with health bar",
			encounter: func() *dnd.Encounter {
				paladin := dnd.NewCharacter("Paladin").WithHealth(25, 25)
				paladin.AdjustHP(-5)
				var creature dnd.Creature = paladin
				return dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(18, []*dnd.Creature{&creature}),
				})
			}(),
			expected: `                                         
  1 creature                             
                                         
 > 18 • Paladin                          
        ████████████████░░░░ 20/25 (Hurt)
                                         
                                         
                                         
                                         
                                         
                                         
                                         `,
		},
		{
			name: "monster with armor class",
			encounter: func() *dnd.Encounter {
				orc := dnd.NewMonster("Orc", dnd.StatBlock{
					ArmorClass: dnd.ArmorClass{Value: 13},
					HitPoints:  dnd.HitPoints{Fixed: 15},
				})
				var creature dnd.Creature = orc
				return dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(14, []*dnd.Creature{&creature}),
				})
			}(),
			expected: `                                            
  1 creature                                
                                            
 > 14 • Orc • 󰒙 13                          
        ████████████████████ 15/15 (Healthy)
                                            
                                            
                                            
                                            
                                            
                                            
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
			monster.SetHP(tt.initialHP)

			var monsterPtr dnd.Creature = monster
			encounter := dnd.NewEncounter("", time.Time{}, []*dnd.InitiativeGroup{
				dnd.NewInitiativeGroup(15, []*dnd.Creature{&monsterPtr}),
			})

			// Simulate hit point adjustment processing using the new method
			groups := encounter.InitiativeGroups()
			creature := groups[0].Creatures()[0]

			// Convert damage to negative, healing to positive (like in encounter_page.go)
			var adjustment int
			if tt.adjustmentType == HitPointDamage {
				adjustment = -tt.adjustment
			} else {
				adjustment = tt.adjustment
			}

			// Apply adjustment using the new method
			(*creature).AdjustHP(adjustment)

			// Check the result
			if got := (*creature).HP(); got != tt.expectedHP {
				t.Errorf("Expected HP %d, got %d", tt.expectedHP, got)
			}
		})
	}
}

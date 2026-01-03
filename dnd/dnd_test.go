package dnd

import "testing"

func TestCreatureInterface(t *testing.T) {
	var creatures []Creature

	char := NewCharacter("Aragorn")
	monster := NewMonster("Orc", StatBlock{HitPoints: HitPoints{Fixed: 15}})

	creatures = append(creatures, char)
	creatures = append(creatures, monster)

	expectedNames := []string{"Aragorn", "Orc"}

	for i, creature := range creatures {
		if creature.Name() != expectedNames[i] {
			t.Errorf("Creature[%d].Name() = %v, want %v", i, creature.Name(), expectedNames[i])
		}
	}
}

func TestCreature_ArmorClass(t *testing.T) {
	tests := []struct {
		name     string
		creature Creature
		wantAC   int
	}{
		{
			name:     "character has no AC (returns 0)",
			creature: NewCharacter("Fighter"),
			wantAC:   0,
		},
		{
			name:     "character with health still has no AC",
			creature: NewCharacter("Rogue").WithHealth(15, 25),
			wantAC:   0,
		},
		{
			name: "monster with AC",
			creature: NewMonster("Orc", StatBlock{
				ArmorClass: ArmorClass{Value: 13},
				HitPoints:  HitPoints{Fixed: 15},
			}),
			wantAC: 13,
		},
		{
			name: "monster with high AC",
			creature: NewMonster("Dragon", StatBlock{
				ArmorClass: ArmorClass{Value: 18},
				HitPoints:  HitPoints{Fixed: 200},
			}),
			wantAC: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.creature.AC(); got != tt.wantAC {
				t.Errorf("ArmorClass() = %v, want %v", got, tt.wantAC)
			}
		})
	}
}

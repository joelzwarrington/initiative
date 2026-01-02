package dnd

import (
	"testing"
)

func TestEncounter_AdvanceTurn(t *testing.T) {
	tests := []struct {
		name      string
		encounter Encounter
		wantRound int
		wantTurn  int
	}{
		{
			name: "advance turn in middle of round",
			encounter: Encounter{
				Round:     1,
				turnIndex: 0,
				InitiativeGroups: []InitiativeGroup{
					{Initiative: 15},
					{Initiative: 10},
					{Initiative: 5},
				},
			},
			wantRound: 1,
			wantTurn:  1,
		},
		{
			name: "advance to next round",
			encounter: Encounter{
				Round:     1,
				turnIndex: 2,
				InitiativeGroups: []InitiativeGroup{
					{Initiative: 15},
					{Initiative: 10},
					{Initiative: 5},
				},
			},
			wantRound: 2,
			wantTurn:  0,
		},
		{
			name: "empty initiative groups",
			encounter: Encounter{
				Round:            1,
				turnIndex:        0,
				InitiativeGroups: []InitiativeGroup{},
			},
			wantRound: 1,
			wantTurn:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.encounter
			e.AdvanceTurn()
			if e.Round != tt.wantRound {
				t.Errorf("AdvanceTurn() round = %v, want %v", e.Round, tt.wantRound)
			}
			if e.GetTurnIndex() != tt.wantTurn {
				t.Errorf("AdvanceTurn() turnIndex = %v, want %v", e.GetTurnIndex(), tt.wantTurn)
			}
		})
	}
}

func TestEncounter_PreviousTurn(t *testing.T) {
	tests := []struct {
		name      string
		encounter Encounter
		wantRound int
		wantTurn  int
	}{
		{
			name: "previous turn in middle of round",
			encounter: Encounter{
				Round:     1,
				turnIndex: 1,
				InitiativeGroups: []InitiativeGroup{
					{Initiative: 15},
					{Initiative: 10},
					{Initiative: 5},
				},
			},
			wantRound: 1,
			wantTurn:  0,
		},
		{
			name: "previous turn to previous round",
			encounter: Encounter{
				Round:     2,
				turnIndex: 0,
				InitiativeGroups: []InitiativeGroup{
					{Initiative: 15},
					{Initiative: 10},
					{Initiative: 5},
				},
			},
			wantRound: 1,
			wantTurn:  2,
		},
		{
			name: "cannot go before round 1, turn 0",
			encounter: Encounter{
				Round:     1,
				turnIndex: 0,
				InitiativeGroups: []InitiativeGroup{
					{Initiative: 15},
					{Initiative: 10},
					{Initiative: 5},
				},
			},
			wantRound: 1,
			wantTurn:  0,
		},
		{
			name: "empty initiative groups",
			encounter: Encounter{
				Round:            1,
				turnIndex:        0,
				InitiativeGroups: []InitiativeGroup{},
			},
			wantRound: 1,
			wantTurn:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.encounter
			e.PreviousTurn()
			if e.Round != tt.wantRound {
				t.Errorf("PreviousTurn() round = %v, want %v", e.Round, tt.wantRound)
			}
			if e.GetTurnIndex() != tt.wantTurn {
				t.Errorf("PreviousTurn() turnIndex = %v, want %v", e.GetTurnIndex(), tt.wantTurn)
			}
		})
	}
}

func TestEncounter_GetTurnIndex(t *testing.T) {
	e := Encounter{turnIndex: 3}
	if got := e.GetTurnIndex(); got != 3 {
		t.Errorf("GetTurnIndex() = %v, want %v", got, 3)
	}
}

func TestNewCharacter(t *testing.T) {
	name := "Frodo"
	char := NewCharacter(name)

	if char.Name != name {
		t.Errorf("NewCharacter() Name = %v, want %v", char.Name, name)
	}
	if char.GetName() != name {
		t.Errorf("NewCharacter().GetName() = %v, want %v", char.GetName(), name)
	}
}

func TestCharacter_WithName(t *testing.T) {
	char := Character{Name: "Frodo"}
	newChar := char.WithName("Gandalf")

	if newChar.Name != "Gandalf" {
		t.Errorf("WithName() Name = %v, want %v", newChar.Name, "Gandalf")
	}
	if char.Name != "Frodo" {
		t.Errorf("WithName() should not modify original, got %v", char.Name)
	}
}

func TestNewMonster(t *testing.T) {
	statBlock := StatBlock{
		HitPoints: HitPoints{Fixed: 25},
		Challenge: Challenge{Rating: 1, ExperiencePoints: 200},
	}

	monster := NewMonster("Goblin", statBlock)

	if monster.Name != "Goblin" {
		t.Errorf("NewMonster() Name = %v, want %v", monster.Name, "Goblin")
	}
	if monster.HitPoints != 25 {
		t.Errorf("NewMonster() HitPoints = %v, want %v", monster.HitPoints, 25)
	}
	if monster.MaximumHitPoints != 25 {
		t.Errorf("NewMonster() MaximumHitPoints = %v, want %v", monster.MaximumHitPoints, 25)
	}
	if monster.GetName() != "Goblin" {
		t.Errorf("NewMonster().GetName() = %v, want %v", monster.GetName(), "Goblin")
	}
}

func TestCreatureInterface(t *testing.T) {
	var creatures []Creature

	char := NewCharacter("Aragorn")
	monster := NewMonster("Orc", StatBlock{HitPoints: HitPoints{Fixed: 15}})

	creatures = append(creatures, &char)
	creatures = append(creatures, &monster)

	expectedNames := []string{"Aragorn", "Orc"}

	for i, creature := range creatures {
		if creature.GetName() != expectedNames[i] {
			t.Errorf("Creature[%d].GetName() = %v, want %v", i, creature.GetName(), expectedNames[i])
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
			creature: func() Creature { c := NewCharacter("Fighter"); return &c }(),
			wantAC:   0,
		},
		{
			name:     "character with health still has no AC",
			creature: func() Creature { c := NewCharacterWithHealth("Rogue", 25); return &c }(),
			wantAC:   0,
		},
		{
			name: "monster with AC",
			creature: func() Creature {
				m := NewMonster("Orc", StatBlock{
					ArmorClass: ArmorClass{Value: 13},
					HitPoints:  HitPoints{Fixed: 15},
				})
				return &m
			}(),
			wantAC: 13,
		},
		{
			name: "monster with high AC",
			creature: func() Creature {
				m := NewMonster("Dragon", StatBlock{
					ArmorClass: ArmorClass{Value: 18},
					HitPoints:  HitPoints{Fixed: 200},
				})
				return &m
			}(),
			wantAC: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.creature.GetArmorClass(); got != tt.wantAC {
				t.Errorf("GetArmorClass() = %v, want %v", got, tt.wantAC)
			}
		})
	}
}

func TestCharacter_Health(t *testing.T) {
	tests := []struct {
		name    string
		char    Character
		wantHP  int
		wantMax int
	}{
		{
			name:    "new character with no health",
			char:    NewCharacter("Test"),
			wantHP:  0,
			wantMax: 0,
		},
		{
			name:    "new character with health",
			char:    NewCharacterWithHealth("Test", 20),
			wantHP:  20,
			wantMax: 20,
		},
		{
			name:    "character with custom health",
			char:    NewCharacter("Test").WithHealth(15, 20),
			wantHP:  15,
			wantMax: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.char.GetHitPoints(); got != tt.wantHP {
				t.Errorf("Character.GetHitPoints() = %v, want %v", got, tt.wantHP)
			}
			if got := tt.char.GetMaximumHitPoints(); got != tt.wantMax {
				t.Errorf("Character.GetMaximumHitPoints() = %v, want %v", got, tt.wantMax)
			}
		})
	}
}

func TestCharacter_SetHitPoints(t *testing.T) {
	tests := []struct {
		name     string
		char     Character
		newHP    int
		wantHP   int
		wantName string
	}{
		{
			name:     "set normal hit points",
			char:     NewCharacterWithHealth("Test", 20),
			newHP:    15,
			wantHP:   15,
			wantName: "Test",
		},
		{
			name:     "set hit points to zero",
			char:     NewCharacterWithHealth("Test", 20),
			newHP:    0,
			wantHP:   0,
			wantName: "Test",
		},
		{
			name:     "set negative hit points (should clamp to 0)",
			char:     NewCharacterWithHealth("Test", 20),
			newHP:    -5,
			wantHP:   0,
			wantName: "Test",
		},
		{
			name:     "set hit points above maximum (should clamp)",
			char:     NewCharacterWithHealth("Test", 20),
			newHP:    25,
			wantHP:   20,
			wantName: "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := tt.char
			char.SetHitPoints(tt.newHP)

			if got := char.GetHitPoints(); got != tt.wantHP {
				t.Errorf("Character.GetHitPoints() after SetHitPoints(%v) = %v, want %v", tt.newHP, got, tt.wantHP)
			}
			if got := char.GetName(); got != tt.wantName {
				t.Errorf("Character.GetName() after SetHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMonster_Health(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	monster := NewMonster("Orc", statBlock)

	if got := monster.GetHitPoints(); got != 15 {
		t.Errorf("Monster.GetHitPoints() = %v, want %v", got, 15)
	}
	if got := monster.GetMaximumHitPoints(); got != 15 {
		t.Errorf("Monster.GetMaximumHitPoints() = %v, want %v", got, 15)
	}
}

func TestMonster_SetHitPoints(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	baseMonster := NewMonster("Orc", statBlock)

	tests := []struct {
		name     string
		monster  Monster
		newHP    int
		wantHP   int
		wantName string
	}{
		{
			name:     "set normal hit points",
			monster:  baseMonster,
			newHP:    10,
			wantHP:   10,
			wantName: "Orc",
		},
		{
			name:     "set hit points to zero",
			monster:  baseMonster,
			newHP:    0,
			wantHP:   0,
			wantName: "Orc",
		},
		{
			name:     "set negative hit points (should clamp to 0)",
			monster:  baseMonster,
			newHP:    -3,
			wantHP:   0,
			wantName: "Orc",
		},
		{
			name:     "set hit points above maximum (should clamp)",
			monster:  baseMonster,
			newHP:    20,
			wantHP:   15,
			wantName: "Orc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mon := tt.monster
			mon.SetHitPoints(tt.newHP)

			if got := mon.GetHitPoints(); got != tt.wantHP {
				t.Errorf("Monster.GetHitPoints() after SetHitPoints(%v) = %v, want %v", tt.newHP, got, tt.wantHP)
			}
			if got := mon.GetName(); got != tt.wantName {
				t.Errorf("Monster.GetName() after SetHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMonster_WithName(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 15}}
	original := NewMonster("Orc", statBlock)

	result := original.WithName("Elite Orc")

	if got := result.GetName(); got != "Elite Orc" {
		t.Errorf("Monster.WithName().GetName() = %v, want %v", got, "Elite Orc")
	}
	if got := result.GetHitPoints(); got != 15 {
		t.Errorf("Monster.WithName().GetHitPoints() = %v, want %v", got, 15)
	}
	if got := result.GetMaximumHitPoints(); got != 15 {
		t.Errorf("Monster.WithName().GetMaximumHitPoints() = %v, want %v", got, 15)
	}
}

func TestCharacter_AdjustHitPoints(t *testing.T) {
	tests := []struct {
		name       string
		char       Character
		adjustment int
		wantHP     int
		wantName   string
	}{
		{
			name:       "healing (positive adjustment)",
			char:       NewCharacterWithHealth("Test", 20).WithHealth(10, 20),
			adjustment: 5,
			wantHP:     15,
			wantName:   "Test",
		},
		{
			name:       "damage (negative adjustment)",
			char:       NewCharacterWithHealth("Test", 20),
			adjustment: -8,
			wantHP:     12,
			wantName:   "Test",
		},
		{
			name:       "healing beyond maximum (should clamp)",
			char:       NewCharacterWithHealth("Test", 20).WithHealth(18, 20),
			adjustment: 5,
			wantHP:     20,
			wantName:   "Test",
		},
		{
			name:       "damage below zero (should clamp to 0)",
			char:       NewCharacterWithHealth("Test", 20).WithHealth(5, 20),
			adjustment: -10,
			wantHP:     0,
			wantName:   "Test",
		},
		{
			name:       "zero adjustment (no change)",
			char:       NewCharacterWithHealth("Test", 20).WithHealth(12, 20),
			adjustment: 0,
			wantHP:     12,
			wantName:   "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := tt.char
			actualAdjustment := char.AdjustHitPoints(tt.adjustment)

			// Test that we get the actual adjustment amount
			_ = actualAdjustment // We could test this if desired

			if got := char.GetHitPoints(); got != tt.wantHP {
				t.Errorf("Character.GetHitPoints() after AdjustHitPoints(%v) = %v, want %v", tt.adjustment, got, tt.wantHP)
			}
			if got := char.GetName(); got != tt.wantName {
				t.Errorf("Character.GetName() after AdjustHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestMonster_AdjustHitPoints(t *testing.T) {
	statBlock := StatBlock{HitPoints: HitPoints{Fixed: 20}}
	baseMonster := NewMonster("Orc", statBlock)

	tests := []struct {
		name       string
		monster    Monster
		adjustment int
		wantHP     int
		wantName   string
	}{
		{
			name:       "healing (positive adjustment)",
			monster:    Monster{Name: "Orc", StatBlock: statBlock, HitPoints: 10, MaximumHitPoints: 20},
			adjustment: 5,
			wantHP:     15,
			wantName:   "Orc",
		},
		{
			name:       "damage (negative adjustment)",
			monster:    baseMonster,
			adjustment: -8,
			wantHP:     12,
			wantName:   "Orc",
		},
		{
			name:       "healing beyond maximum (should clamp)",
			monster:    Monster{Name: "Orc", StatBlock: statBlock, HitPoints: 18, MaximumHitPoints: 20},
			adjustment: 5,
			wantHP:     20,
			wantName:   "Orc",
		},
		{
			name:       "damage below zero (should clamp to 0)",
			monster:    Monster{Name: "Orc", StatBlock: statBlock, HitPoints: 5, MaximumHitPoints: 20},
			adjustment: -10,
			wantHP:     0,
			wantName:   "Orc",
		},
		{
			name:       "zero adjustment (no change)",
			monster:    Monster{Name: "Orc", StatBlock: statBlock, HitPoints: 12, MaximumHitPoints: 20},
			adjustment: 0,
			wantHP:     12,
			wantName:   "Orc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monster := tt.monster
			actualAdjustment := monster.AdjustHitPoints(tt.adjustment)

			// Test that we get the actual adjustment amount
			_ = actualAdjustment // We could test this if desired

			if got := monster.GetHitPoints(); got != tt.wantHP {
				t.Errorf("Monster.GetHitPoints() after AdjustHitPoints(%v) = %v, want %v", tt.adjustment, got, tt.wantHP)
			}
			if got := monster.GetName(); got != tt.wantName {
				t.Errorf("Monster.GetName() after AdjustHitPoints() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

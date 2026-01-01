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

	creatures = append(creatures, char)
	creatures = append(creatures, monster)

	expectedNames := []string{"Aragorn", "Orc"}

	for i, creature := range creatures {
		if creature.GetName() != expectedNames[i] {
			t.Errorf("Creature[%d].GetName() = %v, want %v", i, creature.GetName(), expectedNames[i])
		}
	}
}



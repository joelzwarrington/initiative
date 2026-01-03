package dnd

import (
	"testing"
	"time"
)

// createTestEncounter creates a basic encounter with 3 initiative groups for testing
func createTestEncounter() *Encounter {
	return NewEncounter("", time.Time{}, []*InitiativeGroup{
		NewInitiativeGroup(15, []*Creature{}),
		NewInitiativeGroup(10, []*Creature{}),
		NewInitiativeGroup(5, []*Creature{}),
	})
}

// createEncounterAtLastTurn creates an encounter already at the last turn of round 1
func createEncounterAtLastTurn() *Encounter {
	e := createTestEncounter()
	e.AdvanceTurn() // turnIndex: 1
	e.AdvanceTurn() // turnIndex: 2
	return e
}

// createEncounterAtTurn1 creates an encounter at turnIndex: 1
func createEncounterAtTurn1() *Encounter {
	e := createTestEncounter()
	e.AdvanceTurn() // turnIndex: 1
	return e
}

// createEncounterAtRound2 creates an encounter at round 2, turnIndex: 0
func createEncounterAtRound2() *Encounter {
	e := createTestEncounter()
	e.AdvanceTurn() // turnIndex: 1
	e.AdvanceTurn() // turnIndex: 2
	e.AdvanceTurn() // round: 2, turnIndex: 0
	return e
}

func TestEncounter_AdvanceTurn(t *testing.T) {
	tests := []struct {
		name      string
		encounter *Encounter
		wantRound int
		wantTurn  int
	}{
		{
			name:      "advance turn in middle of round",
			encounter: createTestEncounter(),
			wantRound: 1,
			wantTurn:  1,
		},
		{
			name:      "advance to next round",
			encounter: createEncounterAtLastTurn(),
			wantRound: 2,
			wantTurn:  0,
		},
		{
			name:      "empty initiative groups",
			encounter: NewEncounter("", time.Time{}, []*InitiativeGroup{}),
			wantRound: 1,
			wantTurn:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.encounter
			e.AdvanceTurn()
			if e.Round() != tt.wantRound {
				t.Errorf("AdvanceTurn() round = %v, want %v", e.Round(), tt.wantRound)
			}
			if e.TurnIndex() != tt.wantTurn {
				t.Errorf("AdvanceTurn() turnIndex = %v, want %v", e.TurnIndex(), tt.wantTurn)
			}
		})
	}
}

func TestEncounter_PreviousTurn(t *testing.T) {
	tests := []struct {
		name      string
		encounter *Encounter
		wantRound int
		wantTurn  int
	}{
		{
			name:      "previous turn in middle of round",
			encounter: createEncounterAtTurn1(),
			wantRound: 1,
			wantTurn:  0,
		},
		{
			name:      "previous turn to previous round",
			encounter: createEncounterAtRound2(),
			wantRound: 1,
			wantTurn:  2,
		},
		{
			name:      "cannot go before round 1, turn 0",
			encounter: createTestEncounter(),
			wantRound: 1,
			wantTurn:  0,
		},
		{
			name:      "empty initiative groups",
			encounter: NewEncounter("", time.Time{}, []*InitiativeGroup{}),
			wantRound: 1,
			wantTurn:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.encounter
			e.PreviousTurn()
			if e.Round() != tt.wantRound {
				t.Errorf("PreviousTurn() round = %v, want %v", e.Round(), tt.wantRound)
			}
			if e.TurnIndex() != tt.wantTurn {
				t.Errorf("PreviousTurn() turnIndex = %v, want %v", e.TurnIndex(), tt.wantTurn)
			}
		})
	}
}

func TestEncounter_TurnIndex(t *testing.T) {
	e := NewEncounter("", time.Time{}, []*InitiativeGroup{})
	if got := e.TurnIndex(); got != 0 {
		t.Errorf("TurnIndex() = %v, want %v", got, 0)
	}
}

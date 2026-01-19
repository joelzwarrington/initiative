package dnd

import "time"

// Encounter represents a combat session with multiple participants
type Encounter struct {
	summary string

	startedAt time.Time
	endedAt   time.Time

	round     int
	turnIndex int
	groups    []*InitiativeGroup
}

// InitiativeGroup represents creatures acting on the same initiative count
type InitiativeGroup struct {
	initiative int
	creatures  []*Creature
}

// InitiativeEntry represents a creature entry in the initiative list field
// Used as form data before converting to InitiativeGroups
type InitiativeEntry struct {
	CreatureType string     // "character" or "monster"
	CreatureID   string     // UUID for character, "source:name" for monster
	Name         string     // Display name (custom for monsters)
	Initiative   int        // Initiative value (0 = not set)
	Quantity     int        // 1 for characters, configurable for monsters
	StatBlock    *StatBlock // For monsters only (nil for characters)
}

// Initiative returns the initiative value for this group
func (g *InitiativeGroup) Initiative() int {
	return g.initiative
}

// Creatures returns the creatures in this group
func (g *InitiativeGroup) Creatures() []*Creature {
	return g.creatures
}

// NewInitiativeGroup creates a new initiative group with the given initiative and creatures
func NewInitiativeGroup(initiative int, creatures []*Creature) *InitiativeGroup {
	return &InitiativeGroup{
		initiative: initiative,
		creatures:  creatures,
	}
}

// NewEncounter creates a new encounter with the given parameters
func NewEncounter(summary string, startedAt time.Time, groups []*InitiativeGroup) *Encounter {
	return &Encounter{
		summary:   summary,
		startedAt: startedAt,
		round:     1,
		turnIndex: 0,
		groups:    groups,
	}
}

// Summary returns the encounter's summary
func (e *Encounter) Summary() string {
	return e.summary
}

// StartedAt returns when the encounter started
func (e *Encounter) StartedAt() time.Time {
	return e.startedAt
}

// EndedAt returns when the encounter ended
func (e *Encounter) EndedAt() time.Time {
	return e.endedAt
}

// Round returns the current round number
func (e *Encounter) Round() int {
	return e.round
}

// TurnIndex returns the current turn index
func (e *Encounter) TurnIndex() int {
	return e.turnIndex
}

// InitiativeGroups returns the initiative groups
func (e *Encounter) InitiativeGroups() []*InitiativeGroup {
	return e.groups
}

// UpdateCreature updates a creature in the encounter and returns the actual adjustment made
func (e *Encounter) UpdateCreature(groupIndex, creatureIndex int, adjustment int) int {
	if groupIndex >= len(e.groups) {
		return 0
	}
	creatures := e.groups[groupIndex].Creatures()
	if creatureIndex >= len(creatures) {
		return 0
	}
	return (*creatures[creatureIndex]).AdjustHP(adjustment)
}

// AdvanceTurn moves to the next creature's turn, advancing round if necessary
func (e *Encounter) AdvanceTurn() {
	if len(e.groups) == 0 {
		return
	}

	e.turnIndex++
	if e.turnIndex >= len(e.groups) {
		e.turnIndex = 0
		e.round++
	}
}

// PreviousTurn moves to the previous creature's turn, going back a round if necessary
func (e *Encounter) PreviousTurn() {
	if len(e.groups) == 0 {
		return
	}

	// Don't allow going back if we're at round 1, turn 0
	if e.round == 1 && e.turnIndex == 0 {
		return
	}

	e.turnIndex--
	if e.turnIndex < 0 {
		e.turnIndex = len(e.groups) - 1
		e.round--
	}
}

// Package dnd provides types and data structures to represent Dungeons & Dragons 5e game elements.
package dnd

import (
	"time"
)

// Encounter represents a combat session with multiple participants
type Encounter struct {
	summary string

	startedAt time.Time
	endedAt   time.Time

	round            int
	turnIndex        int
	initiativeGroups []InitiativeGroup
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
func (e *Encounter) InitiativeGroups() []InitiativeGroup {
	return e.initiativeGroups
}

// NewEncounter creates a new encounter with the given parameters
func NewEncounter(summary string, round int, turnIndex int, groups []InitiativeGroup) Encounter {
	return Encounter{
		summary:          summary,
		round:            round,
		turnIndex:        turnIndex,
		initiativeGroups: groups,
	}
}

// WithInitiativeGroups returns a copy of the encounter with different initiative groups
func (e Encounter) WithInitiativeGroups(groups []InitiativeGroup) Encounter {
	return Encounter{
		summary:          e.summary,
		startedAt:        e.startedAt,
		endedAt:          e.endedAt,
		round:            e.round,
		turnIndex:        e.turnIndex,
		initiativeGroups: groups,
	}
}

// WithStartedAt returns a copy of the encounter with a different start time
func (e Encounter) WithStartedAt(startTime time.Time) Encounter {
	return Encounter{
		summary:          e.summary,
		startedAt:        startTime,
		endedAt:          e.endedAt,
		round:            e.round,
		turnIndex:        e.turnIndex,
		initiativeGroups: e.initiativeGroups,
	}
}

// UpdateCreature updates a creature in the encounter and returns the actual adjustment made
func (e *Encounter) UpdateCreature(groupIndex, creatureIndex int, adjustment int) int {
	if groupIndex >= len(e.initiativeGroups) || creatureIndex >= len(e.initiativeGroups[groupIndex].Creatures) {
		return 0
	}
	return e.initiativeGroups[groupIndex].Creatures[creatureIndex].AdjustHitPoints(adjustment)
}

// AdvanceTurn moves to the next creature's turn, advancing round if necessary
func (e *Encounter) AdvanceTurn() {
	if len(e.initiativeGroups) == 0 {
		return
	}

	e.turnIndex++
	if e.turnIndex >= len(e.initiativeGroups) {
		e.turnIndex = 0
		e.round++
	}
}

// PreviousTurn moves to the previous creature's turn, going back a round if necessary
func (e *Encounter) PreviousTurn() {
	if len(e.initiativeGroups) == 0 {
		return
	}

	// Don't allow going back if we're at round 1, turn 0
	if e.round == 1 && e.turnIndex == 0 {
		return
	}

	e.turnIndex--
	if e.turnIndex < 0 {
		e.turnIndex = len(e.initiativeGroups) - 1
		e.round--
	}
}

// InitiativeGroup represents creatures acting on the same initiative count
type InitiativeGroup struct {
	Initiative int
	Creatures  []Creature
}

// Creature represents any entity with stats that can participate in combat
type Creature interface {
	Name() string
	HitPoints() int
	MaximumHitPoints() int
	ArmorClass() int
	SetHitPoints(hp int)
	AdjustHitPoints(amount int) int
}

// Ability represents a D&D ability score with its modifier
type Ability struct {
	Score    int `yaml:"score" json:"score"`
	Modifier int `yaml:"modifier" json:"modifier"`
}

// Challenge represents a monster's challenge rating and XP value
type Challenge struct {
	Rating           int `yaml:"rating" json:"rating"`
	ExperiencePoints int `yaml:"experience_points" json:"experience_points"`
}

// ArmorClass represents a creature's AC value and type
type ArmorClass struct {
	Value int    `yaml:"value" json:"value"`
	Type  string `yaml:"type" json:"type"`
}

// HitPoints represents both fixed HP and the dice roll used to generate it
type HitPoints struct {
	Fixed int    `yaml:"fixed" json:"fixed"`
	Roll  string `yaml:"roll" json:"roll"`
}

// StatBlock contains all mechanical information for a D&D monster
type StatBlock struct {
	Meta string `yaml:"meta" json:"meta"`

	ArmorClass ArmorClass `yaml:"armor_class" json:"armor_class"`
	HitPoints  HitPoints  `yaml:"hit_points" json:"hit_points"`
	Speed      string     `yaml:"speed" json:"speed"`

	STR Ability `yaml:"str" json:"str"`
	DEX Ability `yaml:"dex" json:"dex"`
	CON Ability `yaml:"con" json:"con"`
	INT Ability `yaml:"int" json:"int"`
	WIS Ability `yaml:"wis" json:"wis"`
	CHA Ability `yaml:"cha" json:"cha"`

	SavingThrows map[string]int `yaml:"saving_throws,omitempty" json:"saving_throws,omitempty"`
	Skills       map[string]int `yaml:"skills,omitempty" json:"skills,omitempty"`
	Senses       string         `yaml:"senses,omitempty" json:"senses,omitempty"`
	Languages    string         `yaml:"languages,omitempty" json:"languages,omitempty"`
	Challenge    Challenge      `yaml:"challenge" json:"challenge"`

	Traits           string `yaml:"traits,omitempty" json:"traits,omitempty"`
	Actions          string `yaml:"actions,omitempty" json:"actions,omitempty"`
	LegendaryActions string `yaml:"legendary_actions,omitempty" json:"legendary_actions,omitempty"`
}

// SourceMeta contains metadata about a monster source book
type SourceMeta struct {
	Name    string `yaml:"name" json:"name"`
	Key     string `yaml:"key" json:"key"`
	Version string `yaml:"version" json:"version"`
}

// Source represents a collection of monsters from a specific sourcebook
type Source struct {
	Meta     SourceMeta `yaml:"meta" json:"meta"`
	Monsters []Monster  `yaml:"monsters" json:"monsters"`
}

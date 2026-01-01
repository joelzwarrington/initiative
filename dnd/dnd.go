// Package dnd provides types and data structures to represent Dungeons & Dragons 5e game elements.
package dnd

import "time"

// Core game types
type Encounter struct {
	Summary string

	StartedAt time.Time
	EndedAt   time.Time

	Round            int
	turnIndex        int
	InitiativeGroups []InitiativeGroup
}

func (e *Encounter) GetTurnIndex() int {
	return e.turnIndex
}

func (e *Encounter) AdvanceTurn() {
	if len(e.InitiativeGroups) == 0 {
		return
	}

	e.turnIndex++
	if e.turnIndex >= len(e.InitiativeGroups) {
		e.turnIndex = 0
		e.Round++
	}
}

func (e *Encounter) PreviousTurn() {
	if len(e.InitiativeGroups) == 0 {
		return
	}

	// Don't allow going back if we're at round 1, turn 0
	if e.Round == 1 && e.turnIndex == 0 {
		return
	}

	e.turnIndex--
	if e.turnIndex < 0 {
		e.turnIndex = len(e.InitiativeGroups) - 1
		e.Round--
	}
}

type InitiativeGroup struct {
	Initiative int
	Creatures  []Creature
}

type Creature interface {
	GetName() string
	GetHitPoints() int
	GetMaximumHitPoints() int
	SetHitPoints(hp int)
	AdjustHitPoints(amount int) int
}

// Character represents a player character
type Character struct {
	Name             string `yaml:"name" json:"name"`
	HitPoints        int    `yaml:"hit_points" json:"hit_points"`
	MaximumHitPoints int    `yaml:"maximum_hit_points" json:"maximum_hit_points"`
}

func (c Character) GetName() string {
	return c.Name
}

func (c Character) GetHitPoints() int {
	return c.HitPoints
}

func (c Character) GetMaximumHitPoints() int {
	return c.MaximumHitPoints
}

func (c *Character) SetHitPoints(hp int) {
	if hp < 0 {
		hp = 0
	}
	if hp > c.MaximumHitPoints {
		hp = c.MaximumHitPoints
	}
	c.HitPoints = hp
}

func (c *Character) AdjustHitPoints(amount int) int {
	oldHP := c.HitPoints
	newHP := c.HitPoints + amount
	c.SetHitPoints(newHP)
	return c.HitPoints - oldHP
}

func NewCharacter(name string) Character {
	return Character{Name: name, HitPoints: 0, MaximumHitPoints: 0}
}

func NewCharacterWithHealth(name string, maxHP int) Character {
	return Character{Name: name, HitPoints: maxHP, MaximumHitPoints: maxHP}
}

func (c Character) WithName(name string) Character {
	return Character{Name: name, HitPoints: c.HitPoints, MaximumHitPoints: c.MaximumHitPoints}
}

func (c Character) WithHealth(hp, maxHP int) Character {
	if hp < 0 {
		hp = 0
	}
	if hp > maxHP {
		hp = maxHP
	}
	return Character{Name: c.Name, HitPoints: hp, MaximumHitPoints: maxHP}
}

// Monster represents a creature from the SRD
type Monster struct {
	Name             string    `yaml:"name" json:"name"`
	StatBlock        StatBlock `yaml:"stat_block" json:"stat_block"`
	HitPoints        int       `yaml:"-" json:"-"`
	MaximumHitPoints int       `yaml:"-" json:"-"`
}

func (m Monster) GetName() string {
	return m.Name
}

func (m Monster) GetHitPoints() int {
	return m.HitPoints
}

func (m Monster) GetMaximumHitPoints() int {
	return m.MaximumHitPoints
}

func (m *Monster) SetHitPoints(hp int) {
	if hp < 0 {
		hp = 0
	}
	if hp > m.MaximumHitPoints {
		hp = m.MaximumHitPoints
	}
	m.HitPoints = hp
}

func (m *Monster) AdjustHitPoints(amount int) int {
	oldHP := m.HitPoints
	newHP := m.HitPoints + amount
	m.SetHitPoints(newHP)
	return m.HitPoints - oldHP
}

func NewMonster(name string, statBlock StatBlock) Monster {
	maxHP := statBlock.HitPoints.Fixed
	return Monster{
		Name:             name,
		StatBlock:        statBlock,
		HitPoints:        maxHP,
		MaximumHitPoints: maxHP,
	}
}

func (m Monster) WithName(name string) Monster {
	return Monster{
		Name:             name,
		StatBlock:        m.StatBlock,
		HitPoints:        m.HitPoints,
		MaximumHitPoints: m.MaximumHitPoints,
	}
}

// SRD data types
type Ability struct {
	Score    int `yaml:"score" json:"score"`
	Modifier int `yaml:"modifier" json:"modifier"`
}

type Challenge struct {
	Rating           int `yaml:"rating" json:"rating"`
	ExperiencePoints int `yaml:"experience_points" json:"experience_points"`
}

type ArmorClass struct {
	Value int    `yaml:"value" json:"value"`
	Type  string `yaml:"type" json:"type"`
}

type HitPoints struct {
	Fixed int    `yaml:"fixed" json:"fixed"`
	Roll  string `yaml:"roll" json:"roll"`
}

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

type SourceMeta struct {
	Name    string `yaml:"name" json:"name"`
	Key     string `yaml:"key" json:"key"`
	Version string `yaml:"version" json:"version"`
}

type Source struct {
	Meta     SourceMeta `yaml:"meta" json:"meta"`
	Monsters []Monster  `yaml:"monsters" json:"monsters"`
}

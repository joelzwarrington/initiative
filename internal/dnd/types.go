package dnd

import "time"

// Core game types
type Encounter struct {
	Summary string

	StartedAt time.Time
	EndedAt   time.Time

	Round           int
	InitiativeGroups []InitiativeGroup
}

type InitiativeGroup struct {
	Initiative int
	Creatures  []Creature
}

type Creature interface {
	Name() string
}

// Character represents a player character
type Character struct {
	name string
}

func (c Character) Name() string {
	return c.name
}

func NewCharacter(name string) Character {
	return Character{name: name}
}

func (c Character) WithName(name string) Character {
	return Character{name: name}
}

// Monster represents a creature from the SRD
type Monster struct {
	MonsterName string    `yaml:"name" json:"name"`
	StatBlock   StatBlock `yaml:"stat_block" json:"stat_block"`
}

func (m Monster) Name() string {
	return m.MonsterName
}

func NewMonster(name string, statBlock StatBlock) Monster {
	return Monster{
		MonsterName: name,
		StatBlock:   statBlock,
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
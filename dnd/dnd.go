// Package dnd provides types and data structures to represent Dungeons & Dragons 5e game elements.
package dnd

// Creature represents any entity with stats that can participate in combat
type Creature interface {
	Name() string

	HP() int
	MaxHP() int

	AC() int

	SetHP(hp int)
	AdjustHP(amount int) int
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

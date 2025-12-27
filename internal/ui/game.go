package ui

import "time"

type Encounter struct {
	Summary string

	StartedAt time.Time
	EndedAt   time.Time

	Round          int
	IniativeGroups []IniativeGroup
}

type IniativeGroup struct {
	Iniative  int
	Creatures []Creature
}

type Creature interface {
	Name() string
}

var _ Creature = (*Monster)(nil)

type Monster struct {
	name      string
	statBlock StatBlock
}

func (m Monster) Name() string {
	return m.name
}

var _ Creature = (*Character)(nil)

type Character struct {
	name string
}

func (c Character) Name() string {
	return c.name
}

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

type SRDMetadata struct {
	Name    string `yaml:"name" json:"name"`
	Key     string `yaml:"key" json:"key"`
	Version string `yaml:"version" json:"version"`
}

type SRDMonster struct {
	Name      string    `yaml:"name" json:"name"`
	StatBlock StatBlock `yaml:"stat_block" json:"stat_block"`
}

type SRDDocument struct {
	Meta     SRDMetadata  `yaml:"meta" json:"meta"`
	Monsters []SRDMonster `yaml:"monsters" json:"monsters"`
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

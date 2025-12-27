package ui

import "time"

type Encounter struct {
	Summary string

	StartedAt time.Time
	EndedAt   time.Time

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
	name string
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

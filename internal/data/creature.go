package data

type Creature interface {
	GetName() string
}

var _ Creature = (*Character)(nil)

type Character struct {
	Name string
}

func (c *Character) GetName() string {
	return c.Name
}

var _ Creature = (*NPC)(nil)

type NPC struct {
	Name string
}

func (n *NPC) GetName() string {
	return n.Name
}

var _ Creature = (*Monster)(nil)

type Monster struct {
	Name string
}

func (m *Monster) GetName() string {
	return m.Name
}

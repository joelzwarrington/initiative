package dnd

import "encoding/json"

// Monster represents a creature from the SRD
type Monster struct {
	name      string
	statBlock StatBlock
	hp        int
	maxHP     int
}

// NewMonster creates a monster with the given name and stat block at full health
func NewMonster(name string, statBlock StatBlock) *Monster {
	maxHP := statBlock.HitPoints.Fixed
	return &Monster{
		name:      name,
		statBlock: statBlock,
		hp:        maxHP,
		maxHP:     maxHP,
	}
}

// WithName returns a copy of the monster with a different name
func (m *Monster) WithName(name string) *Monster {
	return &Monster{
		name:      name,
		statBlock: m.statBlock,
		hp:        m.hp,
		maxHP:     m.maxHP,
	}
}

// Name returns the monster's name
func (m *Monster) Name() string {
	return m.name
}

// HP returns the monster's current hit points
func (m *Monster) HP() int {
	return m.hp
}

// MaxHP returns the monster's maximum hit points
func (m *Monster) MaxHP() int {
	return m.maxHP
}

// AC returns the monster's armor class
func (m *Monster) AC() int {
	return m.statBlock.ArmorClass.Value
}

// StatBlock returns the monster's stat block
func (m *Monster) StatBlock() StatBlock {
	return m.statBlock
}

// SetHP sets the monster's current hit points, clamping to valid range
func (m *Monster) SetHP(hp int) {
	if hp < 0 {
		hp = 0
	}
	if hp > m.maxHP {
		hp = m.maxHP
	}
	m.hp = hp
}

// AdjustHP modifies hit points by the given amount and returns actual change
func (m *Monster) AdjustHP(amount int) int {
	oldHP := m.hp
	newHP := m.hp + amount
	m.SetHP(newHP)
	return m.hp - oldHP
}

// Speed returns the monster's speed
func (m *Monster) Speed() string {
	return m.statBlock.Speed
}

// STR returns the monster's Strength ability
func (m *Monster) STR() Ability {
	return m.statBlock.STR
}

// DEX returns the monster's Dexterity ability
func (m *Monster) DEX() Ability {
	return m.statBlock.DEX
}

// CON returns the monster's Constitution ability
func (m *Monster) CON() Ability {
	return m.statBlock.CON
}

// INT returns the monster's Intelligence ability
func (m *Monster) INT() Ability {
	return m.statBlock.INT
}

// WIS returns the monster's Wisdom ability
func (m *Monster) WIS() Ability {
	return m.statBlock.WIS
}

// CHA returns the monster's Charisma ability
func (m *Monster) CHA() Ability {
	return m.statBlock.CHA
}

// --- serialization ---

// monster is used to serialize and deserialize Monster details
type monster struct {
	Name      string    `yaml:"name" json:"name"`
	StatBlock StatBlock `yaml:"stat_block" json:"stat_block"`
}

// MarshalYAML implements the yaml.Marshaler interface
func (m *Monster) MarshalYAML() (any, error) {
	return monster{
		Name:      m.name,
		StatBlock: m.statBlock,
	}, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (m *Monster) UnmarshalYAML(unmarshal func(any) error) error {
	var data monster
	if err := unmarshal(&data); err != nil {
		return err
	}
	maxHP := data.StatBlock.HitPoints.Fixed
	*m = Monster{
		name:      data.Name,
		statBlock: data.StatBlock,
		hp:        maxHP,
		maxHP:     maxHP,
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (m *Monster) MarshalJSON() ([]byte, error) {
	data := monster{
		Name:      m.name,
		StatBlock: m.statBlock,
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *Monster) UnmarshalJSON(b []byte) error {
	var data monster
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	maxHP := data.StatBlock.HitPoints.Fixed
	*m = Monster{
		name:      data.Name,
		statBlock: data.StatBlock,
		hp:        maxHP,
		maxHP:     maxHP,
	}
	return nil
}

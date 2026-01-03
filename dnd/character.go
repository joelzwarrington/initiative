package dnd

import "encoding/json"

// Character represents a player character (PC) or non-player character (NPC)
type Character struct {
	name string

	hp    int
	maxHP int

	speed string

	str Ability
	dex Ability
	con Ability
	int Ability
	wis Ability
	cha Ability
}

// NewCharacter creates a character with the given name
func NewCharacter(name string) *Character {
	return &Character{name: name}
}

// WithName returns a copy of the character with a different name
func (c *Character) WithName(name string) *Character {
	return &Character{name: name, hp: c.hp, maxHP: c.maxHP}
}

// WithHealth returns a copy of the character with different hit points
func (c *Character) WithHealth(hp, maxHP int) *Character {
	if hp < 0 {
		hp = 0
	}
	if hp > maxHP {
		hp = maxHP
	}
	return &Character{name: c.name, hp: hp, maxHP: maxHP}
}

// Name returns the character's name
func (c *Character) Name() string {
	return c.name
}

// HP returns the character's current hit points
func (c *Character) HP() int {
	return c.hp
}

// MaxHP returns the character's maximum hit points
func (c *Character) MaxHP() int {
	return c.maxHP
}

// AC returns the character's armor class
func (c *Character) AC() int {
	return 0 // TODO: Implement character AC system
}

// SetHP sets the character's current hit points, clamping to valid range
func (c *Character) SetHP(hp int) {
	if hp < 0 {
		hp = 0
	}
	if hp > c.maxHP {
		hp = c.maxHP
	}
	c.hp = hp
}

// AdjustHP modifies hit points by the given amount and returns actual change
func (c *Character) AdjustHP(amount int) int {
	oldHP := c.hp
	newHP := c.hp + amount
	c.SetHP(newHP)
	return c.hp - oldHP
}

// Speed returns the character's speed
func (c *Character) Speed() string {
	return c.speed
}

// STR returns the character's Strength ability
func (c *Character) STR() Ability {
	return c.str
}

// DEX returns the character's Dexterity ability
func (c *Character) DEX() Ability {
	return c.dex
}

// CON returns the character's Constitution ability
func (c *Character) CON() Ability {
	return c.con
}

// INT returns the character's Intelligence ability
func (c *Character) INT() Ability {
	return c.int
}

// WIS returns the character's Wisdom ability
func (c *Character) WIS() Ability {
	return c.wis
}

// CHA returns the character's Charisma ability
func (c *Character) CHA() Ability {
	return c.cha
}

// --- serialization ---

// character is used to serialize and deserialize Character details
type character struct {
	Name             string `yaml:"name" json:"name"`
	HitPoints        int    `yaml:"hit_points" json:"hit_points"`
	MaximumHitPoints int    `yaml:"maximum_hit_points" json:"maximum_hit_points"`
}

// MarshalYAML implements the yaml.Marshaler interface
func (c *Character) MarshalYAML() (any, error) {
	return character{
		Name:             c.name,
		HitPoints:        c.hp,
		MaximumHitPoints: c.maxHP,
	}, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *Character) UnmarshalYAML(unmarshal func(any) error) error {
	var data character
	if err := unmarshal(&data); err != nil {
		return err
	}
	*c = Character{
		name:  data.Name,
		hp:    data.HitPoints,
		maxHP: data.MaximumHitPoints,
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (c *Character) MarshalJSON() ([]byte, error) {
	data := character{
		Name:             c.name,
		HitPoints:        c.hp,
		MaximumHitPoints: c.maxHP,
	}
	return json.Marshal(data)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Character) UnmarshalJSON(b []byte) error {
	var data character
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	*c = Character{
		name:  data.Name,
		hp:    data.HitPoints,
		maxHP: data.MaximumHitPoints,
	}
	return nil
}

package data

type Game struct {
	Name string `yaml:"name"`

	Characters []Character `yaml:"characters,omitempty"`
	NPCs       []NPC       `yaml:"npcs,omitempty"`
}

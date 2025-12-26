package data

type Game struct {
	Name string `yaml:"name"`

	Characters map[string]Character `yaml:"characters,omitempty"`
	NPCs       map[string]NPC       `yaml:"npcs,omitempty"`
}

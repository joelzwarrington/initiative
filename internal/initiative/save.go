package initiative

import (
	"os"

	"github.com/joelzwarrington/initiative/dnd"
	"gopkg.in/yaml.v3"
)

type Game struct {
	filepath string `yaml:"-" json:"-"`

	Characters map[string]dnd.Character `yaml:"characters" json:"characters"`
	sources    map[string]*dnd.Source   `yaml:"-" json:"-"`
}

func LoadGame(filepath string, includeSRD bool) (*Game, error) {
	game := &Game{
		Characters: make(map[string]dnd.Character),
		sources:    make(map[string]*dnd.Source),
		filepath:   filepath,
	}

	if _, err := os.Stat(filepath); err == nil {
		if err := game.load(); err != nil {
			return nil, err
		}
	}

	if includeSRD {
		srd, err := dnd.GetSystemReferenceSource()
		if err == nil && srd != nil {
			game.sources[srd.Meta.Key] = srd
		}
	}

	return game, nil
}

func (g *Game) Save() error {
	yamlData, err := yaml.Marshal(g)
	if err != nil {
		return err
	}

	return os.WriteFile(g.filepath, yamlData, 0644)
}

func (g *Game) load() error {
	data, err := os.ReadFile(g.filepath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, g); err != nil {
		return err
	}

	if g.Characters == nil {
		g.Characters = make(map[string]dnd.Character)
	}
	if g.sources == nil {
		g.sources = make(map[string]*dnd.Source)
	}

	return nil
}

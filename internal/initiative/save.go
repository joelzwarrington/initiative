package initiative

import (
	"os"

	"github.com/joelzwarrington/initiative/dnd"
	"gopkg.in/yaml.v3"
)

type game struct {
	filepath string
	sources  map[string]*dnd.Source

	Characters map[string]*dnd.Character `yaml:"characters" json:"characters"`
}

func LoadGame(filepath string, sourceSpecs []string) (*game, error) {
	g := &game{
		Characters: make(map[string]*dnd.Character),
		sources:    make(map[string]*dnd.Source),
		filepath:   filepath,
	}

	if _, err := os.Stat(filepath); err == nil {
		if err := g.load(); err != nil {
			return nil, err
		}
	}

	// Load sources (embedded or file paths)
	for _, sourceSpec := range sourceSpecs {
		source, err := dnd.LoadSource(sourceSpec)
		if err != nil {
			return nil, err
		}
		if source != nil {
			g.sources[source.Meta.Key] = source
		}
	}

	return g, nil
}

func (g *game) Save() error {
	yamlData, err := yaml.Marshal(g)
	if err != nil {
		return err
	}

	return os.WriteFile(g.filepath, yamlData, 0644)
}

func (g *game) load() error {
	data, err := os.ReadFile(g.filepath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, g); err != nil {
		return err
	}

	if g.Characters == nil {
		g.Characters = make(map[string]*dnd.Character)
	}
	if g.sources == nil {
		g.sources = make(map[string]*dnd.Source)
	}

	return nil
}

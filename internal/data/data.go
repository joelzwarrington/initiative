package data

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

type Game struct {
	Name string `yaml:"name"`
}

type Data struct {
	Games    []Game `yaml:"games"`
	filePath string `yaml:"-"`
}

func Load(filePath string) (*Data, error) {
	data := &Data{filePath: filePath}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return data, nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (d *Data) Save() error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(1)

	if err := encoder.Encode(d); err != nil {
		return err
	}

	encoder.Close()
	return os.WriteFile(d.filePath, buf.Bytes(), 0644)
}

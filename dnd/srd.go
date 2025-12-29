package dnd

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed srd.yaml
var srdYAML []byte

// GetSystemReferenceSource returns the complete SRD document
func GetSystemReferenceSource() (*Source, error) {
	var srd Source
	err := yaml.Unmarshal(srdYAML, &srd)
	if err != nil {
		return nil, err
	}
	return &srd, nil
}
package dnd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed *.yaml
var Sources embed.FS

// LoadSourceFromFS loads a source from the embedded filesystem
func LoadSourceFromFS(path string) (*Source, error) {
	data, err := Sources.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file %s: %w", path, err)
	}

	var source Source
	err = yaml.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal source data from %s: %w", path, err)
	}
	return &source, nil
}

// LoadSourceFromFile loads a source from a file path
func LoadSourceFromFile(filepath string) (*Source, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("source file '%s' not found", filepath)
		}
		return nil, fmt.Errorf("failed to read source file %s: %w", filepath, err)
	}

	var source Source
	err = yaml.Unmarshal(data, &source)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal source data from %s: %w", filepath, err)
	}
	return &source, nil
}

// LoadSource loads a source from either embedded FS (by key) or file path
func LoadSource(sourceSpec string) (*Source, error) {
	// Map common keys to embedded files
	embeddedSources := map[string]string{
		"srd": "srd.yaml",
	}

	// Check if it's a known embedded source key
	if embeddedPath, exists := embeddedSources[sourceSpec]; exists {
		return LoadSourceFromFS(embeddedPath)
	}

	// Try as embedded file path first
	source, err := LoadSourceFromFS(sourceSpec)
	if err == nil {
		return source, nil
	}

	// If that fails, try as external file path
	return LoadSourceFromFile(sourceSpec)
}

// ListAvailableSources returns all available source files in the embedded FS
func ListAvailableSources() ([]string, error) {
	var sources []string
	err := fs.WalkDir(Sources, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path != "." {
			sources = append(sources, path)
		}
		return nil
	})
	return sources, err
}

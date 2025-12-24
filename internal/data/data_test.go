package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_NewFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.yaml")

	data, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load should not error on non-existent file: %v", err)
	}

	if data == nil {
		t.Fatal("Load should return non-nil data")
	}

	if len(data.Games) != 0 {
		t.Error("New file should have empty games slice")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.yaml")

	// Create data with games
	data, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	data.Games = []Game{
		{Name: "Test Game 1"},
		{Name: "Test Game 2"},
	}

	// Save data
	err = data.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load it back
	loadedData, err := Load(testFile)
	if err != nil {
		t.Fatalf("Load after save failed: %v", err)
	}

	if len(loadedData.Games) != 2 {
		t.Errorf("Expected 2 games, got %d", len(loadedData.Games))
	}

	if loadedData.Games[0].Name != "Test Game 1" {
		t.Errorf("Expected 'Test Game 1', got '%s'", loadedData.Games[0].Name)
	}

	if loadedData.Games[1].Name != "Test Game 2" {
		t.Errorf("Expected 'Test Game 2', got '%s'", loadedData.Games[1].Name)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.yaml")

	// Write invalid YAML
	err := os.WriteFile(testFile, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err = Load(testFile)
	if err == nil {
		t.Error("Load should error on invalid YAML")
	}
}

package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestApp(t *testing.T) {
	app := newApp()

	// Test View returns expected message
	if got := app.View(); got != "Hello World!" {
		t.Errorf("View() = %q, want %q", got, "Hello World!")
	}

	// Test Init returns nil
	if cmd := app.Init(); cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestAppUpdate(t *testing.T) {
	app := newApp()

	// Test quit keys return commands
	quitKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'q'}},
		{Type: tea.KeyCtrlC},
	}

	for _, key := range quitKeys {
		_, cmd := app.Update(key)
		if cmd == nil {
			t.Errorf("Update(%v) returned nil command, expected quit", key)
		}
	}

	// Test non-quit key returns no command
	normalKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	model, cmd := app.Update(normalKey)

	if cmd != nil {
		t.Errorf("Update(%v) returned command %v, expected nil", normalKey, cmd)
	}

	if model != app {
		t.Error("Update() should return unchanged model for non-quit keys")
	}
}

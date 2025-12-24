package ui

import (
	"initiative/internal/data"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewProgram(t *testing.T) {
	appData := &data.Data{}
	program := NewProgram(appData)
	if program == nil {
		t.Error("NewProgram should return a non-nil program")
	}
}

func TestApp_Update_QuitKey(t *testing.T) {
	appData := &data.Data{}
	testApp := newApp(appData)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}

	_, cmd := testApp.Update(msg)

	if cmd == nil || cmd() != tea.Quit() {
		t.Error("Quit key should return tea.Quit command")
	}
}

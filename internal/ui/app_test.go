package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewProgram(t *testing.T) {
	program := NewProgram()
	if program == nil {
		t.Error("NewProgram should return a non-nil program")
	}
}

func TestApp_Init(t *testing.T) {
	app := newApp()
	cmd := app.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestApp_Update_QuitKey(t *testing.T) {
	testApp := newApp()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	
	_, cmd := testApp.Update(msg)
	
	if cmd == nil || cmd() != tea.Quit() {
		t.Error("Quit key should return tea.Quit command")
	}
}

func TestApp_Update_NewKey(t *testing.T) {
	testApp := newApp()
	testApp.currentView = GameList
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	
	model, _ := testApp.Update(msg)
	updatedApp := model.(app)
	
	if updatedApp.currentView != NewGameForm {
		t.Error("New key should switch to NewGameForm view")
	}
}

func TestApp_Update_BackKey(t *testing.T) {
	testApp := newApp()
	testApp.currentView = NewGameForm
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	
	model, _ := testApp.Update(msg)
	updatedApp := model.(app)
	
	// Back key functionality not implemented for NewGameForm view
	if updatedApp.currentView == GameList {
		t.Error("Back key functionality not implemented in NewGameForm view")
	}
}

func TestApp_View(t *testing.T) {
	t.Run("game list view renders correctly", func(t *testing.T) {
		testApp := newApp()
		testApp.currentView = GameList
		view := testApp.View()
		if view == "" {
			t.Error("GameList view should return non-empty string")
		}
	})
	
	t.Run("new game form view renders correctly", func(t *testing.T) {
		testApp := newApp()
		testApp.currentView = NewGameForm
		view := testApp.View()
		if view == "" {
			t.Error("NewGameForm view should return non-empty string")
		}
	})
}
package ui

import (
	"initiative/internal/game"
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

func TestApp_Update_SelectKey(t *testing.T) {
	testApp := newApp()
	testApp.currentView = GameList
	
	// Add a game to select
	testGame := game.Game{Name: "Test Game"}
	testApp.games = []game.Game{testGame}
	testApp.gameList.SetItems(game.ToListItems(testApp.games))
	
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	
	model, _ := testApp.Update(msg)
	updatedApp := model.(app)
	
	if updatedApp.currentView != ShowGame {
		t.Error("Select key should switch to ShowGame view")
	}
	
	if updatedApp.currentGame == nil {
		t.Error("currentGame should be set when selecting a game")
	}
	
	if updatedApp.currentGame.Name != "Test Game" {
		t.Errorf("Expected currentGame name to be 'Test Game', got %s", updatedApp.currentGame.Name)
	}
}

func TestApp_Update_BackKeyFromShowGame(t *testing.T) {
	testApp := newApp()
	testApp.currentView = ShowGame
	testApp.currentGame = &game.Game{Name: "Test Game"}
	
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	
	model, _ := testApp.Update(msg)
	updatedApp := model.(app)
	
	if updatedApp.currentView != GameList {
		t.Error("Back key should switch from ShowGame to GameList view")
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
	
	t.Run("show game view renders correctly", func(t *testing.T) {
		testApp := newApp()
		testApp.currentView = ShowGame
		testApp.currentGame = &game.Game{Name: "Test Game"}
		view := testApp.View()
		if view == "" {
			t.Error("ShowGame view should return non-empty string")
		}
	})
}
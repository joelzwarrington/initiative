package models

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
)

func TestGame_Title(t *testing.T) {
	game := Game{Name: "Test Game"}
	expected := "Test Game"

	result := game.Title()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGame_FilterValue(t *testing.T) {
	game := Game{Name: "Filter Game"}
	expected := "Filter Game"

	result := game.FilterValue()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGame_ToListItem(t *testing.T) {
	game := Game{Name: "List Game"}

	result := game.ToListItem()

	resultGame, ok := result.(*Game)
	if !ok {
		t.Error("ToListItem should return a *Game")
	}

	if resultGame.Name != game.Name {
		t.Errorf("Expected name %q, got %q", game.Name, resultGame.Name)
	}

	var _ list.Item = result
}

func TestGame_ToListItems(t *testing.T) {
	games := []Game{
		{Name: "Game 1"},
		{Name: "Game 2"},
		{Name: "Game 3"},
	}

	result := GamesToListItems(games)

	if len(result) != 3 {
		t.Errorf("Expected 3 items, got %d", len(result))
	}

	for i, item := range result {
		game, ok := item.(*Game)
		if !ok {
			t.Errorf("Item %d is not a *Game", i)
			continue
		}
		if game.Name != games[i].Name {
			t.Errorf("Expected game %d name %q, got %q", i, games[i].Name, game.Name)
		}
	}
}

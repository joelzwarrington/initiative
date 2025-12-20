package game

import "github.com/charmbracelet/bubbles/list"

type Game struct {
	Name string
}

func (g *Game) Title() string       { return g.Name }
func (g *Game) FilterValue() string { return g.Name }

func (g *Game) ToListItem() list.Item {
	return g
}

func ToListItems(games []Game) []list.Item {
	items := make([]list.Item, len(games))
	for i, game := range games {
		items[i] = game.ToListItem()
	}
	return items
}

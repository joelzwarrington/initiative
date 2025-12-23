package ui

import (
	"fmt"
	"initiative/internal/game"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	GameList view = iota
	NewGameForm
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type gameDelegate struct{}

func (d gameDelegate) Height() int                             { return 1 }
func (d gameDelegate) Spacing() int                            { return 0 }
func (d gameDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d gameDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	g, ok := listItem.(*game.Game)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, g.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type gameList struct {
	list.Model
	currentGame *game.Game
	games       *[]game.Game
}

func newGameList(currentGame *game.Game, games []game.Game) gameList {
	gameList := gameList{
		Model: list.New(
			game.ToListItems(games),
			gameDelegate{},
			0,
			0,
		),
		currentGame: currentGame,
	}

	gameList.Title = "Games"
	gameList.SetShowTitle(true)
	gameList.SetShowStatusBar(false)
	gameList.SetFilteringEnabled(false)
	gameList.SetShowHelp(true)

	gameList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.New,
		}
	}
	
	gameList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.New,
		}
	}

	return gameList
}

type gameForm struct {
	huh.Form
}

func newGameForm() gameForm {
	gameForm := gameForm{
		Form: *huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Name").
					Validate(func(str string) error {
						if str == "" {
							return nil
						}
						return nil
					}),
			),
		).WithTheme(huh.ThemeCharm()),
	}

	return gameForm
}

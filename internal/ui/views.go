package ui

import (
	"fmt"
	"initiative/internal/models"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
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
	ShowGame
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
	g, ok := listItem.(*models.Game)
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
	currentGame *models.Game
	games       *[]models.Game
}

func newGameList(currentGame *models.Game, games []models.Game) gameList {
	gameList := gameList{
		Model: list.New(
			models.GamesToListItems(games),
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
			keys.Select,
		}
	}

	gameList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.New,
			keys.Select,
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

type showGameKeyMap struct{}

func (k showGameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Back, keys.Quit}
}

func (k showGameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keys.Back, keys.Quit},
	}
}

func (a app) showGame() string {
	var s strings.Builder

	if a.currentGame == nil {
		return "No game selected"
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	title := titleStyle.Render(a.currentGame.Name)
	s.WriteString(title)
	s.WriteString("\n\n")

	helpModel := help.New()
	helpView := helpModel.View(showGameKeyMap{})
	s.WriteString(helpView)

	return s.String()
}

package messages

import (
	"initiative/internal/data"

	tea "github.com/charmbracelet/bubbletea"
)

// Navigation commands that views can return to communicate with the app
type NavigateToGameListMsg struct{}
type NavigateToNewGameFormMsg struct{}
type NavigateToShowGameMsg struct {
	Game *data.Game
}
type SaveDataMsg struct{}

func NavigateToGameList() tea.Cmd {
	return func() tea.Msg { return NavigateToGameListMsg{} }
}

func NavigateToNewGameForm() tea.Cmd {
	return func() tea.Msg { return NavigateToNewGameFormMsg{} }
}

func NavigateToShowGame(game *data.Game) tea.Cmd {
	return func() tea.Msg { return NavigateToShowGameMsg{Game: game} }
}

func SaveData() tea.Cmd {
	return func() tea.Msg { return SaveDataMsg{} }
}
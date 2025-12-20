package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit key.Binding
	Back key.Binding
	Help key.Binding
	New  key.Binding
}

// Define common key bindings
var (
	quitKey = key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
	backKey = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	helpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	newKey = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	)
)

var keys = keyMap{
	Quit: quitKey,
	Back: backKey,
	Help: helpKey,
	New:  newKey,
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type encounterFormState int

const (
	statusNormal encounterFormState = iota
	stateDone
)

type EncounterFormModel struct {
	state encounterFormState
	form  *huh.Form
}

func newEncounterFormModel() EncounterFormModel {
	m := EncounterFormModel{}

	// Build character options from game characters
	var characterOptions []huh.Option[string]
	// if game != nil {
	// 	for _, character := range game.Characters {
	// 		characterOptions = append(characterOptions, huh.NewOption(character.Name, character.Name))
	// 	}
	// }
	characterOptions = append(characterOptions, huh.NewOption("joel", "joel"))

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("summary").
				Title("Encounter Summary").
				Placeholder("Epic battle at the castle gates..."),

			huh.NewMultiSelect[string]().
				Key("characters").
				Title("Characters").
				Description("Choose which characters participate in this encounter").
				Options(characterOptions...),

			huh.NewConfirm().
				Key("done").
				Title("Ready to roll?").
				Affirmative("Yes").
				Negative("Wait, no"),
		),
	)

	return m
}

func (m EncounterFormModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m EncounterFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		m.state = stateDone
	}

	return m, cmd
}

func (m EncounterFormModel) View() string {
	return m.form.View()
}

func main() {
	_, err := tea.NewProgram(newEncounterFormModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}

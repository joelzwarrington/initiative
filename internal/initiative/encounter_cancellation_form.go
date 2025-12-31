package initiative

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type cancellationForm struct {
	form      *huh.Form
	width     int
	height    int
	confirmed bool
}

func newCancellationForm(width, height int) *cancellationForm {
	f := &cancellationForm{
		width:  width,
		height: height,
	}

	f.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure you want to end the encounter?").
				Affirmative("Yes!").
				Negative("No.").
				Value(&f.confirmed),
		),
	).WithWidth(width).WithHeight(height)

	return f
}

func (f *cancellationForm) Init() tea.Cmd {
	return f.form.Init()
}

func (f *cancellationForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := f.form.Update(msg)
	f.form = form.(*huh.Form)

	if f.form.State == huh.StateCompleted {
		if f.confirmed {
			return f, tea.Cmd(func() tea.Msg {
				return cancellationConfirmedMsg{}
			})
		} else {
			return f, tea.Cmd(func() tea.Msg {
				return cancellationCancelledMsg{}
			})
		}
	}

	return f, cmd
}

func (f *cancellationForm) View() string {
	return f.form.View()
}

func (f *cancellationForm) SetSize(width, height int) {
	f.width = width
	f.height = height
	f.form = f.form.WithWidth(width).WithHeight(height)
}

// Messages
type cancellationConfirmedMsg struct{}
type cancellationCancelledMsg struct{}

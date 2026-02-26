package workout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
)

type Model struct {
	User *flexcreek.User
}

func New(user *flexcreek.User) Model {
	return Model{
		User: user,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

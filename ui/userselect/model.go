package userselect

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
	"github.com/ekholme/flexcreek/ui/common"
)

type state int

const (
	browsing state = iota
	creating
)

type item struct {
	user flexcreek.User
}

func (i item) Title() string       { return i.user.Username }
func (i item) Description() string { return "Select to log a workout" }
func (i item) FilterValue() string { return i.user.Username }

type Model struct {
	list    list.Model
	input   textinput.Model
	state   state
	service flexcreek.UserService
	err     error
}

func New(us flexcreek.UserService) Model {
	//initialize list and input
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Who is working out?"

	ti := textinput.New()
	ti.Placeholder = "New Username..."

	return Model{
		list:    l,
		input:   ti,
		service: us,
		state:   browsing,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if m.state == creating {
			switch msg.String() {
			case "enter":
				username := m.input.Value()

				//TODO HANDLE THESE ERRORS
				id, _ := m.service.CreateUser(context.Background(), username)
				user, _ := m.service.GetUserById(context.Background(), id)
				return m, func() tea.Msg {
					return common.UserSelectedMsg{User: *user}
				}

			case "esc":
				m.state = browsing
				return m, nil
			}
		}

		switch msg.String() {
		case "n": //Press 'n' for New User
			m.state = creating
			m.input.Focus()
			return m, nil

		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				return m, func() tea.Msg { return common.UserSelectedMsg{User: i.user} }
			}
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

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
	list         list.Model
	input        textinput.Model
	state        state
	service      flexcreek.UserService
	err          error
	SelectedUser *flexcreek.User
	Quitting     bool
}

func New(us flexcreek.UserService) Model {
	//initialize list and input
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Who is working out?"

	ti := textinput.New()
	ti.Placeholder = "New Username..."

	return Model{
		list:     l,
		input:    ti,
		service:  us,
		state:    browsing,
		Quitting: false,
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		users, err := m.service.GetAllUsers(context.Background())
		if err != nil {
			return common.ErrMsg{Err: err}
		}
		return common.UsersMsg{Users: users}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case common.UsersMsg:
		var items []list.Item
		for _, u := range msg.Users {
			items = append(items, item{user: u})
		}
		m.list.SetItems(items)
		return m, nil
	case common.ErrMsg:
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.Quitting = true
			return m, tea.Quit
		}
		if m.state == creating {
			switch msg.String() {
			case "enter":
				username := m.input.Value()
				_, err := m.service.CreateUser(context.Background(), username)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.state = browsing
				m.input.SetValue("")
				return m, m.Init()

			case "esc":
				m.state = browsing
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "n": //Press 'n' for New User
			m.state = creating
			m.input.Focus()
			return m, nil

		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.SelectedUser = &i.user
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

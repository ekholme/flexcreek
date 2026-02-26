package userselect

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
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

type keyMap struct {
	newUser key.Binding
	enter   key.Binding
	quit    key.Binding
}

var keys = keyMap{
	newUser: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new user"),
	),
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select user"),
	),
	quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

type item struct {
	user flexcreek.User
}

func (i item) Title() string       { return i.user.Username }
func (i item) Description() string { return "" }
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
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.newUser, keys.enter}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{keys.newUser, keys.enter} }

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
		if key.Matches(msg, keys.quit) {
			m.Quitting = true
			return m, tea.Quit
		}
		if m.state == creating {
			switch {
			case key.Matches(msg, keys.enter):
				username := m.input.Value()
				_, err := m.service.CreateUser(context.Background(), username)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.state = browsing
				m.input.SetValue("")
				return m, m.Init()

			case msg.String() == "esc":
				m.state = browsing
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, keys.newUser):
			m.state = creating
			m.input.Focus()
			return m, nil

		case key.Matches(msg, keys.enter):
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

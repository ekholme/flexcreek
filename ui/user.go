package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
)

const (
	stateList sessionState = iota
	stateCreateUser
)

// defining interfaces that the user model currently requires
// a bit overkill to break out into 3 single-method interfaces now, but whatever
type UserProvider interface {
	GetAllUsers(ctx context.Context) ([]*flexcreek.User, error)
}

type UserCreator interface {
	CreateUser(ctx context.Context, username string) (int, error)
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, id int) error
}

type UserStore interface {
	UserProvider
	UserCreator
	UserDeleter
}

type UserModel struct {
	store    UserStore
	list     list.Model
	input    textinput.Model
	state    sessionState
	loading  bool
	err      error
	selected *flexcreek.User
}

// constructor for usermodel
func NewUserModel(s UserStore) UserModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a User"

	//text input stuff
	ti := textinput.New()
	ti.Placeholder = "New Username..."
	ti.Focus()

	return UserModel{
		store:   s,
		list:    l,
		input:   ti,
		state:   stateList,
		loading: true,
	}
}

// a command to fetch users from the database
// we wrap this in a team.Cmd to ensure it's non-blocking
func fetchUsersCmd(s UserStore) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		users, err := s.GetAllUsers(ctx)
		if err != nil {
			return err
		}

		return usersLoadedMsg{users}
	}
}

type usersLoadedMsg struct {
	users []*flexcreek.User
}

type userSelectedMsg struct {
	user *flexcreek.User
}

type userItem struct {
	flexcreek.User
}

func (i userItem) Title() string       { return i.Username }
func (i userItem) Description() string { return "Select to view workouts" }
func (i userItem) FilterValue() string { return i.Username }

//bubbletea requires models to implement 3 methods -- Init(), Update(), and View()

func (m UserModel) Init() tea.Cmd {
	return fetchUsersCmd(m.store)
}

func (m UserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.loading = false
		m.err = msg
		return m, nil

	case usersLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.users))
		for i, u := range msg.users {
			items[i] = userItem{*u}
		}
		m.list.SetItems(items)

	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.list.SetSize(h, v)

		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateCreateUser:
			return m.updateForm(msg) //TODO
		}
	}
	return m, cmd
}

func (m UserModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error()
	}

	switch m.state {
	case stateCreateUser:
		return "\n Create New User \n\n" +
			m.input.View() +
			"\n\n (esc to go back)"

	default:
		if m.loading {
			return " Loading users..."
		}

		return "\n" + m.list.View()
	}

}

// helper update functions for different views
func (m UserModel) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		//prevents user selection if we're in a filtering state
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "n":
			m.state = stateCreateUser
			m.input.Focus()
			return m, nil

		case "enter":
			if i, ok := m.list.SelectedItem().(userItem); ok {
				m.selected = &i.User
				return m, func() tea.Msg { return userSelectedMsg{&i.User} }
			}
		}

	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

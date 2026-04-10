package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
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
	loading  bool
	err      error
	selected *flexcreek.User
}

// constructor for usermodel
func NewUserModel(s UserStore) UserModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a User"

	return UserModel{
		store:   s,
		list:    l,
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

type userItem struct {
	flexcreek.User
}

func (i userItem) Title() string       { return i.Username }
func (i userItem) Description() string { return "Select to view workouts" }
func (i userItem) FilterValue() string { return i.Username }

//bubbletea requires models to implement 3 methods -- Init(), Update(), and View()

func (m UserModel) Init() tea.Cmd {
	//todo
}

func (m UserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case usersLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.users))
		for i, u := range msg.users {
			items[i] = userItem{*u}
		}
		m.list.SetItems(items)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(userItem); ok {
				m.selected = &i.User
				return m, func() tea.Msg { return UserSelectedMsg{i.User} }
			}
		case "n":
			// Logic to switch to a "Create User" form state
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m UserModel) View() string {
	//todo
}

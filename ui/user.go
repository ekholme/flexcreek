package ui

import (
	"context"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
)

//handles all interactions with the user model

type UserModel struct {
	list   list.Model
	ctx    *ProgramContext
	err    error
	loaded bool
}

func NewUserModel(ctx *ProgramContext) UserModel {
	//returning an empty list that we'll populate later
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a User"

	return UserModel{
		list: l,
		ctx:  ctx,
	}
}

// bubbletea's list component requires items to satisfy an interface, so I need to wrap my User struct in something that does this
type userItem struct {
	user *flexcreek.User
}

func (i userItem) Title() string       { return i.user.Username }
func (i userItem) Description() string { return strconv.Itoa(i.user.ID) }
func (i userItem) FilterValue() string { return i.user.Username }

// message types for internal communication
type loadedUsersMsg []list.Item
type errMsg error
type UserSelectMsg struct {
	User *flexcreek.User
}

//bubbletea requires models to satisfy an interface with Init(), Update(), and View() methods

func loadUsers(m UserModel) tea.Msg {
	ctx := context.Background()
	users, err := m.ctx.UserSvc.GetAllUsers(ctx)

	if err != nil {
		return errMsg(err)
	}

	items := make([]list.Item, len(users))

	for i, u := range users {
		items[i] = userItem{user: u}
	}

	return loadedUsersMsg(items)
}

// RESUME HERE -- i should be able to return this, just need to parse out how
func (m UserModel) Init() tea.Cmd {
	return loadUsers()
}

func (m UserModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m UserModel) View() string {
	return m.list.View()
}

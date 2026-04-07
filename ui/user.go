package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
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

//RESUME HERE -- define message types for passing messages around internally

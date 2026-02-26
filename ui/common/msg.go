package common

import "github.com/ekholme/flexcreek"

type UserSelectedMsg struct {
	User flexcreek.User
}

type UsersMsg struct {
	Users []flexcreek.User
}

type ErrMsg struct {
	Err error
}

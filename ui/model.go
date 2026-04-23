package ui

import (
	"github.com/ekholme/flexcreek/sqlite"
)

// root model that manages which view is currently active and delegates bubbletea calls to sub-models

type sessionState int

const (
	stateUserManager sessionState = iota
	stateWorkoutManager
)

type RootModel struct {
	state        sessionState
	userModel    UserModel
	workoutModel WorkoutModel
}

// constructor function
func NewRootModel(s *sqlite.Storage) RootModel {
	//TODO
	return RootModel{}
}

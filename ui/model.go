package ui

import "github.com/ekholme/flexcreek"

// root model that manages which view is currently active and delegates bubbletea calls to sub-models

type sessionState int

const (
	stateUserManager sessionState = iota
	stateWorkoutManager
)

type RootModel struct {
	state          sessionState
	userManager    UserModel
	workoutManager WorkoutModel
	ctx            *ProgramContext
}

//constructor function
func NewRootModel(us flexcreek.UserService, ws flexcreek.WorkoutService) RootModel {
	ctx := NewProgramContext(us, ws)
	um := NewUserModel(ctx)
	wm := NewWorkoutModel()

	return RootModel{
		userManager:    um,
		workoutManager: wm,
		ctx:            ctx,
	}
}

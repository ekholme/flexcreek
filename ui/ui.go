package ui

import (
	"github.com/ekholme/flexcreek"
	"github.com/ekholme/flexcreek/ui/userselect"
)

func New(userService flexcreek.UserService, workoutService flexcreek.WorkoutService) MainModel {
	return MainModel{
		state:          userSelectView,
		userselect:     userselect.New(userService),
		userService:    userService,
		workoutService: workoutService,
	}
}

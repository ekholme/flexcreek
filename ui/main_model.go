package ui

import "github.com/ekholme/flexcreek"

type sessionState int

const (
	userSelectView sessionState = iota
	workoutView
)

type MainModel struct {
	state sessionState

	//add component models here once they're defined

	//services
	userService         flexcreek.UserService
	workoutService      flexcreek.WorkoutService
	activityTypeService flexcreek.ActivityTypeService
}

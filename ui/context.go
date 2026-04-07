package ui

import "github.com/ekholme/flexcreek"

// to include a struct that stores services and handles session state

type ProgramContext struct {
	UserSvc    flexcreek.UserService
	WorkoutSvc flexcreek.WorkoutService
	ActiveUser *flexcreek.User
}

func NewProgramContext(us flexcreek.UserService, ws flexcreek.WorkoutService) *ProgramContext {
	return &ProgramContext{
		UserSvc:    us,
		WorkoutSvc: ws,
	}
}

package server

import (
	"net/http"

	"github.com/ekholme/flexcreek"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//services
	UserService         flexcreek.UserService
	WorkoutService      flexcreek.WorkoutService
	ActivityTypeService flexcreek.ActivityTypeService

	//templates
	//todo

	//logging
	//todo

	//config
	//todo
}

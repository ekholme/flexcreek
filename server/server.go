package server

import (
	"encoding/json"
	"net/http"

	"github.com/ekholme/flexcreek"
)

// define a server type
// this will 'carry' all of our stuff around that we want
type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//listening address
	Addr string

	//services
	MovementService flexcreek.MovementService

	//todo -- add other stuff
	//include services used throughout

}

// constructor to create a new server object
func NewServer(addr string, ms flexcreek.MovementService) *Server {
	return &Server{
		Router: http.NewServeMux(),
		Srvr:   &http.Server{},
		Addr:   addr,

		//services
		MovementService: ms,
	}
}

// function to register routes that are part of the application
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("GET /{$}", s.handleIndex)

	//api routes ----------
	s.Router.HandleFunc("GET /api/v1/movement/{id}", s.handleApiGetMovement)
	s.Router.HandleFunc("POST /api/v1/movement/create", s.handleApiCreateMovement)

	//html routes ---------
	//movement routes
	s.Router.HandleFunc("GET /movement/{id}", s.handleGetMovement)
	s.Router.HandleFunc("GET /movement/create", s.handleCreateMovement)
}

// helper function to run the server
func (s *Server) Run() error {
	s.registerRoutes()

	s.Srvr.Handler = s.Router
	s.Srvr.Addr = s.Addr

	return s.Srvr.ListenAndServe()
}

// utility function to write JSON
func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

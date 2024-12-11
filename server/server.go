package server

import "net/http"

//define a server type
//this will 'carry' all of our stuff around that we want
type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//listening address
	Addr string

	//todo -- add other stuff
	//include services used throughout
}

//constructor to create a new server object
func NewServer(addr string) *Server {
	return &Server{
		Router: http.NewServeMux(),
		Srvr:   &http.Server{},
		Addr:   addr,
	}
}

//function to register routes that are part of the application
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("GET /", s.handleIndex)
}

//helper function to run the server
func (s *Server) Run() error {
	s.registerRoutes()

	s.Srvr.Handler = s.Router
	s.Srvr.Addr = s.Addr

	return s.Srvr.ListenAndServe()
}

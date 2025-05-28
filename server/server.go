package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		Router: http.NewServeMux(),
		Srvr: &http.Server{
			Addr: addr,
		},
	}
}

func (s *Server) Run() {
	s.registerRoutes()
	s.Srvr.Handler = s.Router

	fmt.Printf("Starting flexcreek on %s\n", s.Srvr.Addr)

	s.Srvr.ListenAndServe()
}

// register routes
// this will eventually go somewhere else
func (s *Server) registerRoutes() {
	s.Router.HandleFunc("/", s.HandleIndex)
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello from Flexcreek")

	w.Write(msg)
}

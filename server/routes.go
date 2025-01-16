package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Flex Creek"))
}

// Movement Routes ---------

// api movement routes
func (s *Server) handleApiCreateMovement(w http.ResponseWriter, r *http.Request) {
	//TODO
	writeJSON(w, http.StatusOK, "API route to create a movement")
}

func (s *Server) handleApiGetMovement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := make(map[string]int)

	msg["id"] = id

	writeJSON(w, http.StatusOK, msg)
}

// html movement routes
func (s *Server) handleCreateMovement(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Page to create a new movement..."))
}

func (s *Server) handleGetMovement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display movement with id %d...", id)

	w.Write([]byte(msg))
}

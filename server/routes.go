package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ekholme/flexcreek"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Flex Creek"))
}

// Movement Routes ---------

// api movement routes
func (s *Server) handleApiCreateMovement(w http.ResponseWriter, r *http.Request) {

	var m *flexcreek.Movement

	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	mvID, err := s.MovementService.CreateMovement(m)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	msg := make(map[string]int)

	msg["new movement"] = mvID

	writeJSON(w, http.StatusOK, msg)
}

func (s *Server) handleApiGetMovementByID(w http.ResponseWriter, r *http.Request) {
	//todo -- write a validator here
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	m, err := s.MovementService.GetMovementByID(id)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleApiGetMovementByName(w http.ResponseWriter, r *http.Request) {
	nm := r.PathValue("name")

	m, err := s.MovementService.GetMovementByName(nm)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, m)
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

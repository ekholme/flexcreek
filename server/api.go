package server

import (
	"encoding/json"
	"net/http"

	"github.com/ekholme/flexcreek"
)

// api endpoints -------

// movements

func (s *Server) handleApiCreateMovement(r *http.Request, w http.ResponseWriter) {
	var m *flexcreek.Movement

	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}
}

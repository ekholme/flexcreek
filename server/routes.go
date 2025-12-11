package server

import "net/http"

//html routes -----------
//todo

//api routes ------------

func (s *Server) handleApiIndex(w http.ResponseWriter, r *http.Request) {
	msg := make(map[string]string)

	msg["Hello"] = "World"

	writeJSON(w, http.StatusOK, msg)
}

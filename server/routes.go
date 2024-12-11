package server

import "net/http"

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Flex Creek"))
}

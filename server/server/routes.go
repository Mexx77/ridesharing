package server

import "net/http"

func (s *server) routes() {
	http.HandleFunc("/rides", s.addCORSHeader(s.ridesHandler()))
}

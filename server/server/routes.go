package server

import "net/http"

func (s *server) routes() {
	http.HandleFunc("/rides", s.addCORSHeader(s.ridesHandler()))
	http.HandleFunc("/ride", s.addCORSHeader(s.rideHandler()))
}

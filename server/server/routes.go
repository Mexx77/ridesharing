package server

import "net/http"

func (s *server) routes() {
	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)
	http.HandleFunc("/rides", s.addCORSHeader(s.ridesHandler()))
	http.HandleFunc("/ride/add", s.addCORSHeader(s.loggedInOnly(s.rideAddHandler())))
	http.HandleFunc("/ride/update", s.addCORSHeader(s.loggedInOnly(s.rideUpdateHandler())))
	http.HandleFunc("/ride/delete", s.addCORSHeader(s.loggedInOnly(s.rideDeleteHandler())))
	http.HandleFunc("/ride/unconfirmedRides", s.addCORSHeader(s.adminOnly(s.unconfirmedRidesHandler())))
	http.HandleFunc("/users/authenticate", s.addCORSHeader(s.authenticateHandler()))
	http.HandleFunc("/users/refreshToken", s.addCORSHeader(s.refreshTokenHandler()))
	http.HandleFunc("/users/register", s.addCORSHeader(s.registerHandler()))
}

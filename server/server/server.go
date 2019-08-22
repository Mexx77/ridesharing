package server

import (
	"../logging"
	"database/sql"
	"net/http"
	"os"
)

type server struct {
	database  *sql.DB
}

func NewServer() {
	const databaseFile = "./sqlite.db"
	if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
		panic(err)
	}
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := &server{
		database: db,
	}
	server.startHttpServer()
}

func (s *server) startHttpServer() {
	const port = ":8090"
	httpServer := &http.Server{
		Addr:     port,
		ErrorLog: logging.Error,
	}

	logging.Info.Print("Starting insecure http server on port ", port)
	s.routes()
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}

}

func (s *server) addCORSHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}
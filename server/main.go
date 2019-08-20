package main

import (
    "./logging"
    "./server"
    _ "github.com/mattn/go-sqlite3"
)

func main() {

	logging.Init()
	logging.Info.Print("This is ridesharing!")

	server.NewServer()
}
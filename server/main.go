package main

import (
	"github.com/Mexx77/ridesharing/logging"
	"github.com/Mexx77/ridesharing/server"
)

func main() {

	logging.Init()
	logging.Info.Print("This is ridesharing!")

	server.NewServer()
}
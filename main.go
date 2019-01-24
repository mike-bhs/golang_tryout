package main

import (
	"github.com/mike-bhs/golang_tryout/server"
)

func main() {
	serv := server.InitializeServer()

	defer serv.DB.Close()
	defer serv.MessagingClient.Connection.Close()

	// serv.InitializeGracefulShutdown()

	go serv.StartQueues()

	serv.SetRoutes()
	serv.Engine.Run()
	// serv.RunMigrations()
}

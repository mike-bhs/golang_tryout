package main

import (
	"github.com/golang_tryout/server"
)

func main() {
	serv := server.InitializeServer()

	defer serv.DB.Close()

	serv.RunMigrations()
	serv.SetRoutes()

	serv.Engine.Run()
}

package main

import "github.com/golang_tryout/server"

func main() {
	engine := server.StartServer()
	engine.Run()
}

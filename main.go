package main

import (
	"github.com/golang_tryout/server"
)

func main() {
	serv := server.InitializeServer()

	defer serv.DB.Close()
	defer serv.Amqp.Close()

	serv.SetRoutes()
	serv.RunMigrations()
	serv.StartQueues()

	// var gracefulStop = make(chan os.Signal)
	// signal.Notify(gracefulStop, syscall.SIGTERM)
	// signal.Notify(gracefulStop, syscall.SIGINT)
	//
	// go func() {
	// 	sig := <-gracefulStop
	// 	fmt.Printf("caught sig: %+v", sig)
	// 	fmt.Println("Wait for 2 second to finish processing")
	// 	time.Sleep(2 * time.Second)
	// 	os.Exit(0)
	// }()

	serv.Engine.Run()
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mike-bhs/golang_tryout/server"
	"log"
)

func main() {
	hasDbConnection := make(chan bool)
	hasAmqpConnection := make(chan bool)

	serv := &server.Server{Engine: gin.Default()}

	serv.ConnectToDbAsync(hasDbConnection, "root", "", "localhost:3306", "golang_tryout")
	serv.ConnectToRabbitMQAsync(hasAmqpConnection, "guest", "guest", "localhost:5672")

	isDbConnected := <- hasDbConnection

	if isDbConnected == false {
		log.Println("Failed to connect to db")
		return
	}

	isAmqpConnected := <- hasAmqpConnection

	if isAmqpConnected == false {
		log.Println("Failed to connect to RabbitMQ")
		return
	}

	defer serv.DB.Close()
	defer serv.MessagingClient.Connection.Close()

	serv.InitializeGracefulShutdown()

	go serv.StartQueues()

	serv.SetRoutes()
	err := serv.Engine.Run()

	if err != nil {
		log.Println("Failed to start Gin Server")
		return
	}

	// serv.RunMigrations()
}

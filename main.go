package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mike-bhs/golang_tryout/db"
	"github.com/mike-bhs/golang_tryout/rabbitmq"
	"github.com/mike-bhs/golang_tryout/server"
	"log"
)

func main() {
	dbConfig := &db.Config{User: "root", Password: "", Host: "localhost:3306", DbName: "golang_tryout"}
	database := &db.DataBase{Config: dbConfig}

	amqpConfig := &rabbitmq.Config{User: "guest", Password: "guest", Host: "localhost:5672"}
	amqpClient := &rabbitmq.AmqpClient{Config: amqpConfig}

	serv := &server.Server{Engine: gin.Default(), DataBase: database, AmqpClient: amqpClient}

	go database.EstablishConnection()
	go database.MonitorConnection()
	defer database.CloseConnection()

	amqpClient.EstablishConnection()
	defer amqpClient.CloseConnection()

	amqpClient.InitializeGracefulShutdown()

	//go amqpClient.StartQueues()

	serv.SetRoutes()
	err := serv.Engine.Run()

	if err != nil {
		log.Println("Failed to start Gin Server")
		return
	}

	// serv.RunMigrations()
}

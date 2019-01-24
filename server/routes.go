package server

import (
	"github.com/mike-bhs/golang_tryout/app/models"
)

func (serv *Server) SetRoutes() {

	transactions := serv.Engine.Group("/api/transactions")
	{
		transactions.GET("/", serv.GetAllTransactions)
		transactions.GET("/:id", serv.GetTransaction)
		transactions.POST("/", serv.CreateTransaction)
		transactions.PUT("/:id", serv.UpdateTransaction)
	}

	rabbitmq := serv.Engine.Group("/api/rabbitmq")
	{
		rabbitmq.GET("/stop_consumers", serv.StopConsumers)
	}
}

func (serv *Server) RunMigrations() {
	serv.DB.AutoMigrate(&models.Transaction{})
}

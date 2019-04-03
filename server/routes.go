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

	crmEngine := serv.Engine.Group("/api/crm_engine")
	{
		crmEngine.GET("/ibans", serv.GetAllIbans)
	}
}

func (serv *Server) RunMigrations() {
	serv.DB.AutoMigrate(&models.Transaction{})
}

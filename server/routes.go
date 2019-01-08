package server

import (
	"github.com/golang_tryout/app/models"
)

func (serv *Server) SetRoutes() {
	serv.Engine.GET("/transactions", serv.GetTransactions)
}

func (serv *Server) RunMigrations() {
	serv.DB.AutoMigrate(&models.Transaction{})
}

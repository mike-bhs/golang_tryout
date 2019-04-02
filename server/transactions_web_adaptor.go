package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mike-bhs/golang_tryout/app/models"
)

func (serv *Server) StopConsumers(c *gin.Context) {
	serv.MessagingClient.StopConsumers()

	c.JSON(200, gin.H{
		"status":  http.StatusOK,
		"message": "Successfuly stoped consumers",
	})
}

func (serv *Server) GetAllTransactions(c *gin.Context) {
	data := serv.DB.Find(&models.Transactions{}).Value

	c.JSON(200, gin.H{
		"status":       http.StatusOK,
		"transactions": data,
	})
}

func (serv *Server) CreateTransaction(c *gin.Context) {
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	transaction := models.Transaction{Currency: c.PostForm("currency"), Amount: amount}

	serv.DB.Save(&transaction)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Transaction was succesfully created! transaction_id: %d", transaction.ID),
	})
}

func (serv *Server) GetTransaction(c *gin.Context) {
	var transaction models.Transaction
	transactionID := c.Param("id")

	serv.DB.First(&transaction, transactionID)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": transaction})
}

func (serv *Server) UpdateTransaction(c *gin.Context) {
	var transaction models.Transaction
	transactionID := c.Param("id")

	serv.DB.First(&transaction, transactionID)

	amount, present := c.GetPostForm("amount")

	if present {
		amountFlout, _ := strconv.ParseFloat(amount, 64)
		transaction.Amount = amountFlout
	}

	currency, present := c.GetPostForm("currency")

	if present {
		transaction.Currency = currency
	}

	serv.DB.Save(&transaction)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Transaction was succesfully updated! transaction_id: %d", transaction.ID),
	})
}

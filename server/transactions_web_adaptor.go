package server

import "github.com/gin-gonic/gin"

func (serv *Server) GetTransactions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "tansactions",
	})
}

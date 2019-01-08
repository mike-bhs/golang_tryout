package web

import "github.com/gin-gonic/gin"

func GetTransactions(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "tansactions",
	})
}

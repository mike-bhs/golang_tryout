package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang_tryout/lib/web"
)

func StartServer() *gin.Engine {
	engine := gin.Default()

	setRoutes(engine)

	return engine
}

func setRoutes(engine *gin.Engine) {
	engine.GET("/transactions", web.GetTransactions)
}

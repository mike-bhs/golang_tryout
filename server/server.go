package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mike-bhs/golang_tryout/db"
	"github.com/mike-bhs/golang_tryout/rabbitmq"
)

type Server struct {
	Engine         *gin.Engine
	*rabbitmq.AmqpClient
	*db.DataBase
}

func (serv *Server) DB() *gorm.DB {
	return serv.DataBase.DB
}

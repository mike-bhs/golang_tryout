package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/streadway/amqp"
)

type Server struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Amqp   *amqp.Connection
}

func InitializeServer() *Server {
	db, err := gorm.Open("mysql", "root:@tcp(nix-box:3306)/golang_tryout?charset=utf8&parseTime=True&loc=Local")
	amqp, err := amqp.Dial("amqp://guest:guest@nix-box:5672/")

	if err != nil {
		return nil
	}

	return &Server{
		Engine: gin.Default(),
		DB:     db,
		Amqp:   amqp,
	}
}

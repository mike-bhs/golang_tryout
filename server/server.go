package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Server struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func InitializeServer() *Server {
	db, err := gorm.Open("mysql", "root:@tcp(nix-box:3306)/golang_tryout?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return nil
	}

	return &Server{
		Engine: gin.Default(),
		DB:     db,
	}
}

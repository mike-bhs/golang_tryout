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
	*MessagingClient
}

type MessagingClient struct {
	Connection *amqp.Connection
	Consumers  []*MessageConsumer
}

func (mc *MessagingClient) RegisterConsumer(m *MessageConsumer) {
	mc.Consumers = append(mc.Consumers, m)
}

type MessageConsumer struct {
	Queue        amqp.Queue
	Channel      *amqp.Channel
	ConsumerName string
	HandlerFunc  func(amqp.Delivery)
	CloseControl chan (bool)
}

func (consumer *MessageConsumer) CancelConsumer() {
	go func() {
		consumer.CloseControl <- true
	}()
}

func InitializeServer() *Server {
	db, err := gorm.Open("mysql", "root:@tcp(localhost:3306)/golang_tryout?charset=utf8&parseTime=True&loc=Local")
	amqp, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		return nil
	}

	client := &MessagingClient{Connection: amqp, Consumers: []*MessageConsumer{}}

	return &Server{
		Engine:          gin.Default(),
		DB:              db,
		MessagingClient: client,
	}
}

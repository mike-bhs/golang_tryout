package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
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

func (mc *MessagingClient) StopConsumers() {
	log.Println("Stopping RabbitMQ consumers...")

	for _, consumer := range mc.Consumers {
		consumer.CancelConsumer()
	}
}

func (mc *MessagingClient) WaitForConsumersToStop() {
	shutdownChan := make(chan bool)

	go func() {
		for {
			if !mc.HasBusyConsumers() {
				shutdownChan <- true
				break
			}

			log.Println("Waiting for message processing to finish ...")
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-shutdownChan:
		log.Println("RabbitMQ consumers stopped gracefully")
		os.Exit(0)
	case <-time.After(15 * time.Second):
		log.Println("TIMEOUT")
		os.Exit(1)
	}
}

func (mc *MessagingClient) HasBusyConsumers() bool {
	for _, consumer := range mc.Consumers {
		if consumer.IsBusy {
			return true
		}
	}

	return false
}

type MessageConsumer struct {
	Queue        amqp.Queue
	Channel      *amqp.Channel
	ConsumerName string
	IsBusy 		 bool
	handlerFunc  func(amqp.Delivery)
}

func (mc *MessageConsumer) CancelConsumer() {
	mc.Channel.Cancel(mc.ConsumerName, true)
}

func (mc *MessageConsumer) HandleDelivery(d amqp.Delivery) {
	mc.IsBusy = true
	mc.handlerFunc(d)
	mc.IsBusy = false
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

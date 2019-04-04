package server

import (
	"fmt"
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

const ReconnectTimeoutSec = 5 * time.Second
const ReconnectAttemptsAmount = 3

func (serv *Server) ConnectToDbAsync(ch chan bool, user, password, host, dbName string) {
	dbHostUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)

	go func() {
		for i := 1; i < ReconnectAttemptsAmount; i++ {
			log.Printf("Connection to DB attempt %d of %d", i, ReconnectAttemptsAmount)

			db, err := gorm.Open("mysql", dbHostUrl)

			if err != nil {
				log.Println("Failed to connect to database", err)
				time.Sleep(ReconnectTimeoutSec)
				continue
			}

			log.Println("Db connection success")
			serv.DB = db
			ch <- true

			return
		}

		ch <- false
	}()
}

func (serv *Server) ConnectToRabbitMQAsync(ch chan bool, user, password, host string) {
	amqpHostUrl := fmt.Sprintf("amqp://%s:%s@%s/", user, password, host)

	go func() {
		for i := 1; i <= ReconnectAttemptsAmount; i++ {
			log.Printf("Connection to RabbitMQ attempt %d of %d", i, ReconnectAttemptsAmount)

			conn, err := amqp.Dial(amqpHostUrl)

			if err != nil {
				log.Println("Failed to connect to RabbitMQ host", err)
				time.Sleep(ReconnectTimeoutSec)
				continue
			}

			log.Println("Amqp connection success")
			client := &MessagingClient{Connection: conn, Consumers: []*MessageConsumer{}}
			serv.MessagingClient = client
			ch <- true

			return
		}

		ch <- false
	}()
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
	err := mc.Channel.Cancel(mc.ConsumerName, true)

	if err != nil {
		log.Println("Failed to cancel consumer", mc.ConsumerName)
	}
}

func (mc *MessageConsumer) HandleDelivery(d amqp.Delivery) {
	mc.IsBusy = true
	mc.handlerFunc(d)
	mc.IsBusy = false
}

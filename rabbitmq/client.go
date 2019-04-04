package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const ReconnectTimeoutSec = 5 * time.Second

type AmqpClient struct {
	ConnectionIn  *amqp.Connection
	ConnectionOut *amqp.Connection
	Consumers     Consumers
	*Config
}

type Config struct {
	User     string
	Password string
	Host     string
}

func (c *Config) ToUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", c.User, c.Password, c.Host)
}

func (ac *AmqpClient) EstablishConnection() {
	go ac.createConnectionIn()
	//go ac.createConnectionOut()
}

func (ac *AmqpClient) createConnectionIn() {
	log.Println("Establishing RabbitMQ connection for incoming messages")

	connIn, err := amqp.Dial(ac.Config.ToUrl())

	if err != nil {
		log.Println("Failed to connect to RabbitMQ host", err)
		time.Sleep(ReconnectTimeoutSec)

		ac.createConnectionIn()
	}

	ac.ConnectionIn = connIn
	log.Println("Amqp ConnectionIn success")
	ac.StartQueues()
}

//func (ac *AmqpClient) createConnectionOut() {
//	log.Println("Establishing RabbitMQ connection for outcoming messages")
//
//	connOut, err := amqp.Dial(ac.Config.ToUrl())
//
//	if err != nil {
//		log.Println("Failed to connect to RabbitMQ host", err)
//		time.Sleep(ReconnectTimeoutSec)
//
//		ac.createConnectionOut()
//	}
//
//	ac.ConnectionOut = connOut
//	log.Println("Amqp ConnectionOut success")
//}

func (ac *AmqpClient) StartQueues() {
	conn := ac.ConnectionIn

	if ac.IsConnectionEmpty(conn) {
		log.Println("Connection is empty, waiting for connection to open")
		ac.StartQueues()
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed to open a channel")
		ac.StartQueues()
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang_tryout", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	if err != nil {
		log.Println("Failed to create a queue")
		ac.StartQueues()
	}

	consumer := &Consumer{
		Queue:        q,
		Channel:      ch,
		ConsumerName: "golang_tryout_consumer",
		handlerFunc: func(d amqp.Delivery) {
			log.Println("Handling message START: ", string(d.Body))
			time.Sleep(10 * time.Second)
			log.Println("Handling message DONE: ", string(d.Body))
		},
		IsBusy: false,
	}

	ac.RegisterConsumer(consumer)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)
	consumer.ConsumeMessages()
	<-forever
}

func (ac *AmqpClient) RegisterConsumer(c *Consumer) {
	ac.Consumers = append(ac.Consumers, c)
}

func (ac *AmqpClient) HasBusyConsumers() bool {
	for _, consumer := range ac.Consumers {
		if consumer.IsBusy {
			return true
		}
	}

	return false
}

func (ac *AmqpClient) CloseConnection() {
	err := ac.ConnectionIn.Close()

	if err != nil {
		log.Println("Failed to close RabbitMQ connection for incomming messages", err)
	}

	err = ac.ConnectionOut.Close()

	if err != nil {
		log.Println("Failed to close RabbitMQ connection for outcoming messages", err)
	}

	if err == nil {
		log.Println("Successfully closed RabbitMQ connections")
	}
}

func (ac *AmqpClient) IsConnectionEmpty(conn *amqp.Connection) bool {
	return conn == (&amqp.Connection{})
}

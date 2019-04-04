package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	Queue        amqp.Queue
	Channel      *amqp.Channel
	ConsumerName string
	IsBusy       bool
	handlerFunc  func(amqp.Delivery)
}

type Consumers []*Consumer

func (c *Consumer) CancelConsumer() {
	err := c.Channel.Cancel(c.ConsumerName, true)

	if err != nil {
		log.Println("Failed to cancel consumer", c.ConsumerName)
	}
}

func (c *Consumer) HandleDelivery(d amqp.Delivery) {
	c.IsBusy = true
	c.handlerFunc(d)
	c.IsBusy = false
}

func (c *Consumer) ConsumeMessages() {
	deliveriesChan, err := c.Channel.Consume(
		c.Queue.Name,   // queue
		c.ConsumerName, // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		log.Println("Failed to register an consumer")
		c.ConsumeMessages()
	}

	go func() {
		for delivery := range deliveriesChan {
			log.Println("========= Received a message ========")
			c.HandleDelivery(delivery)
			log.Println("========= Processed message a message ========")
		}
	}()
}

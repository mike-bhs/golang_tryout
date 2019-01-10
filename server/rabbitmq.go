package server

import (
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (serv *Server) StartQueues() {
	conn := serv.Amqp
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang_tryout", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("========= Received a message: %s", d.Body)
			log.Printf("Queue: %s, Messages: %d, Counsumers: %d", q.Name, q.Messages, q.Consumers)
			time.Sleep(20 * time.Second)
			log.Printf("Queue: %s, Messages: %d, Counsumers: %d", q.Name, q.Messages, q.Consumers)
			log.Printf("========= Processed a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

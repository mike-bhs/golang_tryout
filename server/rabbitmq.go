package server

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (serv *Server) StartQueues() {
	conn := serv.MessagingClient.Connection

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

	consumer := &MessageConsumer{
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

	serv.MessagingClient.RegisterConsumer(consumer)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)
	consumeMessages(consumer)
	<-forever
}

func consumeMessages(mc *MessageConsumer) {
	deliveriesChan, err := mc.Channel.Consume(
		mc.Queue.Name,   // queue
		mc.ConsumerName, // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for delivery := range deliveriesChan {
			log.Println("========= Received a message ========")
			mc.HandleDelivery(delivery)
			log.Println("========= Processed message a message ========")
		}
	}()
}

func (serv *Server) InitializeGracefulShutdown() {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Wait for 15 second to finish processing")

		conn := serv.MessagingClient.Connection
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		serv.MessagingClient.StopConsumers()
		serv.MessagingClient.WaitForConsumersToStop()
	}()
}

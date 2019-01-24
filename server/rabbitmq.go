package server

import (
	"log"
	"time"

	"github.com/streadway/amqp"
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
		CloseControl: make(chan bool),
		HandlerFunc: func(d amqp.Delivery) {
			log.Println("Handling message: ", string(d.Body))
		},
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
			select {
			case <-mc.CloseControl:
				mc.Channel.Cancel(mc.ConsumerName, true)
				log.Println("CANCELLED CONSUMER")
			default:
				log.Println("========= Received a message ========")
				mc.HandlerFunc(delivery)
				time.Sleep(10 * time.Second)
				log.Println("========= Processed message a message ========")
			}
		}
	}()
}

// func (serv *Server) InitializeGracefulShutdown() {
// 	gracefulStop := make(chan os.Signal)

// 	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

// 	go func() {
// 		sig := <-gracefulStop
// 		log.Printf("caught sig: %+v", sig)
// 		log.Println("Wait for 15 second to finish processing")

// 		conn := serv.Amqp
// 		defer conn.Close()

// 		ch, err := conn.Channel()
// 		failOnError(err, "Failed to open a channel")
// 		defer ch.Close()

// 		// q, err := ch.QueueInspect("golang_tryout")
// 		// failOnError(err, "Failed to inspect a queue")

// 		err = ch.Cancel("golang_tryout_consumer", true)
// 		failOnError(err, "Failed to inspect a queue")

// 		select {
// 		case <-time.After(15 * time.Second):
// 			log.Println("TIMEOUT")
// 			os.Exit(1)
// 		}
// 	}()
// }

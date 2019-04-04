package rabbitmq

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (ac *AmqpClient) InitializeGracefulShutdown() {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Wait for 15 second to finish processing")

		ac.StopConsumers()
		ac.WaitForConsumersToStop()
	}()
}

func (ac *AmqpClient) WaitForConsumersToStop() {
	shutdownChan := make(chan bool)

	go func() {
		for {
			if !ac.HasBusyConsumers() {
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

func (ac *AmqpClient) StopConsumers() {
	log.Println("Stopping RabbitMQ consumers...")

	for _, consumer := range ac.Consumers {
		consumer.CancelConsumer()
	}
}

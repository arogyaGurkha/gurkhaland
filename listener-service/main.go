package main

import (
	"fmt"
	"github.com/arogyaGurkha/gurkhaland/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// connect to RabbitMQ
	rabbitConn, err := connectToAMQP()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// listen for messages
	log.Println("Listening and consuming RabbitMQ messages..")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch queue, consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connectToAMQP() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// wait for rabbit
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ!")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// raise to power 2
		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}

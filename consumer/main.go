package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	rabbitURL := os.Getenv("QPC_RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"my-data",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	quit := make(chan bool)

	go func() {
		for m := range msgs {
			js, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}

			// Instead of just printing, we could write to a file and send to s3 here
			// Or stream it to another queue, such as kinesis, sqs, kafka, etc
			fmt.Printf("%s\n", js)
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")

	<-quit
}

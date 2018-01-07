package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func main() {
	rabbitURL := os.Getenv("QPC_RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	r := mux.NewRouter()
	r.Handle("/send", sendHandler(rabbitURL, "testing", "my-data")).Methods("POST")

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}
	log.Fatal(s.ListenAndServe())
}

// Message -
type Message struct {
	Name string `json:"name"`
}

func sendHandler(rabbitURL, s, qn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("sendHandler called")
		w.Header().Set("Content-Type", "application/json")

		var m Message
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &m)

		err := sendMessage(rabbitURL, m.Name, qn)
		if err != nil {
			fmt.Printf("error: %+v\n", err)
			w.Write([]byte("err"))
		}

		fmt.Println("sent")
		w.Write([]byte("Message Send"))
	}
}

// This could be in a package called producer ?
func sendMessage(rabbitURL, s, qn string) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qn,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(s),
		},
	)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

// TODO - MAYBE

// Producer - maybe
type Producer interface {
	SendMessage(s, qn string, ch *amqp.Channel)
}

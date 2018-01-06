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
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	r := mux.NewRouter()
	r.Handle("/send", sendHandler("testing", "my-data", ch)).Methods("POST")

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

func sendHandler(s, qn string, ch *amqp.Channel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var m Message
		b, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(b, &m)

		err := sendMessage(m.Name, qn, ch)
		if err != nil {
			w.Write([]byte("Error"))
		}
		w.Write([]byte("Message Send"))
	}
}

// This could be in a package called producer ?
func sendMessage(s, qn string, ch *amqp.Channel) error {
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

	fmt.Println(s)

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

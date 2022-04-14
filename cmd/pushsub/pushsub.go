package main

import (
	"log"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/jetstream"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Print("Push Sub!")

	// Simple Sync Subscriber
	js := jetstream.Connect()

	c := make(chan string)

	go func() {
		js.Subscribe("ORDERS.async", func(m *nats.Msg) {
			c <- string(m.Data)
		})
	}()

	counter := 0
	for msg := range c {
		counter++
		log.Printf("%d- MSG: %s", counter, msg)
	}
}

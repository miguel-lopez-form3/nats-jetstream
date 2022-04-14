package main

import (
	"log"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/jetstream"
)

func main() {
	log.Print("Pull Sub!")
	js := jetstream.Connect()

	// Simple Pull Consumer
	sub, err := js.PullSubscribe("ORDERS.*", "MONITOR")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := sub.Fetch(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range msgs {
		log.Printf("Received a JetStream message: %s\n", string(m.Data))
	}

	// Unsubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()
}

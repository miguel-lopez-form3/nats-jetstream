package main

import (
	"fmt"
	"log"
	"time"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/jetstream"
)

func main() {
	log.Print("AsyncPub!")
	js := jetstream.Connect()

	// Simple Async Stream Publisher
	for i := 0; i < 500; i++ {
		js.PublishAsync("ORDERS.async", []byte("async hello"))
	}
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}

	log.Print("Published async successully")
}

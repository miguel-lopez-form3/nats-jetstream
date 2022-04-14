package main

import (
	"log"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/jetstream"
)

func main() {
	log.Print("SyncPub!")
	js := jetstream.Connect()

	js.Publish("ORDERS.sync", []byte("sync hello"))

	log.Print("Published sync successully")
}

package main

import (
	"log"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/certs"
)

func main() {
	err := certs.CreateCertKeyPairOnFileSystem("server")
	if err != nil {
		log.Fatalf("failed to create the server certs: %v", err)
	}
	err = certs.CreateCertKeyPairOnFileSystem("client")
	if err != nil {
		log.Fatalf("failed to create the client certs: %v", err)
	}
	log.Println("created the certs successfully!")
}

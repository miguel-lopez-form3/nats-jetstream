package main

import (
	"log"

	"github.com/miguel-lopez-form3/nats-jetstream/internal/certs"
)

func main() {
	parentCA, err := certs.CreateCAToFileSystem()
	if err != nil {
		log.Fatalf("failed to create the CA: %v", err)
	}
	err = certs.CreateCertKeyPairOnFileSystem("server", parentCA)
	if err != nil {
		log.Fatalf("failed to create the server certs: %v", err)
	}
	err = certs.CreateCertKeyPairOnFileSystem("client", parentCA)
	if err != nil {
		log.Fatalf("failed to create the client certs: %v", err)
	}
	log.Println("created the certs successfully!")
}

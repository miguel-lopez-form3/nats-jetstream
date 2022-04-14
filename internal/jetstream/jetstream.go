package jetstream

import (
	"crypto/tls"
	"log"

	"github.com/nats-io/nats.go"
)

func Connect() nats.JetStreamContext {
	certFile := "./certs/client/client.crt"
	keyFile := "./certs/client/client.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("error parsing X509 certificate/key pair: %v", err)
	}

	config := &tls.Config{
		ServerName:   "nats",
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	nc, err := nats.Connect("tls://localhost:4222", nats.Secure(config))
	if err != nil {
		log.Fatalf("failed to connect to the NATS server: %v", err)
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatalf("failed to initialise JetStream: %v", err)
	}
	setup(js)
	return js
}

func setup(js nats.JetStreamContext) {
	for name := range js.StreamNames() {
		if name == "ORDERS" {
			return
		}
	}
	// Create a Stream
	js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
}

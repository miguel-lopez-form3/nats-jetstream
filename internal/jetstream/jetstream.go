package jetstream

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Connect() nats.JetStreamContext {
	/*
		certFile := "./certs/f-client.pem"
		keyFile := "./certs/f-client-key.pem"
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("error parsing X509 certificate/key pair: %v", err)
		}

		rootPEM, err := ioutil.ReadFile("./certs/rootCA.pem")
		if err != nil || rootPEM == nil {
			log.Fatalf("failed to read root certificate")
		}

		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM([]byte(rootPEM))
		if !ok {
			log.Fatal("failed to parse root certificate")
		}

		config := &tls.Config{
			ServerName:   "localhost",
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
			MinVersion:   tls.VersionTLS12,
		}
	*/

	nc, err := nats.Connect("nats://localhost:4222")
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

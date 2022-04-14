package main

import "github.com/miguel-lopez-form3/nats-jetstream/internal/certs"

func main() {
	certs.CreateCertKeyPairOnFileSystem("system")
	certs.CreateCertKeyPairOnFileSystem("client")
}

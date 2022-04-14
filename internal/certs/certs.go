package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func CreateCertKeyPairOnFileSystem(organisationId string) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	cert, err := generateSelfSignedCertificateAsBytes(key, "Form3")
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("certs/%s", organisationId)

	certFilename := fmt.Sprintf("%s/%s.crt", dir, organisationId)
	certOut, err := os.Create(certFilename)
	if err != nil {
		return err
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyFilename := fmt.Sprintf("%s/%s.key", dir, organisationId)
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		return err
	}
	err = keyOut.Close()

	return err
}

func generateSelfSignedCertificateAsBytes(key *rsa.PrivateKey, orgName string) ([]byte, error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{orgName},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"London"},
			StreetAddress: []string{"7 Harp Ln"},
			PostalCode:    []string{"EC3R 6DP"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv4(0, 0, 0, 0), net.IPv6loopback},
		NotBefore:    time.Now().Add(time.Minute * -5),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"nats", "localhost"},
	}

	return x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
}

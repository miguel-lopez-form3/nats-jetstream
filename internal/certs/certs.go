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

type Cert struct {
	Certificate      *x509.Certificate
	PrivKey          *rsa.PrivateKey
	CertificateBytes []byte
}

func CreateCAToFileSystem() (*Cert, error) {
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Form3"},
			Country:       []string{"UK"},
			Province:      []string{""},
			Locality:      []string{"London"},
			StreetAddress: []string{"7 Harp Ln"},
			PostalCode:    []string{"EC3R 6DP"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certOut, err := os.Create("certs/ca/rootCA.pem")
	if err != nil {
		return nil, err
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	if err != nil {
		return nil, err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(caPrivKey)
	keyOut, err := os.OpenFile("certs/ca/rootCA-key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0400)
	if err != nil {
		return nil, err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		return nil, err
	}
	err = keyOut.Close()
	if err != nil {
		return nil, err
	}

	return &Cert{
		Certificate:      ca,
		PrivKey:          caPrivKey,
		CertificateBytes: caBytes,
	}, nil
}

func CreateCertKeyPairOnFileSystem(organisationId string, parentCert *Cert) error {
	cert, err := generateSelfSignedCertificateAsBytes("Form3", parentCert)
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("certs/%s", organisationId)

	certFilename := fmt.Sprintf("%s/%s.crt", dir, organisationId)
	certOut, err := os.Create(certFilename)
	if err != nil {
		return err
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert.CertificateBytes})
	if err != nil {
		return err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(cert.PrivKey)
	keyFilename := fmt.Sprintf("%s/%s.key", dir, organisationId)
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		return err
	}
	err = keyOut.Close()

	return err
}

func generateSelfSignedCertificateAsBytes(orgName string, parentCert *Cert) (*Cert, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

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
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now().Add(time.Minute * -5),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		DNSNames:     []string{"nats", "localhost"},
	}

	var parent *x509.Certificate
	if parentCert == nil {
		parent = template
	} else {
		parent = parentCert.Certificate
	}

	var priv *rsa.PrivateKey
	if parentCert == nil {
		priv = key
	} else {
		priv = parentCert.PrivKey
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, &key.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	return &Cert{
		Certificate:      template,
		PrivKey:          key,
		CertificateBytes: certBytes,
	}, nil
}

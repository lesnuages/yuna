package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
)

func getExtKeyUsage(objectType string) x509.ExtKeyUsage {
	switch objectType {
	case "client":
		return x509.ExtKeyUsageClientAuth
	case "server":
		return x509.ExtKeyUsageServerAuth
	default:
		return x509.ExtKeyUsageServerAuth
	}
}

func GenerateCertificate(objType string, data []string) (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 365 * 24)
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: data, // hide data in Org name
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{getExtKeyUsage(objType)},
		BasicConstraintsValid: true,
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("Marshal failed: %s", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	return tls.X509KeyPair(pemCert, pemKey)
}

package common

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var HttpClient http.Client

var tlsConfig *tls.Config

func init() {
	rootCAs := x509.NewCertPool()
	caCertPath := os.Getenv("VARUNA_LOCAL_CA")

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatal("Failed to read root CA certificate:", err)
	}

	if !rootCAs.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed to append root CA certificate")
	}

	tlsConfig = &tls.Config{
		RootCAs: rootCAs,
	}

	tr := &http3.Transport{
		TLSClientConfig: tlsConfig,
		QUICConfig:      &quic.Config{},
	}
	HttpClient = http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
}

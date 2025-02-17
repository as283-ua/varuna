package common

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"time"
)

var HttpClient http.Client

var tlsConfig *tls.Config

func init() {
	rootCAs := x509.NewCertPool()
	caCertPath := "/home/andrejs/snap/code/181/.local/share/mkcert/rootCA.pem"

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

	HttpClient = http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
	}
}

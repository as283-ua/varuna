package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"server/handler"

	"github.com/quic-go/quic-go/http3"
)

const (
	addr = "0.0.0.0:4433"
	cert = "../varuna.cauth.crt"
	key  = "../varuna.cauth.key"
)

func main() {
	router := http.NewServeMux()

	router.Handle("GET /hello", http.HandlerFunc(handler.HandleHello))

	fmt.Println("Server started on https://" + addr)

	if useHttp3 := true; useHttp3 {
		server := &http3.Server{
			Addr:      addr,
			Handler:   router,
			TLSConfig: TLSConfig(),
		}

		if err := server.ListenAndServeTLS(cert, key); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServeTLS(addr, cert, key, router); err != nil {
			panic(err)
		}
	}

}

func TLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}
}

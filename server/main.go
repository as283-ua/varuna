package main

import (
	"fmt"
	"net/http"
	"server/handler"
)

func main() {
	router := http.NewServeMux()

	router.Handle("GET /hello", http.HandlerFunc(handler.HandleHello))

	fmt.Println("Server started on https://127.0.0.1/")
	if err := http.ListenAndServeTLS("127.0.0.1:443", "127.0.0.1.pem", "127.0.0.1-key.pem", router); err != nil {
		panic(err)
	}
}

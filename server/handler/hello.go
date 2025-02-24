package handler

import (
	"fmt"
	"net/http"
)

func HandleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Connected %v\n", r.RemoteAddr)
	w.Write([]byte("hello"))
}

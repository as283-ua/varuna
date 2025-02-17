package main

import (
	"client/common"
	"fmt"
	"log"
)

func main() {
	client := common.HttpClient

	resp, err := client.Get("https://127.0.0.1/hello")

	if err != nil {
		log.Fatal(err.Error())
	}

	bytes := make([]byte, 1024)

	resp.Body.Read(bytes)
	fmt.Println(string(bytes))
}

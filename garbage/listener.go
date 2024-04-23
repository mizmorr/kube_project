package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Kubia server starting...")
	handler := func(response http.ResponseWriter, request *http.Request) {
		fmt.Println("Received request from", request.RemoteAddr)
		response.WriteHeader(http.StatusOK)
		name, _ := os.Hostname()
		response.Write([]byte("You've hit " + name + "\n"))
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

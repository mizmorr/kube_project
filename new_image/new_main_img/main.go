package main

import (
	"fmt"
	"net/http"
)

func logRequest(req *http.Request) {
	uri := req.URL.String()
	method := req.Method
	header := req.Header
	fmt.Println(uri, "\n", method, "\n", header, "\n---")
}

func main() {

	http.HandleFunc("/get_service", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.Host)
		logRequest(r)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
		logRequest(r)
	})

	fmt.Println("Server is listening port...")
	http.ListenAndServe(":9080", nil)
}

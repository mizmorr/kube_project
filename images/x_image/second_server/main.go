package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "server second")
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
	})
	fmt.Println("Server second is listening...")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/productpage", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "About Page")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index Page")
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}

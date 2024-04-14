package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func logRequest(req *http.Request) {
	uri := req.URL.String()
	method := req.Method
	header := req.Header
	fmt.Println(uri, "\n", method, "\n", header, "\n---")
}
func main() {

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "about.html")
		logRequest(r)

	})
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index Page")
		logRequest(r)
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
		logRequest(r)
	})
	http.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		port := os.Getenv("SERVER2_SERVICE_PORT")
		host := os.Getenv("SERVER2_SERVICE_HOST")
		resp, err := http.Get("http://" + host + ":" + port + "/graph")
		if err != nil {
			fmt.Println("some error occurs", err)
			fmt.Fprint(w, "some error occurs: ", err)

		} else {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Fprintln(w, string(b))
		}

		logRequest(r)

	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func logRequest(req *http.Request) {
	uri := req.URL.String()
	method := req.Method
	header := req.Header
	fmt.Println(uri, "\n", method, "\n", header, "\n---")
}
func make_slave_request(slave_num, port string, client *http.Client, w http.ResponseWriter) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://slave%s:%s/hi", slave_num, port), nil)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s\n", bodyText)

}
func main() {
	http.HandleFunc("/hi_to_slave/", func(w http.ResponseWriter, r *http.Request) {
		uri := strings.SplitAfter(r.URL.Path, "/hi_to_slave/")
		slavery_num := uri[len(uri)-1]
		port := os.Getenv(fmt.Sprintf("SLAVE%s_SERVICE_PORT", slavery_num))
		if port == "" {
			fmt.Fprint(w, "no such slave")
		} else {
			client := &http.Client{}

			make_slave_request(slavery_num, port, client, w)

			fmt.Println("slave request")
		}

	})
	http.HandleFunc("/hi_every1", func(w http.ResponseWriter, r *http.Request) {
		slave_num := 1
		for (os.Getenv(fmt.Sprintf("SLAVE%d_SERVICE_PORT", slave_num))) != "" {
			slave_num++
		}
		client := &http.Client{}

		for i := 1; i < slave_num; i++ {

			port := os.Getenv(fmt.Sprintf("SLAVE%s_SERVICE_PORT", fmt.Sprint((i))))
			make_slave_request(fmt.Sprint(i), port, client, w)
			fmt.Printf("curl to slave%d\n", i)

		}
		fmt.Println("hello request done")

	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "about.html")
		logRequest(r)

	})
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Index Page of main server")
		logRequest(r)
	})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
		logRequest(r)
	})
	http.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
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
	listen_port := os.Getenv("SERVER_SERVICE_PORT")
	fmt.Println("Server is listening port...", listen_port)
	http.ListenAndServe(":"+listen_port, nil)
}

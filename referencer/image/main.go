package main

import (
	"encoding/json"
	"fmt"
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

type Req_body struct {
	Name string `json:"name"`
}

func make_map() map[string]string {
	res := make(map[string]string)
	k, err := os.ReadFile("markup.txt")
	if err != nil {
		panic(err)
	}
	arr := strings.Split(strings.Trim(string(k), "\n"), "\n")
	for _, s := range arr {
		curr := strings.Split(s, ":")
		res[curr[0]] = curr[1]
	}
	return res
}
func main() {
	reference_map := make_map()
	http.HandleFunc("/get_service", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req Req_body
		err := decoder.Decode(&req)
		if err != nil {
			panic(err)
		}
		name := req.Name
		cluster_num := reference_map[name]
		fmt.Fprintln(w, cluster_num)
		logRequest(r)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
		logRequest(r)
	})

	fmt.Println("Server is listening port...")
	http.ListenAndServe(":9080", nil)
}

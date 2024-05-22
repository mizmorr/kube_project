package main

import (
	"encoding/json"
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
	log.Println(uri, "\n", method, "\n", header, "\n---")
}

type Req_body struct {
	Name string `json:"name"`
}

func get_slaves(name string) (res []string, err error) {
	json_file, err := os.Open("slaves.json")
	if err != nil {
		return nil, err
	}
	var k map[string]interface{}
	defer json_file.Close()
	byte_value, _ := io.ReadAll(json_file)
	json.Unmarshal(byte_value, &k)
	result := k[name]
	if result != nil {
		m_ap := result.(map[string]interface{})
		for k := range m_ap {
			res = append(res, k)
		}
		return
	} else {
		return nil, fmt.Errorf("no slaves found")
	}

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
func slaves_req_handler(r *http.Request, res http.ResponseWriter, ref map[string]string) {
	decoder := json.NewDecoder(r.Body)
	var req Req_body
	err := decoder.Decode(&req)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	name := req.Name
	cluster_num := ref[name]
	slaves, err := get_slaves(cluster_num)
	if err != nil {
		log.Fatal(err.Error())
		fmt.Fprint(res, "")
	} else {
		logRequest(r)
		fmt.Fprint(res, slaves)
	}

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

	http.HandleFunc("/get_slaves", func(w http.ResponseWriter, r *http.Request) {
		slaves_req_handler(r, w, reference_map)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "metrics is good")
		logRequest(r)
	})
	fmt.Println("Server is listening port...")
	http.ListenAndServe(":9080", nil)
}

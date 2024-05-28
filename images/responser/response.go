package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Json struct {
	Cluster string   `json:"cluster"`
	Slaves  []string `json:"slaves"`
}
type Req_body struct {
	Name string `json:"name"`
}

func make_response() (string, string) {
	read, err := os.ReadFile("/etc/hosts")

	if err != nil {
		log.Fatal(err.Error())
	}
	str := strings.Split(strings.TrimSpace(regexp.MustCompile(`["\n", "\t"]`).ReplaceAllString(string(read), " ")), " ")
	num := strings.Split(str[len(str)-1], "-")[0]
	whoami := str[len(str)-1] + "! My IP is - " + str[len(str)-2]
	return fmt.Sprintf("Hello from %s.\n", whoami), num
}

func unpack_response(resp *http.Response) (string, []string) {
	decoder := json.NewDecoder(resp.Body)
	var target Json
	err := decoder.Decode(&target)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	return target.Cluster, target.Slaves
}

func post_req_maker(uri string, req_str interface{}) *http.Request {
	req_body := req_str
	marshalled, err := json.Marshal(req_body)
	if err != nil {
		log.Fatalf("impossible to marshall: %s", err)
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func do_post(num, resource string) (string, []string) {
	port := os.Getenv("REFERENCER_SERVICE_PORT")
	host := os.Getenv("REFERENCER_SERVICE_HOST")
	uri := "http://" + host + ":" + port + "/" + resource
	req_body := Req_body{Name: num}
	req := post_req_maker(uri, req_body)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)

	defer res.Body.Close()

	return unpack_response(res)

}

func slave_request(slaves []string, w http.ResponseWriter) {
	for _, v := range slaves {
		cluster_uri := "http://cluster-" + v + ":9080/hi"
		req, _ := http.NewRequest("GET", cluster_uri, nil)
		client := &http.Client{}
		for i := 0; i < 5; i++ {
			res, err := client.Do(req)
			if err != nil {
				log.Fatalf("impossible to send request: %s", err)
			}
			resBody, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("status Code: %d", res.StatusCode)
			fmt.Fprint(w, "we've that slave - "+string(resBody))
		}
	}
}

func main() {
	response, my_num := make_response()
	cluster_num, slaves := do_post(my_num, "get_info")
	log.Printf("I am %s, my slaves are: %v\n", cluster_num, slaves)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "go hit mates\n")
		log.Print(r.Method, ", from ", r.URL.Path, "\n")
		cluster_uri := "http://cluster-" + cluster_num + ":9080/hi_mates"
		req, err := http.NewRequest("GET", cluster_uri, nil)
		if err != nil {
			log.Fatalf("impossible to build request: %s", err)
		}
		client := &http.Client{}
		for i := 0; i < 5; i++ {
			res, err := client.Do(req)
			if err != nil {
				log.Fatalf("impossible to send request: %s", err)
			}
			resBody, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("status Code: %d", res.StatusCode)
			fmt.Fprint(w, "we've hit him - "+string(resBody))
		}
		if slaves != nil {
			slave_request(slaves, w)
		}
	})
	http.HandleFunc("/hi_mates", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, response)
		log.Print(r.Method, ", from ", r.URL.Path, "\n")
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":9080", nil)

}

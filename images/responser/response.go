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
func do_post(num, resource string) string {
	port := os.Getenv("REFERENCER_SERVICE_PORT")
	host := os.Getenv("REFERENCER_SERVICE_HOST")
	uri := "http://" + host + ":" + port + "/" + resource
	type Req_body struct {
		Name string `json:"name"`
	}
	req_body := Req_body{Name: num}
	marshalled, err := json.Marshal(req_body)
	if err != nil {
		log.Fatalf("impossible to marshall: %s", err)
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	log.Printf("res body: %s", string(resBody))
	return string(resBody)
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
	post_resp := do_post(my_num, "get_slaves")
	is_slaves_exist := post_resp == ""
	cluster_num := strings.TrimSpace(do_post(my_num, "get_service"))
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
		if is_slaves_exist {
			slaves := strings.Split(strings.TrimSpace(post_resp), " ")
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

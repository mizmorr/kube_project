package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

func make_request(name string, client *http.Client, w http.ResponseWriter) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://cluster-%s:9080/hi", name), nil)
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

	names := map[int64]string{0: "zero", 1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth", 6: "sixth", 7: "seventh", 8: "eighth", 9: "nineth"}

	var wg sync.WaitGroup
	http.HandleFunc("/hi_every_one/", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		wg.Add(len(names))
		for _, name := range names {
			go func(n string) {
				for i := 0; i < 4; i++ {
					make_request(n, client, w)
				}
				print(n)
				wg.Done()
			}(name)
		}
		wg.Wait()
	})
	http.ListenAndServe(":9080", nil)

}

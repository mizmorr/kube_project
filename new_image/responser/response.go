package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func make_response() string {
	read, err := os.ReadFile("/etc/hosts")

	if err != nil {
		log.Fatal(err.Error())
	}
	str := strings.Split(strings.TrimSpace(regexp.MustCompile(`["\n", "\t"]`).ReplaceAllString(string(read), " ")), " ")
	whoami := str[len(str)-1] + "! My IP is - " + str[len(str)-2]
	return fmt.Sprintf("Hello from %s.\n", whoami)
}
func main() {
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, make_response())
		fmt.Print(r.Method, ", from ", r.URL.Path, "\n")
	})
	fmt.Println("Server is listening...")
	http.ListenAndServe(":9080", nil)

}

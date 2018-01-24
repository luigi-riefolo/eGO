package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func main() {
	http.HandleFunc("/", hello)

	ln, err := net.Listen("tcp", ":0")

	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	}

	go http.Serve(ln, nil)

	url := fmt.Sprintf("http://0.0.0.0:%d", ln.Addr().(*net.TCPAddr).Port)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error in GET: %s", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	ln.Close()
}

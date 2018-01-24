package main

import (
	"fmt"
	"log"
	"net/http"
)

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	fmt.Fprintln(w, h)

	// [gzip, deflate]
	_ = r.Header["Accept-Encoding"]
	// gzip, deflate
	_ = r.Header.Get("Accept-Encoding")
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/headers", headers)
	log.Fatal(server.ListenAndServe())
}

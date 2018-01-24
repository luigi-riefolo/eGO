package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"golang.org/x/net/http2"
)

// MyHandler ...
type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
func main() {
	handler := MyHandler{}
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: &handler,
	}
	err := http2.ConfigureServer(&server, &http2.Server{})
	if err != nil {
		log.Fatal(err)
	}

	/*
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	*/
	//path := filepath.Dir(os.Args[0])
	_, filename, _, _ := runtime.Caller(0)
	serverCert := filepath.Clean(fmt.Sprintf("%s/../../certs/server.pem", filename))
	serverKey := filepath.Clean(fmt.Sprintf("%s/../../certs/server.key", filename))

	log.Fatal(server.ListenAndServeTLS(serverCert, serverKey))
}

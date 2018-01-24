package main

import (
	"io"
	"log"
	"net/http"
)

func maxClients(h http.Handler, n int) http.Handler {
	sema := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sema <- struct{}{}

		// Free the slot
		defer func() { <-sema }()

		h.ServeHTTP(w, r)
	})
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//res := getExpensiveResource()
		_, err := io.WriteString(w, "OK!!!")
		if err != nil {
			log.Fatalln(err)
		}
	})

	http.Handle("/", maxClients(handler, 5))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//http.Error(w, "", http.StatusTooManyRequests)
//ipAddr := strings.Split(r.RemoteAddr, ":")[0]

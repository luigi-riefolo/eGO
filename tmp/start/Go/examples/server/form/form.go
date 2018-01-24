package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func processFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploaded")
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func processFileTwo(w http.ResponseWriter, r *http.Request) {
	log.Fatal(r.ParseMultipartForm(1024))
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func process(w http.ResponseWriter, r *http.Request) {
	log.Fatal(r.ParseForm())
	fmt.Fprintln(w, r.Form)

	fmt.Fprintln(w, r.PostForm)

	log.Fatal(r.ParseMultipartForm(1024))
	fmt.Fprintln(w, r.MultipartForm)

	fmt.Fprintln(w, r.FormValue("hello"))
	fmt.Fprintln(w, r.Form)

	fmt.Fprintln(w, r.PostFormValue("hello"))
	fmt.Fprintln(w, r.PostForm)

}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/process", process)
	log.Fatal(server.ListenAndServe())
}

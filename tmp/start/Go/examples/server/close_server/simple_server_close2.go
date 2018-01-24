package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// ListenAndServeWithClose ...
func ListenAndServeWithClose(addr string, handler http.Handler) (sc io.Closer, err error) {

	var listener net.Listener

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	return listener, nil
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func main() {

	redirHandler := http.RedirectHandler("http://www.example.com", http.StatusTemporaryRedirect)

	srvCLoser, err := ListenAndServeWithClose(":8080", redirHandler)
	if err != nil {
		log.Fatalln("ListenAndServeWithClose Error - ", err)
	}

	// Do Stuff

	// Close HTTP Server
	err = srvCLoser.Close()
	if err != nil {
		log.Fatalln("Server Close Error - ", err)
	}

	log.Println("Server Closed")
}

package main

import (
	"log"
	"sync"

	// Project packages
	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/server"
)

var (
	conf    config.Config
	service config.Service
)

func init() {

	var err error

	conf, err = config.GetConfig()
	if err != nil {
		log.Fatal("Config file could not be parsed")
	}

	// get service configuration
}

func main() {

	srv := server.NewServer(service)

	// register the gRPC service server

	var wg sync.WaitGroup
	srv.Listen()

	wg.Add(1)

	// wait until the server shuts down
	wg.Wait()
}

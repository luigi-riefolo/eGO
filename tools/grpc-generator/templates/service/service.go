package main

import (
	"log"

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

	srv := server.NewServer(service.Server.Port)

	log.Printf("Starting the %s service: %s:%d",
		service.Name, service.Server.Host, service.Server.Port)

	// register the gRPC service server

	log.Fatal(srv.Listen())
}

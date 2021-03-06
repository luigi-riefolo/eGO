package server

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartPrometheus ...
func (srv *Server) StartPrometheus() {
	log.Println("Initializing Prometheus")

	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(srv.gRPC)

	// Register Prometheus metrics handler.
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		// TODO: use error log
		log.Println(http.ListenAndServe(":8080", nil))
	}()
}

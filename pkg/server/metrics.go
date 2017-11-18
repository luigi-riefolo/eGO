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

/*
func metricTest() {
	promhttp
// Try it once more, this time with a help string.
pushCounter = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "repository_pushes",
    Help: "Number of pushes to external repository.",
})
err = prometheus.Register(pushCounter)
if err != nil {
    fmt.Println("Push counter couldn't be registered AGAIN, no counting will happen:", err)
    return
}
}
*/

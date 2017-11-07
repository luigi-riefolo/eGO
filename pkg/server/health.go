package server

import (
	"context"
	"log"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	healthCheckInterval = 10 * time.Second
)

/*
HealthCheckResponse_UNKNOWN     HealthCheckResponse_ServingStatus = 0
    HealthCheckResponse_SERVING     HealthCheckResponse_ServingStatus = 1
    HealthCheckResponse_NOT_SERVING HealthCheckResponse_ServingStatus = 2
*/
func (srv *Server) setServicesHealth() {

	log.Printf("Services health monitor started")

	for service := range srv.gRPC.GetServiceInfo() {
		srv.health.SetServingStatus(service, healthpb.HealthCheckResponse_SERVING)
		log.Printf("Service %s is serving\n", service)
	}

	srv.startHealthMonitor()
}

func (srv *Server) startHealthMonitor() {
	ticker := time.NewTicker(healthCheckInterval)
	srv.stopHealthMonitor = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				srv.servicesHealth()

			case <-srv.stopHealthMonitor:
				ticker.Stop()
				return
			}
		}
	}()
}

func (srv *Server) servicesHealth() {
	for service := range srv.gRPC.GetServiceInfo() {
		req := &healthpb.HealthCheckRequest{
			Service: service,
		}
		resp, err := srv.health.Check(context.TODO(), req)
		if err != nil {
			log.Printf("Health monitor failed to check status for service %s: %v\n",
				service, err)
		}

		if resp.Status != healthpb.HealthCheckResponse_SERVING {
			status := healthpb.HealthCheckResponse_ServingStatus_name[int32(resp.GetStatus())]
			log.Printf("Service %s is not serving, state: %#v\n",
				service, status)
		}
	}
}

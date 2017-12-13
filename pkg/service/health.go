package service

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
func (svc *Service) setServicesHealth() {

	for service := range svc.gRPC.GetServiceInfo() {
		svc.health.SetServingStatus(service, healthpb.HealthCheckResponse_SERVING)
		// TODO: use debug log
		//log.Printf("Service health info %s is serving\n", service)
	}

	svc.startHealthMonitor()
	log.Printf("%s server health monitor started", svc.name)
}

func (svc *Service) startHealthMonitor() {
	ticker := time.NewTicker(healthCheckInterval)
	svc.stopHealthMonitor = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				svc.servicesHealth()

			case <-svc.stopHealthMonitor:
				ticker.Stop()
				return
			}
		}
	}()
}

func (svc *Service) servicesHealth() {
	for service := range svc.gRPC.GetServiceInfo() {
		req := &healthpb.HealthCheckRequest{
			Service: service,
		}
		resp, err := svc.health.Check(context.TODO(), req)
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

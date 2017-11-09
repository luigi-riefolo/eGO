package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DO WE NEED AN INTERFACE?????
type gateway interface {
	LoadEndPoint(ctx context.Context, path string, registerFunction registerFromEndpoint)
}

// registerFromEndpoint ...
type registerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) (err error)

// Gateway represents an implementation of the GatewayServer interface.
type Gateway struct {
	Creds    credentials.TransportCredentials
	DialOpts []grpc.DialOption
	Host     string
	Port     int
	//mux        *http.ServeMux
	ListenAddr string
	Mux        *runtime.ServeMux
	Services   map[string]*server.Server
}

// LoadEndpoint ...
func (gw *Gateway) LoadEndpoint(ctx context.Context, serviceConf config.Service, regFn registerFromEndpoint) {

	addr := fmt.Sprintf(":%d", serviceConf.Server.Port)

	if err := regFn(ctx, gw.Mux, addr, gw.DialOpts); err != nil {
		log.Fatalf("failed to register endpoint: %v", err)
	}
}

/*
func (gw *gateway) initStaticContentHandlers() {
	// Static content handlers
	swaggerHandler := http.FileServer(http.Dir("swagger-ui"))
	gw.mux.Handle("/help/", http.StripPrefix("/help/", swaggerHandler))
}
*/

// ListenAndServe ...
func (gw *Gateway) ListenAndServe() error {

	gw.HandleSig()

	log.Printf("Gateway listening on %s\n", gw.ListenAddr)

	return http.ListenAndServe(gw.ListenAddr, gw.Mux)
}

// HandleSig ...
func (gw *Gateway) HandleSig() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	stopServices := func(sigStr string) {
		for _, srv := range gw.Services {
			for s := range srv.GetgRPCServer().GetServiceInfo() {
				log.Printf("Shutting down %s\n", s)
			}
			srv.GetgRPCServer().Stop()
		}
	}
	done := make(chan struct{})
	go func() {
		for {
			s := <-sig
			log.Printf("Handling %s signal\n", s.String())

			switch s {

			case syscall.SIGINT:
				stopServices(s.String())
				done <- struct{}{}

			case syscall.SIGTERM:
				stopServices(s.String())
				done <- struct{}{}

			case syscall.SIGQUIT:
				stopServices(s.String())
				done <- struct{}{}

			default:
				log.Printf("Received unknown signal: %s\n", s.String())
			}
		}
	}()

	go func() {
		<-done
		log.Println("Shutting down gateway")
		os.Exit(1)
	}()
}

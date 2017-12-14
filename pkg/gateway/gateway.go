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
	"github.com/luigi-riefolo/eGO/pkg/service"
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
	Creds      credentials.TransportCredentials
	DialOpts   []grpc.DialOption
	Host       string
	Port       int
	ListenAddr string
	Mux        *runtime.ServeMux
	Services   service.List
}

// New creates a gateway.
func New(ctx context.Context, conf config.Service, opts []grpc.DialOption, regFn registerFromEndpoint) (Gateway, error) {

	gw := Gateway{
		Mux:      runtime.NewServeMux(),
		DialOpts: opts,
	}

	//loadCerts(gw)

	gw.ListenAddr = fmt.Sprintf(":%d", conf.Server.GatewayPort)

	//	if err := regFn(ctx, gw.Mux, gw.ListenAddr, gw.DialOpts); err != nil {
	//		return nil, fmt.Errorf("failed to register gateway endpoints: %v", err)
	//	}

	//gw.Services = svcList

	return gw, nil
}

// LoadEndpoint ...
func (gw *Gateway) LoadEndpoint(ctx context.Context, serviceConf config.Service, regFn registerFromEndpoint) {

	addr := fmt.Sprintf(":%d", serviceConf.Server.GatewayPort)
	gw.ListenAddr = addr

	if err := regFn(ctx, gw.Mux, addr, gw.DialOpts); err != nil {
		log.Fatalf("failed to register endpoint: %v", err)
	}
}

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

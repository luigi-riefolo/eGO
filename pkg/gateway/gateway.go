package gateway

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/luigi-riefolo/alfa/pkg/config"
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
}

// LoadEndpoint ...
func (gw *Gateway) LoadEndpoint(ctx context.Context, serviceConf config.Service, regFn registerFromEndpoint) {

	//addr := fmt.Sprintf("%s:%d", serviceConf.Server.Host, serviceConf.Server.Port)
	addr := fmt.Sprintf(":%d", serviceConf.Server.Port)
	log.Printf("%s %s", serviceConf.Name, addr)

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

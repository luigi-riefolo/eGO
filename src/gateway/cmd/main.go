package main

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// Project packages
	"github.com/luigi-riefolo/eGO/src/alfa"
	"github.com/luigi-riefolo/eGO/src/gateway"
	gatewaypb "github.com/luigi-riefolo/eGO/src/gateway/pb"

	"github.com/luigi-riefolo/eGO/pkg/config"
	gw "github.com/luigi-riefolo/eGO/pkg/gateway"
	"github.com/luigi-riefolo/eGO/pkg/service"
)

var (
	conf config.Config
)

func init() {

	var err error
	conf, err = config.GetConfig()
	if err != nil {
		log.Fatal("Config file could not be parsed")
	}
}

// runEndPoints ...
func runEndPoints() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	service.ContextTimeout(ctx)

	// set up the gateway instance
	gwSrv := gw.Gateway{}
	gwSrv.Mux = runtime.NewServeMux()
	gwSrv.DialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(config.BackoffDelay),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxMsgSize)),
	}
	/*	grpc.WithTransportCredentials(gw.creds)}
		grpc.WithStatsHandler()
		grpc.WithUnaryInterceptor())
	*/

	//loadCerts(gw)

	gwSrv.ListenAddr = fmt.Sprintf(":%d", conf.Gateway.Server.GatewayPort)

	log.Println("Loading service endpoints")

	// List of gateway endpoints
	gwSrv.LoadEndpoint(ctx, conf.Gateway, gatewaypb.RegisterGatewayHandlerFromEndpoint)

	gwSrv.Services = map[string]*service.Service{
		// List of services
		conf.Gateway.Name: gateway.Serve(ctx, conf),
		conf.Alfa.Name:    alfa.Serve(ctx, conf),
	}

	return gwSrv.ListenAndServe()
}

func main() {

	if err := runEndPoints(); err != nil {
		log.Fatal(err)
	}
}

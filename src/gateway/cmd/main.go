package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// Project packages

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/luigi-riefolo/eGO/pkg/config"
	gw_pkg "github.com/luigi-riefolo/eGO/pkg/gateway"
	"github.com/luigi-riefolo/eGO/pkg/service"
	"github.com/luigi-riefolo/eGO/src/gateway"
	gatewaypb "github.com/luigi-riefolo/eGO/src/gateway/pb"
)

// TODO: create clients
// load and reload vars
// metrics
// auth
// trace ids
// check which log lib is best
// CORS: do we always need it or activate it for specific endpoints?????????

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

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//server.ContextTimeout(ctx)

	// set up the gateway instance
	gw := gw_pkg.Gateway{}
	gw.Mux = runtime.NewServeMux()
	gw.DialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(config.BackoffDelay),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxMsgSize)),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
	}
	/*	grpc.WithTransportCredentials(gw.creds)}
		grpc.WithStatsHandler()
		grpc.WithUnaryInterceptor())
	*/

	//loadCerts(gw)

	gw.ListenAddr = fmt.Sprintf(":%d", conf.Alfa.Server.GatewayPort)

	log.Println("Loading service endpoints")

	// List of gateway endpoints
	gw.LoadEndpoint(ctx, conf.Gateway, gatewaypb.RegisterGatewayHandlerFromEndpoint)

	gw.Services = map[string]*service.Service{
		// List of services
		//conf.Alfa.Name:    alfa.Serve(ctx, conf),
		conf.Gateway.Name: gateway.Serve(ctx, conf),
	}

	log.Fatal(gw.ListenAndServe())
}

/*
func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	service.ContextTimeout(ctx)

	// set up the gateway instance
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(config.BackoffDelay),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxMsgSize)),
		//		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		//		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		/*	grpc.WithTransportCredentials(gw.creds)}
			grpc.WithStatsHandler()
			grpc.WithUnaryInterceptor())

	}

services := service.List{
		// List of services
		conf.Gateway.Name: gateway.Serve(ctx, conf),
		conf.Alfa.Name:    alfa.Serve(ctx, conf),
	}
	gwSrv, err := gw.New(ctx, conf.Gateway, opts, pb.RegisterGatewayHandlerFromEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(gwSrv.ListenAndServe())
}
*/

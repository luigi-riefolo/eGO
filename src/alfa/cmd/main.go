package main

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// Project packages
	"github.com/luigi-riefolo/eGO/src/alfa"
	alfapb "github.com/luigi-riefolo/eGO/src/alfa/pb"
	"github.com/luigi-riefolo/eGO/src/beta"

	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/gateway"
	"github.com/luigi-riefolo/eGO/pkg/server"
)

// TODO: create clients
// load and reload vars
// metrics
// auth
// trace ids
// check which log lib is best
// CORS: do we always need it or activate it for specific endpoints?????????

var (
	conf    config.Config
	service string
	version string
	build   string
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

	server.ContextTimeout(ctx)

	// set up the gateway instance
	gw := gateway.Gateway{}
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
	gw.LoadEndpoint(ctx, conf.Alfa, alfapb.RegisterAlfaServiceHandlerFromEndpoint)

	gw.Services = map[string]*server.Server{
		// List of services
		conf.Alfa.Name: alfa.Serve(ctx, conf),
		conf.Beta.Name: beta.Serve(ctx, conf),
	}

	log.Printf("Starting '%s' service, version '%s' (%s)\n", service, version, build)

	return gw.ListenAndServe()
}

func main() {

	if err := runEndPoints(); err != nil {
		log.Fatal(err)
	}
}

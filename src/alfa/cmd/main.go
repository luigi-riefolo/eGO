package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// Project packages
	alfa "github.com/luigi-riefolo/alfa/src/alfa"
	alfapb "github.com/luigi-riefolo/alfa/src/alfa/pb"
	beta "github.com/luigi-riefolo/alfa/src/beta"

	"github.com/luigi-riefolo/alfa/pkg/config"
	"github.com/luigi-riefolo/alfa/pkg/gateway"
	"github.com/luigi-riefolo/alfa/pkg/server"
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

	log.Printf("Starting %s", conf.Alfa.Name)
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

	// List of services
	alfa.Serve(ctx, conf)
	beta.Serve(ctx, conf)

	log.Printf("Gateway listening on %s\n", gw.ListenAddr)

	log.Fatal(http.ListenAndServe(gw.ListenAddr, gw.Mux))

	return nil
}

func main() {

	if err := runEndPoints(); err != nil {
		log.Fatal(err)
	}
}

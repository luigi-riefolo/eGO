package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/luigi-riefolo/eGO/pkg/config"
)

// Server ...
type Server struct {
	name              string
	gRPC              *grpc.Server
	lis               net.Listener
	health            *health.Server
	stopHealthMonitor chan struct{}
}

// GetgRPCServer ...
func (srv *Server) GetgRPCServer() *grpc.Server {
	return srv.gRPC
}

// NewServer creates a new gRPC server.
func NewServer(conf config.Service) *Server {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: recovery handlers have to be last in the chain
	// so that other middleware (e.g. logging) can operate
	// on the recovered state instead of being directly
	// affected by any panic
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(RecoveryInterceptor),
	}
	opts := []grpc.ServerOption{
		grpc.RPCCompressor(grpc.NewGZIPCompressor()),
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		grpc.MaxMsgSize(config.MaxMsgSize),
		grpc_middleware.WithUnaryServerChain(
			//grpc_opentracing.UnaryServerInterceptor(),

			grpc_prometheus.UnaryServerInterceptor,
			//       grpc_auth.UnaryServerInterceptor(myAuthFunction),
			//grpc_validator
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	}

	srv := &Server{
		name:   conf.Name,
		gRPC:   grpc.NewServer(opts...),
		lis:    l,
		health: health.NewServer(),
	}

	return srv
}

// Listen starts listening on a port.
func (srv *Server) Listen() {

	srv.setServicesHealth()

	go func() {
		if err := srv.gRPC.Serve(srv.lis); err != nil {
			log.Printf("Cannot listen and serve: %v", err)
		}
	}()

	log.Printf("Service %s listening on: %s", srv.name, srv.lis.Addr().String())
}

// HandleSig ...
func (srv *Server) HandleSig(sig string) {

	log.Printf("Handling %s signal in %s\n", sig, srv.name)
	for s := range srv.GetgRPCServer().GetServiceInfo() {
		log.Printf("Shutting down %s\n", s)
	}
	srv.gRPC.Stop()
}

// ContextTimeout ...
// TODO: refactor
func ContextTimeout(ctx context.Context) {
	//	for {
	select {
	case <-ctx.Done():
		if ctx.Err() == context.Canceled {
			log.Printf("Context: %v", ctx.Err())
		} else if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Context: %v", ctx.Err())
		}
	default:
		time.Sleep(1 * time.Second)
	}
	//	}
}

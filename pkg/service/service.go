package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/luigi-riefolo/eGO/pkg/auth"
	"github.com/luigi-riefolo/eGO/pkg/config"
)

// Service represents a service.
type Service struct {
	name              string
	gRPC              *grpc.Server
	lis               net.Listener
	health            *health.Server
	stopHealthMonitor chan struct{}
}

// List represents a map of service. The key is the
// name of service and the value its struct representation.
type List map[string]*Service

// New creates a new gRPC service.
func New(conf config.Service) *Service {

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
		//grpc.InTapHandle(ratelimit.RateLimiter),
		grpc.RPCCompressor(grpc.NewGZIPCompressor()),
		grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		grpc.MaxMsgSize(config.MaxMsgSize),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			LoggingInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			TimeoutInterceptor,
			grpc_prometheus.UnaryServerInterceptor,
			grpc_auth.UnaryServerInterceptor(auth.Authorize),
			//		grpc_ctxtags.UnaryServerInterceptor(
			//			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			//		grpc_zap.UnaryServerInterceptor(llog.Logger, llog.LogOpts...),
			//grpc_validator
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
	}

	svc := &Service{
		name:   conf.Name,
		gRPC:   grpc.NewServer(opts...),
		lis:    l,
		health: health.NewServer(),
	}

	return svc
}

// GetgRPCServer ...
func (svc *Service) GetgRPCServer() *grpc.Server {
	return svc.gRPC
}

// Listen starts listening on a port.
func (svc *Service) Listen() {

	svc.setServicesHealth()

	go func() {
		if err := svc.gRPC.Serve(svc.lis); err != nil {
			log.Printf("Cannot listen and serve: %v", err)
		}
	}()

	log.Printf("Service %s listening on: %s", svc.name, svc.lis.Addr().String())
}

// HandleSig handles signalling events.
func (svc *Service) HandleSig(sig string) {

	log.Printf("Handling %s signal in %s\n", sig, svc.name)
	for s := range svc.GetgRPCServer().GetServiceInfo() {
		log.Printf("Shutting down %s\n", s)
	}
	svc.gRPC.Stop()
}

// ContextTimeout ...
func ContextTimeout(ctx context.Context) {
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
}

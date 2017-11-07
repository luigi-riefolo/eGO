package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/luigi-riefolo/eGO/pkg/config"
)

// Server ...
type Server struct {
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
func NewServer(port int) *Server {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
		gRPC:   grpc.NewServer(opts...),
		lis:    l,
		health: health.NewServer(),
	}

	srv.HandleSig()

	return srv
}

// Listen starts listening on a port.
func (srv *Server) Listen() {

	srv.setServicesHealth()

	// initialize Prometheus metrics
	//srv.StartPrometheus()

	go func() {
		if err := srv.gRPC.Serve(srv.lis); err != nil {
			log.Printf("Cannot listen and serve: %v", err)
		}
	}()
}

// HandleSig ...
func (srv *Server) HandleSig() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	done := make(chan int)
	go func() {
		for {
			s := <-sig
			switch s {

			case syscall.SIGINT:
				log.Println("Handling interrupt signal")
				srv.gRPC.GracefulStop()
				done <- 0

			case syscall.SIGTERM:
				log.Println("Handling termination signal")
				srv.gRPC.GracefulStop()
				done <- 0

			case syscall.SIGQUIT:
				log.Println("Handling quit signal")
				srv.gRPC.GracefulStop()
				done <- 0

			default:
				log.Println("Unknown signal")
				done <- 1
			}
		}
	}()
}

// ContextTimeout ...
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

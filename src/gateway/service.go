package gateway

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/service"
	alfapb "github.com/luigi-riefolo/eGO/src/alfa/pb"
	pb "github.com/luigi-riefolo/eGO/src/gateway/pb"
)

// Service represents the Gateway service.
type Service struct {
	conf       config.Config
	alfaClient alfapb.AlfaServiceClient
}

// New initialises the Gateway service server.
func New(ctx context.Context, conf config.Config) pb.GatewayServer {
	gateway := &Service{
		conf: conf,
	}

	/*	c, err := client.Get(ctx, conf.Alfa)
		if err != nil {
			log.Fatalf("%v", err)
		}
		gateway.alfaClient = c.(alfapb.AlfaServiceClient)
	*/
	return gateway
}

/*
POC
a gRPC server can handle multiple services by registering them.
need to test which approach is faster:
- each service runs as stand-alone
- a server handles multiple services

type Gamma struct {
}

func (s *Gamma) TestOp(ctx context.Context, req *empty.Empty) (*alfapb.Message, error) {
	return &alfapb.Message{}, nil
}

func NewGammaService() pb.GammaServiceServer {
	gammaSrv := &Gamma{}
	return gammaSrv
}
*/

// Serve starts listening and serving requests.
func Serve(ctx context.Context, conf config.Config) *service.Service {

	svc := service.New(conf.Gateway)

	pb.RegisterGatewayServer(
		svc.GetgRPCServer(),
		New(ctx, conf))

	//pb.RegisterGammaServiceServer(srv.GetgRPCServer(), NewGammaService())

	svc.Listen()

	// initialize Prometheus metrics
	svc.StartPrometheus()

	return svc
}

// Get ...
func (s *Service) Get(ctx context.Context, req *empty.Empty) (*pb.Message, error) {

	fmt.Println("GET")
	err := grpc.SendHeader(ctx, metadata.New(map[string]string{
		"foo": "foo1",
		"bar": "bar1",
	}))
	if err != nil {
		return nil, err
	}
	msg := &pb.Message{
		Msg: "Hi there!!!",
	}
	return msg, nil
}

// Set ...
func (s *Service) Set(ctx context.Context, req *pb.Message) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Test ...
func (s *Service) Test(ctx context.Context, req *empty.Empty) (*alfapb.Message, error) {

	msg, err := s.alfaClient.Test(ctx, req)
	if err != nil {
		log.Printf("Test: %v", err)
	}

	return msg, nil
}

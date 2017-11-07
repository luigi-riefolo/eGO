package alfa

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/luigi-riefolo/eGO/pkg/client"
	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/server"
	pb "github.com/luigi-riefolo/eGO/src/alfa/pb"
	betapb "github.com/luigi-riefolo/eGO/src/beta/pb"
)

// Service represents the Alfa service.
type Service struct {
	conf       config.Config
	betaClient betapb.BetaServiceClient
}

// NewAlfaService initialises the Alfa service server.
func NewAlfaService(ctx context.Context, conf config.Config) pb.AlfaServiceServer {
	alfa := &Service{
		conf: conf,
	}

	// TODO: REMOVE ME
	opts := []grpc.DialOption{}

	c, err := client.Get(ctx, conf.Beta, opts...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	alfa.betaClient = c.(betapb.BetaServiceClient)

	return alfa
}

/*
POC
a gRPC server can handle multiple services by registering them.
need to test which approach is faster:
- each service runs as stand-alone
- a server handles multiple services

type Gamma struct {
}

func (s *Gamma) TestOp(ctx context.Context, req *empty.Empty) (*betapb.Message, error) {
	return &betapb.Message{}, nil
}

func NewGammaService() pb.GammaServiceServer {
	gammaSrv := &Gamma{}
	return gammaSrv
}
*/

// Serve starts listening and serving requests.
func Serve(ctx context.Context, conf config.Config) {

	srv := server.NewServer(conf.Alfa.Server.Port)

	alfaService := NewAlfaService(ctx, conf)

	pb.RegisterAlfaServiceServer(srv.GetgRPCServer(), alfaService)

	//pb.RegisterGammaServiceServer(srv.GetgRPCServer(), NewGammaService())

	srv.Listen()
}

// Get ...
func (s *Service) Get(ctx context.Context, req *empty.Empty) (*pb.Message, error) {

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
func (s *Service) Test(ctx context.Context, req *empty.Empty) (*betapb.Message, error) {

	msg, err := s.betaClient.Test(ctx, req)
	if err != nil {
		log.Printf("Test: %v", err)
	}

	return msg, nil
}

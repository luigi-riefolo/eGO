package beta

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/luigi-riefolo/alfa/pkg/client"
	"github.com/luigi-riefolo/alfa/pkg/config"
	"github.com/luigi-riefolo/alfa/pkg/server"
	pb "github.com/luigi-riefolo/alfa/src/beta/pb"
	omegapb "github.com/luigi-riefolo/alfa/src/omega/pb"
)

// Service implements a BetaServiceServer
type Service struct {
	conf        config.Config
	omegaClient omegapb.OmegaServiceClient
}

// NewService initialises the Alfa service server.
func NewService(ctx context.Context, conf config.Config) pb.BetaServiceServer {
	beta := &Service{
		conf: conf,
	}

	opts := []grpc.DialOption{}
	c, err := client.Get(ctx, conf.Omega, opts...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	beta.omegaClient = c.(omegapb.OmegaServiceClient)

	return beta
}

// Serve starts listening and serving requests.
func Serve(ctx context.Context, serviceConf config.Config) {

	srv := server.NewServer(serviceConf.Beta.Server.Port)

	betaServer := NewService(ctx, serviceConf)

	pb.RegisterBetaServiceServer(srv.GetgRPCServer(), betaServer)

	srv.Listen()
}

// Test ...
func (s *Service) Test(ctx context.Context, req *empty.Empty) (*pb.Message, error) {
	msg := &pb.Message{
		Msg: "OK BETA!!!",
	}
	/*
		_, err := s.omegaClient.Echo(ctx, &omegapb.Message{Msg: "YES"})
		if err != nil {
			log.Printf("Omega client err: %#v\n", err)
		}
	*/
	return msg, nil
}

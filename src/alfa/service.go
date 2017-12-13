package alfa

import (
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/service"
	pb "github.com/luigi-riefolo/eGO/src/alfa/pb"
	omegapb "github.com/luigi-riefolo/eGO/src/omega/pb"
)

// Service implements a AlfaServiceServer
type Service struct {
	conf        config.Config
	omegaClient omegapb.OmegaServiceClient
}

// New initialises the Alfa service.
func New(ctx context.Context, conf config.Config) pb.AlfaServiceServer {
	alfa := &Service{
		conf: conf,
	}
	/*
		c, err := client.Get(ctx, conf.Omega)
		if err != nil {
			log.Fatalf("%v", err)
		}
		alfa.omegaClient = c.(omegapb.OmegaServiceClient)
	*/
	return alfa
}

// Serve starts listening and serving requests.
func Serve(ctx context.Context, conf config.Config) *service.Service {

	svc := service.New(conf.Alfa)

	pb.RegisterAlfaServiceServer(
		svc.GetgRPCServer(),
		New(ctx, conf))

	svc.Listen()

	return svc
}

// Test ...
func (s *Service) Test(ctx context.Context, req *empty.Empty) (*pb.Message, error) {
	msg := &pb.Message{
		Msg: "OK ALFA!!!",
	}

	/*	_, err := s.omegaClient.Echo(ctx, &omegapb.Message{Msg: "YES"})
		if err != nil {
			log.Printf("Omega client err: %#v\n", err)
		}
	*/
	return msg, nil
}

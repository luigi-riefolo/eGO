package omega

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/luigi-riefolo/eGO/pkg/config"
	pb "github.com/luigi-riefolo/eGO/src/omega/pb"
)

// Service implements an OmegaServiceServer.
type Service struct {
	conf config.Config
	// List of clients
}

// NewService ...
func NewService(conf config.Config) pb.OmegaServiceServer {
	omega := &Service{
		conf: conf,
	}
	return omega
}

// Echo ...
func (s *Service) Echo(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	fmt.Printf("OMEGA ECHO\n")
	return msg, nil
}

// Dump ...
/*func (s *Service) Dump(ctx context.Context, stream pb.OmegaService_DumpServer) error {

	for k, v := range []string{"a", "b"} {
		err := stream.Send(&pb.DumpItem{
			Key: fmt.Sprintf("%d", k),
			Val: v,
		})
		if err != nil {
			log.Println("ERROR: ", err)
		}
	}
	return nil
}
*/

package client

import (
	"google.golang.org/grpc"

	pb "github.com/luigi-riefolo/eGO/src/omega/pb"
)

// New creates an Omega client.
func New(conn *grpc.ClientConn) pb.OmegaServiceClient {
	return pb.NewOmegaServiceClient(conn)
}

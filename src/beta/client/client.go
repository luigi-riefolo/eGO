package client

import (
	"google.golang.org/grpc"

	pb "github.com/luigi-riefolo/eGO/src/beta/pb"
)

// New creates a Beta client.
func New(conn *grpc.ClientConn) pb.BetaServiceClient {
	return pb.NewBetaServiceClient(conn)
}

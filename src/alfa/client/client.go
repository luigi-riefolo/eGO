package client

import (
	"google.golang.org/grpc"

	pb "github.com/luigi-riefolo/eGO/src/alfa/pb"
)

// New creates a Alfa client.
func New(conn *grpc.ClientConn) pb.AlfaServiceClient {
	return pb.NewAlfaServiceClient(conn)
}

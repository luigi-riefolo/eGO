package client

import (
	"context"
	"fmt"
	"log"

	"github.com/luigi-riefolo/eGO/pkg/config"
	betaclient "github.com/luigi-riefolo/eGO/src/beta/client"
	omegaclient "github.com/luigi-riefolo/eGO/src/omega/client"
	"google.golang.org/grpc"
)

var (
	// DefaultDialOpts is a list of default DialOptions.
	DefaultDialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithCompressor(grpc.NewGZIPCompressor()),
		grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		grpc.WithTimeout(config.ClientRequestTimeout),
		grpc.WithMaxMsgSize(config.MaxMsgSize),
		//func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption
		//func WithStatsHandler(h stats.Handler) DialOption
		//func WithTransportCredentials(creds credentials.TransportCredentials) DialOption
		//func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption
	}
)

func newConn(ctx context.Context, conf config.Service, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)

	// Set up a connection to the server.
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}

	// TODO: close client connections on interrupts
	//conn.Close()

	log.Printf("Client for %s connected on: %s", conf.Name, target)

	return conn, nil
}

// Get ...
func Get(ctx context.Context, conf config.Service, opts ...grpc.DialOption) (interface{}, error) {

	var client interface{}

	if len(opts) == 0 {
		opts = DefaultDialOpts
	}

	conn, err := newConn(ctx, conf, opts...)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialise the %s client connection: %v",
			conf.Name, err)
	}

	// apparently the connecting state is incremental and that's why
	// WaitForStateChange returns quick or maybe it doesn't work yet,
	// for the func is experimental
	/*
		sourceState := conn.GetState()
		fmt.Println("STATE : ", sourceState.String())
		s := conn.WaitForStateChange(ctx, sourceState)
		fmt.Println("CHANGE STATE: ", s)
	*/

	switch conf.ShortName {

	case "beta":
		client = betaclient.New(conn)
	case "omega":
		client = omegaclient.New(conn)
	default:
		return nil, fmt.Errorf("Client %s not supported", conf.ShortName)
	}

	return client, err
}

package client

import (
	"context"

	"google.golang.org/grpc"
)

// RateLimiterInterceptor ...
func RateLimiterInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	return nil
}

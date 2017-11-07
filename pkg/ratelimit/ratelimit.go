package ratelimit

import (
	"golang.org/x/net/context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
)

// Tap  defines the function handles which are executed on
// the transport layer of gRPC-Go and related information
type Tap struct {
	lim *rate.Limiter
}

// NewTap ...
func NewTap() *Tap {
	// 150 tokens per second
	// 5 token
	return &Tap{rate.NewLimiter(150, 5)}
}

// ServerInHandle ...
func (t *Tap) ServerInHandle(ctx context.Context, info *tap.Info) (context.Context, error) {
	//info.FullMethodName
	//t.lim.Reserve()
	//t.lim.Wait()
	if !t.lim.Allow() {
		return nil, status.Errorf(codes.ResourceExhausted, "service is over rate limit")
	}
	return ctx, nil
}

// RateLimiter ...
func RateLimiter() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxConcurrentStreams(64),
		grpc.InTapHandle(NewTap().ServerInHandle),
	}
	// l = netutil.LimitListener(l, 1024)
	//srv.Serve(l)
}

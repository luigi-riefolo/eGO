package service

import (
	"context"
	"log"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/tap"

	"github.com/luigi-riefolo/eGO/pkg/auth"
	"github.com/luigi-riefolo/eGO/pkg/config"
	"github.com/luigi-riefolo/eGO/pkg/errors"
	"github.com/luigi-riefolo/eGO/pkg/util"
)

var t = rate.NewLimiter(1, 1)

// RateLimiter limits the amount of request for a user.
func RateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	if !t.Allow() {
		return nil, errors.Err(codes.ResourceExhausted)
	}
	//	fmt.Printf("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA: %v -- %v\n", t.Limit(), t.Burst())

	return ctx, nil
}

// RecoveryInterceptor intercepts panic signals and tries to recover.
func RecoveryInterceptor(p interface{}) error {
	log.Println("Trying to recover from panic")
	if r := recover(); r != nil {
		return grpc.Errorf(codes.Internal, "Could not recover from panic: %v", r)
	}
	return nil
}

// TimeoutInterceptor cancels a request if it does
// not complete within 'requestTimeout' seconds.
func TimeoutInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		config.ServerRequestTimeout)
	defer cancel()

	return handler(ctx, req)
}

// LoggingInterceptor ...
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Printf("%s\n", info.FullMethod)
	return handler(ctx, req)
}

// UnaryInterceptor ...
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	newCtx := context.WithValue(ctx, util.CtxKey("user_id"), "USER_ID")
	return handler(newCtx, req)
}

// AuthInterceptor ...
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := auth.Authorize(ctx)
	if err != nil {
		return nil, errors.Err(codes.PermissionDenied, err)
	}
	return handler(newCtx, req)
}

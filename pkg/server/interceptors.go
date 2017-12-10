package server

import (
	"context"
	"log"

	"github.com/luigi-riefolo/eGO/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ctxKey string

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

// UnaryInterceptor ...
// TODO: example of UnaryInterceptor, use it for auth
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	newCtx := context.WithValue(ctx, ctxKey("user_id"), "USER_ID")
	return handler(newCtx, req)
}

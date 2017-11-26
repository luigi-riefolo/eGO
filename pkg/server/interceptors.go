package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const requestTimeout = 10 * time.Second

type cxtKey string

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
		requestTimeout)
	defer cancel()

	return handler(ctx, req)
}

// UnaryInterceptor ...
// TODO: example of UnaryInterceptor, use it for auth
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var k cxtKey
	k = "user_id"
	newCtx := context.WithValue(ctx, k, "USER_ID")
	return handler(newCtx, req)
}

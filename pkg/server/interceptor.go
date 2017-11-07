package server

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// RecoveryInterceptor intercepts panic signals and tries to recover.
func RecoveryInterceptor(p interface{}) error {
	log.Println("Trying to recover from panic")
	if r := recover(); r != nil {
		return grpc.Errorf(codes.Internal, "Could not recover from panic: %v", r)
	}
	return nil
}

// TimeoutInterceptor ...
func TimeoutInterceptor(p interface{}) error {
	//	ctx, cancel := context.WithTimeout(
	//		context.Background(),
	//		config.RequestTimeout)
	//	defer cancel()

	return nil
}

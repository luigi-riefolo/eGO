package errors

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	rateLimitErr       = grpc.Errorf(codes.ResourceExhausted, "service is over rate limit")
	unauthenticatedErr = grpc.Errorf(codes.Unauthenticated, "invalid auth token")
	unauthorizedErr    = grpc.Errorf(codes.PermissionDenied, "user is not authorized")
	unknownErr         = grpc.Errorf(codes.Internal, "unknown error")
)

func fmtErr(err error, msg ...interface{}) error {

	desc := err.Error()
	for _, d := range msg {
		desc += fmt.Sprintf(": %s", d)
	}

	err = fmt.Errorf("%s", desc)

	return err
}

// Err returns a formatted error.
func Err(code codes.Code, msg ...interface{}) error {

	var err error
	switch code {
	case codes.PermissionDenied:
		err = unauthorizedErr
	case codes.ResourceExhausted:
		err = rateLimitErr
	case codes.Unauthenticated:
		err = unauthenticatedErr

	case codes.Unknown:
		fallthrough
	default:
		err = unknownErr
	}

	if len(msg) > 0 {
		err = fmtErr(err, msg...)
	}

	return err
}

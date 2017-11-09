package lib

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"cloud.google.com/go/trace"
)

//func getNearbyPoints(ctx context.Context, lat, lon float64) []geoindex.Point {
//span := trace.FromContext(ctx).NewChild("geo.getNearbyPoints")
//defer span.Finish()

// EnableGRPCTracingDialOption enables tracing of requests that are sent over a gRPC connection.
// Modified version of: https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/trace/trace.go#L242-L265
var EnableGRPCTracingDialOption = grpc.WithUnaryInterceptor(grpc.UnaryClientInterceptor(clientInterceptor))

func clientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// trace current request w/ child span
	span := trace.FromContext(ctx).NewChild(method)
	defer span.Finish()

	// new metadata, or copy of existing
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	// append trace header to context metadata
	// header specification: https://cloud.google.com/trace/docs/faq
	headerKey := ""
	md[headerKey] = append(
		md[headerKey], fmt.Sprintf("%s/%d;o=1", span.TraceID(), 0),
	)
	ctx = metadata.NewIncomingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

// EnableGRPCTracingServerOption enables parsing google trace header from metadata
// and adds a new child span to the incoming request context.
func EnableGRPCTracingServerOption(traceClient *trace.Client) grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor(traceClient))
}

func serverInterceptor(traceClient *trace.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// fetch metadata from request context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		headerKey := ""
		header := strings.Join(md[headerKey], "")

		// create new child span from google trace header, add to
		// current request context
		span := traceClient.SpanFromHeader(info.FullMethod, header)
		defer span.Finish()
		ctx = trace.NewContext(ctx, span)

		return handler(ctx, req)
	}
}

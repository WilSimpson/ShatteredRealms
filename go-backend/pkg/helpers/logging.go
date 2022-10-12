package helpers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryLogRequest() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Info(info.FullMethod)
		return handler(ctx, req)
	}
}
func StreamLogRequest() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Info(info.FullMethod)
		return handler(srv, stream)
	}
}

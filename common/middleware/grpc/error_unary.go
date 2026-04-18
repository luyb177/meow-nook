package grpcmw

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"google.golang.org/grpc"
)

// ErrorUnaryServer converts application errors (*errorx.AppError etc.)
// returned by handlers into gRPC status errors with details.
// Place it early in interceptor chain so downstream interceptors see
// the converted gRPC error.
func ErrorUnaryServer() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, errorx.ToGRPC(err)
		}
		return resp, nil
	}
}

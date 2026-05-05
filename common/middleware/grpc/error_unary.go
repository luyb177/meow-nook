package grpcmw

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"google.golang.org/grpc"
)

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

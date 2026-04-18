package grpcmw

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/luyb177/meow-nook/common/logger"
)

// PropagateRequestIDUnaryClient
// NOTE: 这个就是用来注入 request_id 的，但是没找到对应的地方使用
func PropagateRequestIDUnaryClient() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if rid := logger.RequestIDFromContext(ctx); rid != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

package grpcmw

import (
	"context"

	"github.com/luyb177/meow-nook/common/logger"
	"google.golang.org/grpc/metadata"
)

func InjectRequestID(ctx context.Context) context.Context {
	if rid := logger.RequestIDFromContext(ctx); rid != "" {
		return metadata.AppendToOutgoingContext(ctx, "x-request-id", rid)
	}
	return ctx
}

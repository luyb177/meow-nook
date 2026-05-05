package grpcmw

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/luyb177/meow-nook/common/logger"
)

func AccessLogUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		var rid string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("x-request-id"); len(vals) > 0 {
				rid = vals[0]
			}
		}
		ctx = logger.WithRequestID(ctx, rid)

		if tid := logger.TraceIDFromOTel(ctx); tid != "" {
			ctx = logger.WithTraceID(ctx, tid)
		}

		resp, err := handler(ctx, req)

		code := status.Code(err).String()

		remote := ""
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			remote = p.Addr.String()
		}

		l := logger.FromContext(ctx)
		fields := []zap.Field{
			zap.String("grpc_method", info.FullMethod),
			zap.String("grpc_code", code),
			zap.Int64("latency_ms", time.Since(start).Milliseconds()),
			zap.String("peer", remote),
		}

		// todo 这里可以按照 grpc code 来区分日志级别，目前先统一用 info 级别
		if err != nil {
			l.Info("grpc access", append(fields, zap.Error(err))...)
		} else {
			l.Info("grpc access", fields...)
		}
		return resp, err
	}
}

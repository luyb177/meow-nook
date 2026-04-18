package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/grpcx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/gateway/internal/config"
	"github.com/luyb177/meow-nook/service/gateway/internal/handler"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file path")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// Initialize structured logger.
	if err := logger.Init(logger.Config{
		Level:    c.Log.Level,
		Encoding: c.Log.Encoding,
	}); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer logger.Sync()

	// Register a unified error handler so that *AppError values are serialized
	// as { "code": N, "msg": "…" } instead of a plain string.
	// gRPC status errors carrying an ErrorDetail are decoded first so that
	// business codes and messages from upstream services reach the client
	// unchanged.
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		// 1. Try to extract business error from a gRPC status (with or without
		//    ErrorDetail details in the status).
		if ae := grpcx.ExtractAppError(err); ae != nil {
			return ae.HTTPStatus(), map[string]interface{}{
				"code": ae.Code,
				"msg":  ae.Msg,
			}
		}
		// 2. Plain *AppError returned directly (e.g. from gateway-local logic).
		if ae, ok := err.(*errorx.AppError); ok {
			return ae.HTTPStatus(), map[string]interface{}{
				"code": ae.Code,
				"msg":  ae.Msg,
			}
		}
		// 3. Unknown error – treat as internal.
		return http.StatusInternalServerError, map[string]interface{}{
			"code": errorx.CodeInternalError,
			"msg":  "服务器内部错误",
		}
	})

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	logger.Info(fmt.Sprintf("meow-nook gateway starting on %s:%d", c.Host, c.Port))
	server.Start()
}

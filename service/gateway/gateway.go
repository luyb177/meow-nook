package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/luyb177/meow-nook/common/errorx"
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

	// Initialise structured logger.
	if err := logger.Init(logger.Config{
		Level:    c.Log.Level,
		Encoding: c.Log.Encoding,
	}); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	defer logger.Sync()

	// Register a unified error handler so that *AppError values are serialised
	// as { "code": N, "msg": "…" } instead of a plain string.
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		if ae, ok := err.(*errorx.AppError); ok {
			return ae.HTTPStatus(), map[string]interface{}{
				"code": ae.Code,
				"msg":  ae.Msg,
			}
		}
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

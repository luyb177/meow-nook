// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"flag"
	"fmt"

	"github.com/luyb177/meow-nook/common/logger"
	httpmw "github.com/luyb177/meow-nook/common/middleware/http"
	"github.com/luyb177/meow-nook/service/gateway/internal/config"
	"github.com/luyb177/meow-nook/service/gateway/internal/handler"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "service/gateway/etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	err := logger.Init(c.Logger)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(httpmw.AccessLog)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logger.Info("Starting server...", zap.String("listen_on", fmt.Sprintf("%s:%d", c.Host, c.Port)))
	server.Start()
}

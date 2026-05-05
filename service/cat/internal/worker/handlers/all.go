package handlers

import (
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

func BuildRegistry(svcCtx *svc.ServiceContext) *kafka.Registry {
	r := kafka.NewRegistry()

	// 以后新增任务类型就在这里继续 r.Register(...)
	return r
}

package handlers

import (
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
)

// 一个集中注册的地方

func BuildRegistry(svcCtx *svc.ServiceContext) *kafka.Registry {
	r := kafka.NewRegistry()

	th := NewTestHandler()
	r.Register(task.TypeSuccessTest, th.OnSuccess)
	r.Register(task.TypeFailTest, th.OnFail)

	vh := NewVerifyCodeHandler(svcCtx)
	r.Register(task.TypeSendVerifyCode, vh.SendVerifyCode)

	// 以后新增任务类型就在这里继续 r.Register(...)
	return r
}

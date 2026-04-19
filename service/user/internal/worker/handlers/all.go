package handlers

import (
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
)

// 一个集中注册的地方

func BuildRegistry() *kafka.Registry {
	r := kafka.NewRegistry()

	th := TestHandlers{}
	r.Register(task.TypeSuccessTest, th.OnSuccess)
	r.Register(task.TypeFailTest, th.OnFail)

	// 以后新增任务类型就在这里继续 r.Register(...)
	return r
}

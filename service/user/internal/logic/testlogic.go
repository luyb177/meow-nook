package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/mq/kafka"
	"github.com/luyb177/meow-nook/service/user/internal/pkg/task"
	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"
)

type TestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestLogic {
	return &TestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Test Auth
func (l *TestLogic) Test(in *v1.TestReq) (*v1.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t1 := kafka.NewTypedTask(task.TypeSuccessTest, "test-success", task.TestPayload{Name: "test-success"})
	if err := l.svcCtx.KafkaProducer.Dispatch(ctx, t1); err != nil {
		return nil, errorx.WrapInternal("发送 success_test 任务失败", err)
	}

	t2 := kafka.NewTypedTask(task.TypeFailTest, "test-fail", task.TestPayload{Name: "test-fail"})
	if err := l.svcCtx.KafkaProducer.Dispatch(ctx, t2); err != nil {
		return nil, errorx.WrapInternal("发送 fail_test 任务失败", err)
	}

	return &v1.Response{}, nil
}

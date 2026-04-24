package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCatTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCatTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatTasksLogic {
	return &ListCatTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 正式任务列表
func (l *ListCatTasksLogic) ListCatTasks(in *v1.ListCatTasksRequest) (*v1.ListCatTasksResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ListCatTasksResponse{}, nil
}

package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimCatTaskLogic {
	return &ClaimCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请认领任务
func (l *ClaimCatTaskLogic) ClaimCatTask(in *v1.ClaimCatTaskRequest) (*v1.ClaimCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ClaimCatTaskResponse{}, nil
}

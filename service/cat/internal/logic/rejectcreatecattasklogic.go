package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectCreateCatTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectCreateCatTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectCreateCatTaskLogic {
	return &RejectCreateCatTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectCreateCatTaskLogic) RejectCreateCatTask(in *v1.RejectCreateCatTaskRequest) (*v1.RejectCreateCatTaskResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.RejectCreateCatTaskResponse{}, nil
}

package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/task/internal/svc"
	"github.com/luyb177/meow-nook/service/task/pb/task/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClaimTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimTaskLogic {
	return &ClaimTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClaimTaskLogic) ClaimTask(in *v1.ClaimTaskReq) (*v1.ClaimTaskResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ClaimTaskResp{}, nil
}

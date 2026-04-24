package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveUpdateCatInfoLogic {
	return &ApproveUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApproveUpdateCatInfoLogic) ApproveUpdateCatInfo(in *v1.ApproveUpdateCatInfoRequest) (*v1.ApproveUpdateCatInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApproveUpdateCatInfoResponse{}, nil
}

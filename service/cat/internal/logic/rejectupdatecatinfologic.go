package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectUpdateCatInfoLogic {
	return &RejectUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectUpdateCatInfoLogic) RejectUpdateCatInfo(in *v1.RejectUpdateCatInfoRequest) (*v1.RejectUpdateCatInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.RejectUpdateCatInfoResponse{}, nil
}

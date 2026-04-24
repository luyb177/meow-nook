package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUpdateCatApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUpdateCatApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUpdateCatApplyLogic {
	return &ListUpdateCatApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUpdateCatApplyLogic) ListUpdateCatApply(in *v1.ListUpdateCatApplyRequest) (*v1.ListUpdateCatApplyResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ListUpdateCatApplyResponse{}, nil
}

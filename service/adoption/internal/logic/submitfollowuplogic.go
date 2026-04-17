package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFollowUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitFollowUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFollowUpLogic {
	return &SubmitFollowUpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitFollowUpLogic) SubmitFollowUp(in *v1.SubmitFollowUpReq) (*v1.SubmitFollowUpResp, error) {
	// todo: add your logic here and delete this line

	return &v1.SubmitFollowUpResp{}, nil
}

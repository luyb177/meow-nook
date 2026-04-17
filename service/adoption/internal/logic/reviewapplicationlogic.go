package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/adoption/internal/svc"
	"github.com/luyb177/meow-nook/service/adoption/pb/adoption/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReviewApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReviewApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReviewApplicationLogic {
	return &ReviewApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReviewApplicationLogic) ReviewApplication(in *v1.ReviewApplicationReq) (*v1.ReviewApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ReviewApplicationResp{}, nil
}

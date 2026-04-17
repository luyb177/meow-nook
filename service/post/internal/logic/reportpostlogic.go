package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/post/internal/svc"
	"github.com/luyb177/meow-nook/service/post/pb/post/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportPostLogic {
	return &ReportPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReportPostLogic) ReportPost(in *v1.ReportPostReq) (*v1.ReportPostResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ReportPostResp{}, nil
}

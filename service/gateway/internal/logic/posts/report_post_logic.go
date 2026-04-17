// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package posts

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportPostLogic {
	return &ReportPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportPostLogic) ReportPost(req *types.ReportPostReq) error {
	// todo: add your logic here and delete this line

	return nil
}

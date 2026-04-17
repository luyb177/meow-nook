package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPointLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointLogsLogic {
	return &ListPointLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPointLogsLogic) ListPointLogs(in *v1.ListPointLogsReq) (*v1.ListPointLogsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.ListPointLogsResp{}, nil
}

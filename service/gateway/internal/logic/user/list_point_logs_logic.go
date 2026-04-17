// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPointLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPointLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPointLogsLogic {
	return &ListPointLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPointLogsLogic) ListPointLogs() (resp *types.PointLogsResp, err error) {
	// todo: add your logic here and delete this line

	return
}

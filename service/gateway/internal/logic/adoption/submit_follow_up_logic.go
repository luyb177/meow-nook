// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFollowUpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFollowUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFollowUpLogic {
	return &SubmitFollowUpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitFollowUpLogic) SubmitFollowUp(req *types.SubmitFollowUpReq) error {
	// todo: add your logic here and delete this line

	return nil
}

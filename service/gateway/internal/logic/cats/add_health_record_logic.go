// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package cats

import (
	"context"

	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHealthRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHealthRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHealthRecordLogic {
	return &AddHealthRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddHealthRecordLogic) AddHealthRecord(req *types.AddHealthRecordReq) (resp *types.HealthRecord, err error) {
	// todo: add your logic here and delete this line

	return
}

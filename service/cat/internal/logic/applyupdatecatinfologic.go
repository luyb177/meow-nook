package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyUpdateCatInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyUpdateCatInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyUpdateCatInfoLogic {
	return &ApplyUpdateCatInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 申请修改猫咪档案信息
func (l *ApplyUpdateCatInfoLogic) ApplyUpdateCatInfo(in *v1.ApplyUpdateCatInfoRequest) (*v1.ApplyUpdateCatInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.ApplyUpdateCatInfoResponse{}, nil
}

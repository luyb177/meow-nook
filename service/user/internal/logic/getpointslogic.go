package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/user/internal/svc"
	"github.com/luyb177/meow-nook/service/user/pb/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPointsLogic {
	return &GetPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Points
func (l *GetPointsLogic) GetPoints(in *v1.GetPointsReq) (*v1.GetPointsResp, error) {
	// todo: add your logic here and delete this line

	return &v1.GetPointsResp{}, nil
}

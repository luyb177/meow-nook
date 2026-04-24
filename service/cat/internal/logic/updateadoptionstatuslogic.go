package logic

import (
	"context"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdoptionStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAdoptionStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdoptionStatusLogic {
	return &UpdateAdoptionStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAdoptionStatusLogic) UpdateAdoptionStatus(in *v1.UpdateAdoptionStatusRequest) (*v1.UpdateAdoptionStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.UpdateAdoptionStatusResponse{}, nil
}

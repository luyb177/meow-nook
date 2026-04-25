package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
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
	// TODO(permission): check that in.OperatorId has admin role

	if _, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId); err != nil {
		return nil, errorx.ErrCatNotFound
	}

	if err := l.svcCtx.Repo.Cat.UpdateCat(l.ctx, in.CatId, map[string]interface{}{
		"adoption_status": in.AdoptionStatus,
	}); err != nil {
		logger.Error("UpdateAdoptionStatus: update cat failed")
		return nil, errorx.WrapInternal("更新领养状态失败", err)
	}

	if err := l.svcCtx.Repo.Cat.UpsertCatAdoption(l.ctx, in.CatId, in.AdoptionStatus); err != nil {
		logger.Error("UpdateAdoptionStatus: upsert adoption failed")
		return nil, errorx.WrapInternal("更新领养记录失败", err)
	}

	return &v1.UpdateAdoptionStatusResponse{Success: true}, nil
}

package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"gorm.io/gorm"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type UpdateReturnStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReturnStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReturnStatusLogic {
	return &UpdateReturnStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReturnStatusLogic) UpdateReturnStatus(in *v1.UpdateReturnStatusRequest) (*v1.UpdateReturnStatusResponse, error) {
	if in.AdoptionId == 0 || in.OperatorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id and operator_id is required", errorx.ErrBadRequest)
	}

	returnedAt := time.Now()
	if in.ReturnedAt != nil {
		returnedAt = in.ReturnedAt.AsTime()
	}

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		if err := l.svcCtx.Repo.Adoption.UpdateReturnStatus(
			l.ctx,
			in.AdoptionId,
			in.Returned,
			in.ReturnedToUserId,
			in.ReturnReason,
			in.Photos,
			returnedAt,
			tx,
		); err != nil {
			return errorx.WrapDBUpdate("更新退回状态失败", err)
		}

		adoptInfo, err := l.svcCtx.Repo.Adoption.GetAdoptionByID(l.ctx, in.AdoptionId, tx)
		if err != nil {
			return errorx.WrapDBQuery("查询领养记录失败", err)
		}

		if in.Returned {
			if err := l.svcCtx.Repo.Cat.UpdateCatAdoptionStatus(l.ctx, adoptInfo.CatID, cat.CatAdoptionStatusPending, 0, tx); err != nil {
				return errorx.WrapDBUpdate("更新猫咪领养状态失败", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateReturnStatusResponse{
		Message: "更新退回状态成功",
	}, nil
}

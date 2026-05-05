package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"gorm.io/gorm"

	adoptionRepo "github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type DirectAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDirectAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DirectAdoptLogic {
	return &DirectAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DirectAdoptLogic) DirectAdopt(in *v1.DirectAdoptRequest) (*v1.DirectAdoptResponse, error) {
	if in.CatId == 0 || in.AdopterId == 0 || in.CreatorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id, adopter_id and creator_id are required", errorx.ErrBadRequest)
	}

	var adoptionID uint64

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		now := time.Now()
		record := &adoptionRepo.Adoption{
			CatID:       in.CatId,
			AdopterID:   in.AdopterId,
			CreatorID:   in.CreatorId,
			ApplyID:     0,
			Status:      adoptionRepo.AdoptionStatusActive,
			AgreementNo: in.AgreementNo,
			AgreedAt:    &now,
			AdoptedAt:   &now,
			Note:        in.Note,
		}

		if err := l.svcCtx.Repo.Adoption.CreateAdoption(l.ctx, record, tx); err != nil {
			return errorx.WrapDBInsert("创建领养记录失败", err)
		}
		adoptionID = record.ID

		if err := l.svcCtx.Repo.Cat.UpdateCatAdoptionStatus(l.ctx, in.CatId, cat.CatAdoptionStatusAdopted, in.AdopterId, tx); err != nil {
			return errorx.WrapDBUpdate("更新猫咪领养状态失败", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.DirectAdoptResponse{
		AdoptionId: adoptionID,
		Message:    "直接领养成功",
	}, nil
}

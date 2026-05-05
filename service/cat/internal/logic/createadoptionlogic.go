package logic

import (
	"context"
	"time"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	adoptionRepo "github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type CreateAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdoptionLogic {
	return &CreateAdoptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdoptionLogic) CreateAdoption(in *v1.CreateAdoptionRequest) (*v1.CreateAdoptionResponse, error) {
	if in.ApplyId == 0 || in.CreatorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id and creator_id are required", errorx.ErrBadRequest)
	}

	var adoptionID uint64

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *gorm.DB) error {
		apply, err := l.svcCtx.Repo.Adoption.GetApplyByID(l.ctx, in.ApplyId, tx)
		if err != nil {
			return err
		}

		if apply.Status != adoptionRepo.AdoptApplyStatusApproved {
			return errorx.Wrap(errorx.CodeBadRequest, "领养申请未审核通过，无法创建领养记录", errorx.ErrBadRequest)
		}

		now := time.Now()
		record := &adoptionRepo.Adoption{
			CatID:       apply.CatID,
			AdopterID:   apply.ApplicantID,
			ApplyID:     apply.ID,
			CreatorID:   in.CreatorId,
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

		if err := l.svcCtx.Repo.Adoption.UpdateStatus(l.ctx, apply.ID, adoptionRepo.AdoptApplyStatusCompleted, map[string]any{
			"adoption_id": adoptionID,
		}, tx); err != nil {
			return errorx.WrapDBUpdate("更新领养申请状态失败", err)
		}

		if err := l.svcCtx.Repo.Cat.UpdateCatAdoptionStatus(l.ctx, apply.CatID, cat.CatAdoptionStatusAdopted, apply.ApplicantID, tx); err != nil {
			return errorx.WrapDBUpdate("更新猫咪领养状态失败", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAdoptionResponse{
		AdoptionId: adoptionID,
		Message:    "创建领养记录成功",
	}, nil
}

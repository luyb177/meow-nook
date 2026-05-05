package logic

import (
	"context"
	"errors"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"gorm.io/gorm"

	"github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type ApplyAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAdoptLogic {
	return &ApplyAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyAdoptLogic) ApplyAdopt(in *v1.ApplyAdoptRequest) (*v1.ApplyAdoptResponse, error) {
	if in.ApplicantId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "applicant_is is required", errorx.ErrBadRequest)
	}
	if in.CatId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id is required", errorx.ErrBadRequest)
	}

	// 1. 检查猫咪是否存在
	catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, in.CatId)
	if err != nil {
		return nil, err
	}
	if catInfo == nil {
		return nil, errorx.WrapDBQuery("cat not found", errorx.ErrNotFound)
	}

	// 2. 检查猫咪是否可领养
	if catInfo.AdoptionStatus != cat.CatAdoptionStatusPending {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "猫咪不可领养", adoption.ErrCatNotAvailable)
	}

	// 3. 检查是否重复申请
	_, err = l.svcCtx.Repo.Adoption.GetApplyByCatAndUser(l.ctx, in.CatId, in.ApplicantId, []string{
		adoption.AdoptApplyStatusRejected,
		adoption.AdoptApplyStatusCancelled,
		adoption.AdoptApplyStatusExpired,
	})
	if err == nil {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "重复申请领养该小猫", adoption.ErrDuplicateApply)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.WrapDBQuery("查询领养申请失败", err)
	}

	// 4. TODO: 查询用户积分
	// 当前先占位，后续接 user service
	applicantCreditScore := int32(0)

	apply := &adoption.AdoptApplication{
		CatID:                in.CatId,
		ApplicantID:          in.ApplicantId,
		ApplyReason:          in.ApplyReason,
		ContactPhone:         in.ContactPhone,
		ContactWechat:        in.ContactWechat,
		ApplicantCreditScore: int(applicantCreditScore),
		Status:               adoption.AdoptApplyStatusPending,
	}

	if err = l.svcCtx.Repo.Adoption.CreateApply(l.ctx, apply); err != nil {
		return nil, err
	}

	return &v1.ApplyAdoptResponse{
		ApplyId:              apply.ID,
		Status:               apply.Status,
		Message:              "申请成功，请等待审核",
		ApplicantCreditScore: applicantCreditScore,
	}, nil
}

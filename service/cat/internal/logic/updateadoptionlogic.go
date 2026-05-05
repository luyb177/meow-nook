package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type UpdateAdoptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdoptionLogic {
	return &UpdateAdoptionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdoptionLogic) UpdateAdoption(in *v1.UpdateAdoptionRequest) (*v1.UpdateAdoptionResponse, error) {
	if in.AdoptionId == 0 || in.OperatorId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id and operator_id are required", errorx.ErrBadRequest)
	}

	values := map[string]any{}
	if in.AgreementNo != "" {
		values["agreement_no"] = in.AgreementNo
	}
	if in.Note != "" {
		values["note"] = in.Note
	}

	if len(values) == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "no fields to update", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.UpdateAdoption(l.ctx, in.AdoptionId, values); err != nil {
		return nil, errorx.WrapDBUpdate("更新领养记录失败", err)
	}

	return &v1.UpdateAdoptionResponse{
		Message: "更新成功",
	}, nil
}

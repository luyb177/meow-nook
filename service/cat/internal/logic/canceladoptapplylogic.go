package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type CancelAdoptApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelAdoptApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAdoptApplyLogic {
	return &CancelAdoptApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAdoptApplyLogic) CancelAdoptApply(in *v1.CancelAdoptApplyRequest) (*v1.CancelAdoptApplyResponse, error) {
	if in.ApplyId == 0 || in.ApplicantId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id and applicant_id are required", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.CancelApply(l.ctx, in.ApplyId, in.ApplicantId, in.Reason); err != nil {
		return nil, errorx.WrapDBUpdate("取消申请失败", err)
	}

	return &v1.CancelAdoptApplyResponse{
		Message: "取消申请成功",
	}, nil
}

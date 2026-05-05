package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
)

type RejectAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectAdoptLogic {
	return &RejectAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectAdoptLogic) RejectAdopt(in *v1.RejectAdoptRequest) (*v1.RejectAdoptResponse, error) {
	if in.ApplyId == 0 || in.ReviewerId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id and reviewer_id are required", errorx.ErrBadRequest)
	}
	if in.RejectReason == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "reject_reason is required", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.RejectApply(l.ctx, in.ApplyId, in.ReviewerId, in.RejectReason); err != nil {
		return nil, errorx.WrapDBUpdate("驳回申请失败", err)
	}

	return &v1.RejectAdoptResponse{
		Message: "驳回成功",
	}, nil
}

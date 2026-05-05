package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type ApproveAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveAdoptLogic {
	return &ApproveAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveAdoptLogic) ApproveAdopt(in *v1.ApproveAdoptRequest) (*v1.ApproveAdoptResponse, error) {
	if in.ApplyId == 0 || in.ReviewerId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id and reviewer_id are required", errorx.ErrBadRequest)
	}

	if err := l.svcCtx.Repo.Adoption.ApproveApply(l.ctx, in.ApplyId, in.ReviewerId, 7); err != nil {
		return nil, errorx.WrapDBUpdate("同意申请失败", err)
	}

	apply, err := l.svcCtx.Repo.Adoption.GetApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询领养申请失败", err)
	}

	resp := &v1.ApproveAdoptResponse{
		Message: "审核通过",
	}

	if apply.ExpiresAt != nil {
		resp.ExpiresAt = timestamppb.New(*apply.ExpiresAt)
	}

	return resp, nil
}

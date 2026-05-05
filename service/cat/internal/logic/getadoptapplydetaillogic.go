package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type GetAdoptApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptApplyDetailLogic {
	return &GetAdoptApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdoptApplyDetailLogic) GetAdoptApplyDetail(in *v1.GetAdoptApplyDetailRequest) (*v1.GetAdoptApplyDetailResponse, error) {
	if in.ApplyId == 0 || in.RequesterId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id and requester_id are required", errorx.ErrBadRequest)
	}

	apply, err := l.svcCtx.Repo.Adoption.GetApplyByID(l.ctx, in.ApplyId)
	if err != nil {
		return nil, err
	}

	// 基础权限：申请人本人可看
	// TODO: 后续补管理员权限校验
	if apply.ApplicantID != in.RequesterId {
		return nil, errorx.Wrap(errorx.CodePermissionDenied, "permission denied", errorx.ErrForbidden)
	}

	resp := &v1.GetAdoptApplyDetailResponse{
		Apply: &v1.AdoptApplicationVO{
			Id:                   apply.ID,
			CatId:                apply.CatID,
			ApplicantId:          apply.ApplicantID,
			ApplyReason:          apply.ApplyReason,
			ContactPhone:         apply.ContactPhone,
			ContactWechat:        apply.ContactWechat,
			ApplicantCreditScore: int32(apply.ApplicantCreditScore),
			Status:               apply.Status,
			RejectReason:         apply.RejectReason,
			ReviewerId:           apply.ReviewerID,
			AdoptionId:           apply.AdoptionID,
			CreatedAt:            timestamppb.New(apply.CreatedAt),
			UpdatedAt:            timestamppb.New(apply.UpdatedAt),
		},
	}

	if apply.ReviewedAt != nil {
		resp.Apply.ReviewedAt = timestamppb.New(*apply.ReviewedAt)
	}
	if apply.ApprovedAt != nil {
		resp.Apply.ApprovedAt = timestamppb.New(*apply.ApprovedAt)
	}
	if apply.ExpiresAt != nil {
		resp.Apply.ExpiresAt = timestamppb.New(*apply.ExpiresAt)
	}

	catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, apply.CatID)
	if err == nil && catInfo != nil {
		resp.Apply.CatName = catInfo.Name
		resp.Apply.CatAvatar = catInfo.Avatar
		resp.Cat = &v1.CatBriefVO{
			Id:     catInfo.ID,
			Name:   catInfo.Name,
			Avatar: catInfo.Avatar,
			Gender: catInfo.Gender,
		}
	}

	return resp, nil
}

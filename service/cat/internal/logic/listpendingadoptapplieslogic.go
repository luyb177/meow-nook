package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	adoptionRepo "github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type ListPendingAdoptAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingAdoptAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingAdoptAppliesLogic {
	return &ListPendingAdoptAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPendingAdoptAppliesLogic) ListPendingAdoptApplies(in *v1.ListPendingAdoptAppliesRequest) (*v1.ListPendingAdoptAppliesResponse, error) {
	filter := &adoptionRepo.AdoptApplicationFilter{
		CatID:    in.CatId,
		Status:   adoptionRepo.AdoptApplyStatusPending,
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Adoption.ListApplies(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询待审核领养申请列表失败", err)
	}

	items := make([]*v1.PendingApplyVO, 0, len(list))
	for _, v := range list {
		item := &v1.PendingApplyVO{
			ApplyId:              v.ID,
			CatId:                v.CatID,
			ApplicantId:          v.ApplicantID,
			ApplicantCreditScore: int32(v.ApplicantCreditScore),
			ApplyReason:          v.ApplyReason,
			ContactPhone:         v.ContactPhone,
			ContactWechat:        v.ContactWechat,
			CreatedAt:            timestamppb.New(v.CreatedAt),
		}

		catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, v.CatID)
		if err == nil && catInfo != nil {
			item.CatName = catInfo.Name
			item.CatAvatar = catInfo.Avatar
			item.CatGender = catInfo.Gender
		}

		items = append(items, item)
	}

	return &v1.ListPendingAdoptAppliesResponse{
		Items:    items,
		Total:    int32(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

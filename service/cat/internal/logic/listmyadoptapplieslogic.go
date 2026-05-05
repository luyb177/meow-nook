package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	adoptionRepo "github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type ListMyAdoptAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyAdoptAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyAdoptAppliesLogic {
	return &ListMyAdoptAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyAdoptAppliesLogic) ListMyAdoptApplies(in *v1.ListMyAdoptAppliesRequest) (*v1.ListMyAdoptAppliesResponse, error) {
	filter := &adoptionRepo.AdoptApplicationFilter{
		ApplicantID: in.ApplicantId,
		Status:      in.Status,
		Page:        int(in.Page),
		PageSize:    int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Adoption.ListApplies(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询领养申请列表失败", err)
	}

	items := make([]*v1.AdoptApplyItemVO, 0, len(list))
	for _, v := range list {
		item := &v1.AdoptApplyItemVO{
			ApplyId:      v.ID,
			CatId:        v.CatID,
			ApplyReason:  v.ApplyReason,
			Status:       v.Status,
			RejectReason: v.RejectReason,
			CreatedAt:    timestamppb.New(v.CreatedAt),
		}
		if v.ExpiresAt != nil {
			item.ExpiresAt = timestamppb.New(*v.ExpiresAt)
		}

		catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, v.CatID)
		if err == nil && catInfo != nil {
			item.CatName = catInfo.Name
			item.CatAvatar = catInfo.Avatar
			item.CatGender = catInfo.Gender
		}

		items = append(items, item)
	}

	return &v1.ListMyAdoptAppliesResponse{
		Items:    items,
		Total:    int32(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

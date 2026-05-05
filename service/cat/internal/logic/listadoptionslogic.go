package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	adoptionRepo "github.com/luyb177/meow-nook/service/cat/internal/repo/adoption"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type ListAdoptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAdoptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdoptionsLogic {
	return &ListAdoptionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAdoptionsLogic) ListAdoptions(in *v1.ListAdoptionsRequest) (*v1.ListAdoptionsResponse, error) {
	filter := &adoptionRepo.AdoptionFilter{
		CatID:     in.CatId,
		AdopterID: in.AdopterId,
		Status:    in.Status,
		Page:      int(in.Page),
		PageSize:  int(in.PageSize),
	}

	list, total, err := l.svcCtx.Repo.Adoption.ListAdoptions(l.ctx, filter)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询领养列表失败", err)
	}

	items := make([]*v1.AdoptionItemVO, 0, len(list))
	for _, v := range list {
		item := &v1.AdoptionItemVO{
			AdoptionId:  v.ID,
			CatId:       v.CatID,
			AdopterId:   v.AdopterID,
			Status:      v.Status,
			AgreementNo: v.AgreementNo,
			IsReturned:  v.IsReturned,
		}
		if v.AdoptedAt != nil {
			item.AdoptedAt = timestamppb.New(*v.AdoptedAt)
		}

		catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, v.CatID)
		if err == nil && catInfo != nil {
			item.CatName = catInfo.Name
			item.CatAvatar = catInfo.Avatar
		}

		items = append(items, item)
	}

	return &v1.ListAdoptionsResponse{
		Items:    items,
		Total:    int32(total),
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}

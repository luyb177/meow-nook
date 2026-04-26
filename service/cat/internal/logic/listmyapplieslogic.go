// internal/logic/list_my_applies_logic.go
package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListMyAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyAppliesLogic {
	return &ListMyAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyAppliesLogic) ListMyApplies(in *v1.ListMyAppliesRequest) (*v1.ListMyAppliesResponse, error) {
	var status *string
	if in.Status != "" {
		status = &in.Status
	}

	list, total, err := l.svcCtx.Repo.Cat.ListApplies(l.ctx, cat.ApplyListFilter{
		ApplicantUserID: in.ApplicantUserId,
		Status:          status,
		Page:            int(in.Page),
		PageSize:        int(in.PageSize),
	})

	if err != nil {
		return nil, errorx.WrapDBQuery("查询申请列表失败", err)
	}

	items := make([]*v1.ListMyAppliesResponse_ApplyItem, 0, len(list))
	for _, apply := range list {
		items = append(items, &v1.ListMyAppliesResponse_ApplyItem{
			ApplyId:   apply.ID,
			Name:      apply.Name,
			Status:    apply.Status,
			CatId:     apply.CatID,
			CreatedAt: timestamppb.New(apply.CreatedAt),
		})
	}

	return &v1.ListMyAppliesResponse{
		List:  items,
		Total: total,
	}, nil
}

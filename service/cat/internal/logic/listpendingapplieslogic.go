// internal/logic/list_pending_applies_logic.go
package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/repo/cat"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListPendingAppliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPendingAppliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPendingAppliesLogic {
	return &ListPendingAppliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPendingAppliesLogic) ListPendingApplies(in *v1.ListPendingAppliesRequest) (*v1.ListPendingAppliesResponse, error) {
	status := cat.ApplyStatusPending

	list, total, err := l.svcCtx.Repo.Cat.ListApplies(l.ctx, cat.ApplyListFilter{
		Status:   &status,
		Keyword:  in.Keyword,
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	})

	if err != nil {
		return nil, errorx.WrapDBQuery("查询待审核申请列表失败", err)
	}

	items := make([]*v1.ListPendingAppliesResponse_ApplyItem, 0, len(list))
	for _, apply := range list {
		items = append(items, &v1.ListPendingAppliesResponse_ApplyItem{
			ApplyId:          apply.ID,
			Name:             apply.Name,
			Gender:           apply.Gender,
			AgeStage:         apply.AgeStage,
			DiscoveryAddress: apply.DiscoveryAddress,
			ApplicantUserId:  apply.ApplicantUserID,
			CreatedAt:        timestamppb.New(apply.CreatedAt),
		})
	}

	return &v1.ListPendingAppliesResponse{
		List:  items,
		Total: total,
	}, nil
}

package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	"github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatsLogic {
	return &ListCatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 猫咪列表
func (l *ListCatsLogic) ListCats(in *v1.ListCatsRequest) (*v1.ListCatsResponse, error) {
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	lastID := decodeCursor(in.Cursor)

	cats, hasMore, err := l.svcCtx.Repo.Cat.ListCats(l.ctx, lastID, pageSize)
	if err != nil {
		logger.Error("ListCats: query failed")
		return nil, errorx.WrapInternal("查询猫咪列表失败", err)
	}

	items := make([]*v1.CatItem, 0, len(cats))
	for _, c := range cats {
		items = append(items, &v1.CatItem{
			CatId:          c.ID,
			CatCode:        c.CatCode,
			Name:           c.Name,
			Gender:         c.Gender,
			BodySize:       c.BodySize,
			AgeStage:       c.AgeStage,
			AdoptionStatus: c.AdoptionStatus,
			Avatar:         c.Avatar,
			CreatedAt:      timestamppb.New(c.CreatedAt),
		})
	}

	nextCursor := ""
	if hasMore && len(cats) > 0 {
		nextCursor = encodeCursor(cats[len(cats)-1].ID)
	}

	return &v1.ListCatsResponse{
		List:       items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

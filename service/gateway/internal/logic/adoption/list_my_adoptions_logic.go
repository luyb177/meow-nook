package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	. "github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ListMyAdoptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyAdoptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyAdoptionsLogic {
	return &ListMyAdoptionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyAdoptionsLogic) ListMyAdoptions(req *types.ListMyAdoptionsReq) (*types.ListAdoptionsResp, error) {
	logger.Info("ListMyAdoptionsLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}
	// todo get userID from token
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.ListMyAdoptions(l.ctx, &catpb.ListMyAdoptionsRequest{
		UserId:   userID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	items := make([]types.AdoptionItemVO, 0, len(resp.Items))
	for _, v := range resp.Items {
		items = append(items, types.AdoptionItemVO{
			AdoptionId:  v.AdoptionId,
			CatId:       v.CatId,
			CatName:     v.CatName,
			CatAvatar:   v.CatAvatar,
			AdopterId:   v.AdopterId,
			AdopterName: v.AdopterName,
			Status:      v.Status,
			AgreementNo: v.AgreementNo,
			AdoptedAt:   PBTimeToString(v.AdoptedAt),
			IsReturned:  v.IsReturned,
		})
	}

	return &types.ListAdoptionsResp{
		Items:    items,
		Total:    int64(resp.Total),
		Page:     resp.Page,
		PageSize: resp.PageSize,
	}, nil
}

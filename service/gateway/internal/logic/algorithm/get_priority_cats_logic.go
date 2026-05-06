package algorithm

import (
	"context"

	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPriorityCatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPriorityCatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPriorityCatsLogic {
	return &GetPriorityCatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPriorityCatsLogic) GetPriorityCats(req *types.GetPriorityCatsReq) (*types.GetPriorityCatsResp, error) {
	// 设置默认值
	topN := req.TopN
	if topN <= 0 {
		topN = 10
	}

	// 调用 cat 服务
	resp, err := l.svcCtx.CatRPC.GetPriorityCats(l.ctx, &catpb.GetPriorityCatsRequest{
		TopN: int32(topN),
	})
	if err != nil {
		return nil, err
	}

	// 转换响应
	cats := make([]types.PriorityCatInfo, len(resp.Cats))
	for i, c := range resp.Cats {
		cats[i] = types.PriorityCatInfo{
			CatID:     c.CatId,
			CatName:   c.CatName,
			CatCode:   c.CatCode,
			Score:     c.Score,
			Details:   c.Details,
			Longitude: c.Longitude,
			Latitude:  c.Latitude,
		}
	}

	return &types.GetPriorityCatsResp{
		Cats:  cats,
		Total: int(resp.Total),
	}, nil
}

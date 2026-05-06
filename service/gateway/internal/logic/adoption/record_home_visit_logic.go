// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordHomeVisitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordHomeVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordHomeVisitLogic {
	return &RecordHomeVisitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RecordHomeVisit
func (l *RecordHomeVisitLogic) RecordHomeVisit(req *types.RecordHomeVisitReq) (*types.RecordHomeVisitResp, error) {
	logger.Info("RecordHomeVisitLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.RecordHomeVisit(l.ctx, &catpb.RecordHomeVisitRequest{
		AdoptionId: req.AdoptionId,
		Remark:     req.Remark,
		Photos:     req.Photos,
		VisitorId:  uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RecordHomeVisitResp{Message: resp.Message}, nil
}

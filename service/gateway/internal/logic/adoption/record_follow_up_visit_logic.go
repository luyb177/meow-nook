// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordFollowUpVisitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordFollowUpVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordFollowUpVisitLogic {
	return &RecordFollowUpVisitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RecordFollowUpVisit
func (l *RecordFollowUpVisitLogic) RecordFollowUpVisit(req *types.RecordFollowUpVisitReq) (*types.RecordFollowUpVisitResp, error) {
	logger.Info("RecordFollowUpVisitLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	resp, err := l.svcCtx.CatRPC.RecordFollowUpVisit(l.ctx, &catpb.RecordFollowUpVisitRequest{
		AdoptionId: req.AdoptionId,
		VisitType:  req.VisitType,
		Remark:     req.Remark,
		Photos:     req.Photos,
		VisitorId:  userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RecordFollowUpVisitResp{Message: resp.Message}, nil
}

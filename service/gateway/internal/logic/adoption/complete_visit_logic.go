package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type CompleteVisitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteVisitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteVisitLogic {
	return &CompleteVisitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteVisitLogic) CompleteVisit(req *types.RecordFollowUpVisitReq) (*types.RecordFollowUpVisitResp, error) {
	logger.Info("CompleteVisitLogic called")
	//
	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	if req.AdoptionId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id is required", errorx.ErrBadRequest)
	}
	if req.VisitType < 1 || req.VisitType > 4 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "invalid visit_type", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.CompleteVisit(l.ctx, &catpb.RecordFollowUpVisitRequest{
		AdoptionId: req.AdoptionId,
		VisitType:  req.VisitType,
		Remark:     req.Remark,
		Photos:     req.Photos,
		VisitorId:  userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RecordFollowUpVisitResp{
		Message: resp.Message,
	}, nil
}

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/logic"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type RejectAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectAdoptLogic {
	return &RejectAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectAdoptLogic) RejectAdopt(req *types.RejectAdoptReq) (*types.RejectAdoptResp, error) {
	logger.Info("RejectAdoptLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	if req.ApplyId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id is required", errorx.ErrBadRequest)
	}
	if req.RejectReason == "" {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "拒绝原因不能为空", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.RejectAdopt(l.ctx, &catpb.RejectAdoptRequest{
		ApplyId:      req.ApplyId,
		RejectReason: req.RejectReason,
		ReviewerId:   uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.RejectAdoptResp{
		Message: resp.Message,
	}, nil
}

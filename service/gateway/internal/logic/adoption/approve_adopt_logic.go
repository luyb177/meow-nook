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

type ApproveAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveAdoptLogic {
	return &ApproveAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveAdoptLogic) ApproveAdopt(req *types.ApproveAdoptReq) (*types.ApproveAdoptResp, error) {
	logger.Info("ApproveAdoptLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	if req.ApplyId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.ApproveAdopt(l.ctx, &catpb.ApproveAdoptRequest{
		ApplyId:    req.ApplyId,
		Note:       req.Note,
		ReviewerId: userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApproveAdoptResp{
		ApplyId:   resp.ApplyId,
		Status:    resp.Status,
		Message:   resp.Message,
		ExpiresAt: PBTimeToString(resp.ExpiresAt),
	}, nil
}

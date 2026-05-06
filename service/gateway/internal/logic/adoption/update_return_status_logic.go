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

type UpdateReturnStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReturnStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReturnStatusLogic {
	return &UpdateReturnStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateReturnStatus
func (l *UpdateReturnStatusLogic) UpdateReturnStatus(req *types.UpdateReturnStatusReq) (*types.UpdateReturnStatusResp, error) {
	logger.Info("UpdateReturnStatusLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	resp, err := l.svcCtx.CatRPC.UpdateReturnStatus(l.ctx, &catpb.UpdateReturnStatusRequest{
		AdoptionId:       req.AdoptionId,
		Returned:         req.Returned,
		ReturnedToUserId: req.ReturnedToUserId,
		ReturnReason:     req.ReturnReason,
		Photos:           req.Photos,
		OperatorId:       uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.UpdateReturnStatusResp{Message: resp.Message}, nil
}

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

type UpdateAdoptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdoptionLogic {
	return &UpdateAdoptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdoptionLogic) UpdateAdoption(req *types.UpdateAdoptionReq) (*types.UpdateAdoptionResp, error) {
	logger.Info("UpdateAdoptionLogic called")

	userID, err := logic.GetUserID(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	}

	if req.AdoptionId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.UpdateAdoption(l.ctx, &catpb.UpdateAdoptionRequest{
		AdoptionId:  req.AdoptionId,
		AgreementNo: req.AgreementNo,
		Note:        req.Note,
		OperatorId:  uint64(userID),
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.UpdateAdoptionResp{
		Message: resp.Message,
	}, nil
}

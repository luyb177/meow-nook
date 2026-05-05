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

type CreateAdoptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAdoptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdoptionLogic {
	return &CreateAdoptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateAdoption
func (l *CreateAdoptionLogic) CreateAdoption(req *types.CreateAdoptionReq) (*types.CreateAdoptionResp, error) {
	logger.Info("CreateAdoptionLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}

	// todo get userID from token
	userID := uint64(1)

	if req.ApplyId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "apply_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.CreateAdoption(l.ctx, &catpb.CreateAdoptionRequest{
		ApplyId:     req.ApplyId,
		AgreementNo: req.AgreementNo,
		Note:        req.Note,
		CreatorId:   userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.CreateAdoptionResp{
		AdoptionId: resp.AdoptionId,
		Message:    resp.Message,
	}, nil
}

package adoption

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/logger"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

type ApplyAdoptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyAdoptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyAdoptLogic {
	return &ApplyAdoptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyAdoptLogic) ApplyAdopt(req *types.ApplyAdoptReq) (*types.ApplyAdoptResp, error) {
	logger.Info("ApplyAdoptLogic called")

	//userID, err := ctxutil.GetUserID(l.ctx)
	//if err != nil {
	//	return nil, errorx.Wrap(errorx.CodeUnauthorized, "未登录", err)
	//}
	// todo get userID from token
	userID := uint64(1)

	if req.CatId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "cat_id is required", errorx.ErrBadRequest)
	}

	resp, err := l.svcCtx.CatRPC.ApplyAdopt(l.ctx, &catpb.ApplyAdoptRequest{
		CatId:         req.CatId,
		ApplyReason:   req.ApplyReason,
		ContactPhone:  req.ContactPhone,
		ContactWechat: req.ContactWechat,
		ApplicantId:   userID,
	})
	if err != nil {
		return nil, errorx.FromGRPC(err)
	}

	return &types.ApplyAdoptResp{
		ApplyId:              resp.ApplyId,
		Status:               resp.Status,
		Message:              resp.Message,
		ApplicantCreditScore: resp.ApplicantCreditScore,
	}, nil
}

package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	v1 "github.com/luyb177/meow-nook/service/cat/pb/cat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/luyb177/meow-nook/service/cat/internal/svc"
)

type GetAdoptionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdoptionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdoptionDetailLogic {
	return &GetAdoptionDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdoptionDetailLogic) GetAdoptionDetail(in *v1.GetAdoptionDetailRequest) (*v1.GetAdoptionDetailResponse, error) {
	if in.AdoptionId == 0 || in.RequesterId == 0 {
		return nil, errorx.Wrap(errorx.CodeBadRequest, "adoption_id and requester_id are required", errorx.ErrBadRequest)
	}

	info, err := l.svcCtx.Repo.Adoption.GetAdoptionByID(l.ctx, in.AdoptionId)
	if err != nil {
		return nil, errorx.WrapDBQuery("查询领养信息失败", err)
	}

	// 基础权限：领养人本人可看
	// TODO: 后续补管理员权限
	if info.AdopterID != in.RequesterId {
		return nil, errorx.Wrap(errorx.CodePermissionDenied, "领养人本人可看", errorx.ErrForbidden)
	}

	resp := &v1.GetAdoptionDetailResponse{
		Adoption: &v1.AdoptionVO{
			Id:              info.ID,
			CatId:           info.CatID,
			AdopterId:       info.AdopterID,
			Status:          info.Status,
			AgreementNo:     info.AgreementNo,
			IsReturned:      info.IsReturned,
			ReturnReason:    info.ReturnReason,
			Note:            info.Note,
			CreatedAt:       timestamppb.New(info.CreatedAt),
			UpdatedAt:       timestamppb.New(info.UpdatedAt),
			HomeVisitUserId: info.HomeVisitUserID,
		},
	}

	if info.AgreedAt != nil {
		resp.Adoption.AgreedAt = timestamppb.New(*info.AgreedAt)
	}
	if info.AdoptedAt != nil {
		resp.Adoption.AdoptedAt = timestamppb.New(*info.AdoptedAt)
	}
	if info.HomeVisitAt != nil {
		resp.Adoption.HomeVisitAt = timestamppb.New(*info.HomeVisitAt)
	}
	if info.VisitOneWeekAt != nil {
		resp.Adoption.VisitOneWeekAt = timestamppb.New(*info.VisitOneWeekAt)
	}
	if info.VisitOneMonthAt != nil {
		resp.Adoption.VisitOneMonthAt = timestamppb.New(*info.VisitOneMonthAt)
	}
	if info.VisitThreeMonthAt != nil {
		resp.Adoption.VisitThreeMonthAt = timestamppb.New(*info.VisitThreeMonthAt)
	}
	if info.VisitSixMonthAt != nil {
		resp.Adoption.VisitSixMonthAt = timestamppb.New(*info.VisitSixMonthAt)
	}
	if info.ReturnedAt != nil {
		resp.Adoption.ReturnedAt = timestamppb.New(*info.ReturnedAt)
	}

	catInfo, err := l.svcCtx.Repo.Cat.GetCatByID(l.ctx, info.CatID)
	if err == nil && catInfo != nil {
		resp.Adoption.CatName = catInfo.Name
		resp.Adoption.CatAvatar = catInfo.Avatar
		resp.Cat = &v1.CatBriefVO{
			Id:     catInfo.ID,
			Name:   catInfo.Name,
			Avatar: catInfo.Avatar,
			Gender: catInfo.Gender,
		}
	}

	resp.Adopter = &v1.UserBriefVO{
		Id: info.AdopterID,
	}

	return resp, nil
}

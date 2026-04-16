package logic

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/cat/internal/svc"
	catpb "github.com/luyb177/meow-nook/service/cat/pb/cat"
)

// ──────────────────────────────────────────────
// CreateCat
// ──────────────────────────────────────────────

type CreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatLogic {
	return &CreateCatLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateCatLogic) CreateCat(req *catpb.CreateCatReq) (*catpb.CreateCatResp, error) {
	if req.Name == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "猫咪名称不能为空")
	}
	// TODO: persist to DB
	return &catpb.CreateCatResp{CatId: 1}, nil
}

// ──────────────────────────────────────────────
// GetCat
// ──────────────────────────────────────────────

type GetCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatLogic {
	return &GetCatLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCatLogic) GetCat(req *catpb.GetCatReq) (*catpb.GetCatResp, error) {
	// TODO: fetch from DB
	return &catpb.GetCatResp{Cat: &catpb.CatInfo{Id: req.CatId}}, nil
}

// ──────────────────────────────────────────────
// UpdateCat
// ──────────────────────────────────────────────

type UpdateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatLogic {
	return &UpdateCatLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateCatLogic) UpdateCat(req *catpb.UpdateCatReq) error {
	// TODO: update in DB
	return nil
}

// ──────────────────────────────────────────────
// DeleteCat
// ──────────────────────────────────────────────

type DeleteCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCatLogic {
	return &DeleteCatLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteCatLogic) DeleteCat(req *catpb.DeleteCatReq) error {
	// TODO: delete from DB
	return nil
}

// ──────────────────────────────────────────────
// ListCats
// ──────────────────────────────────────────────

type ListCatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCatsLogic {
	return &ListCatsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListCatsLogic) ListCats(req *catpb.ListCatsReq) (*catpb.ListCatsResp, error) {
	// TODO: query DB with filters and pagination
	return &catpb.ListCatsResp{Cats: []*catpb.CatInfo{}, Total: 0}, nil
}

// ──────────────────────────────────────────────
// Rescue records
// ──────────────────────────────────────────────

type AddRescueRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRescueRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRescueRecordLogic {
	return &AddRescueRecordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddRescueRecordLogic) AddRescueRecord(req *catpb.AddRescueRecordReq) (*catpb.AddRescueRecordResp, error) {
	if req.Status == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "救助状态不能为空")
	}
	// TODO: persist to DB
	return &catpb.AddRescueRecordResp{RecordId: 1}, nil
}

type ListRescueRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRescueRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRescueRecordsLogic {
	return &ListRescueRecordsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListRescueRecordsLogic) ListRescueRecords(req *catpb.ListRescueRecordsReq) (*catpb.ListRescueRecordsResp, error) {
	// TODO: query from DB
	return &catpb.ListRescueRecordsResp{Records: []*catpb.RescueRecord{}}, nil
}

// ──────────────────────────────────────────────
// Health records
// ──────────────────────────────────────────────

type AddHealthRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHealthRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHealthRecordLogic {
	return &AddHealthRecordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddHealthRecordLogic) AddHealthRecord(req *catpb.AddHealthRecordReq) (*catpb.AddHealthRecordResp, error) {
	if req.RecordType == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "医疗记录类型不能为空")
	}
	// TODO: persist to DB
	return &catpb.AddHealthRecordResp{RecordId: 1}, nil
}

type ListHealthRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHealthRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHealthRecordsLogic {
	return &ListHealthRecordsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListHealthRecordsLogic) ListHealthRecords(req *catpb.ListHealthRecordsReq) (*catpb.ListHealthRecordsResp, error) {
	// TODO: query from DB
	return &catpb.ListHealthRecordsResp{Records: []*catpb.HealthRecord{}}, nil
}

// ──────────────────────────────────────────────
// Stats
// ──────────────────────────────────────────────

type GetStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatsLogic {
	return &GetStatsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetStatsLogic) GetStats(req *catpb.GetStatsReq) (*catpb.GetStatsResp, error) {
	// TODO: aggregate from DB
	return &catpb.GetStatsResp{}, nil
}

// ──────────────────────────────────────────────
// Heatmap
// ──────────────────────────────────────────────

type GetHeatmapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHeatmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHeatmapLogic {
	return &GetHeatmapLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetHeatmapLogic) GetHeatmap(req *catpb.GetHeatmapReq) (*catpb.GetHeatmapResp, error) {
	// TODO: aggregate geo data from DB
	return &catpb.GetHeatmapResp{Points: []*catpb.HeatmapPoint{}}, nil
}

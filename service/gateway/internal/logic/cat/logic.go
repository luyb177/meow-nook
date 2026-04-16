package cat

import (
	"context"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/service/gateway/internal/svc"
	"github.com/luyb177/meow-nook/service/gateway/internal/types"
)

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

func (l *ListCatsLogic) ListCats(req *types.ListCatsReq) (*types.ListCatsResp, error) {
	// TODO: call cat gRPC service
	return &types.ListCatsResp{Cats: []types.CatInfo{}, Total: 0}, nil
}

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

func (l *CreateCatLogic) CreateCat(req *types.CreateCatReq) (*types.CreateCatResp, error) {
	if req.Name == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "猫咪名称不能为空")
	}
	// TODO: call cat gRPC service
	return &types.CreateCatResp{CatId: 1}, nil
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

func (l *GetCatLogic) GetCat(catID int64) (*types.CatInfo, error) {
	// TODO: call cat gRPC service
	return &types.CatInfo{Id: catID}, nil
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

func (l *UpdateCatLogic) UpdateCat(catID int64, req *types.UpdateCatReq) error {
	// TODO: call cat gRPC service
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

func (l *DeleteCatLogic) DeleteCat(catID int64) error {
	// TODO: call cat gRPC service
	return nil
}

// ──────────────────────────────────────────────
// GetCatStats
// ──────────────────────────────────────────────

type GetCatStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCatStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatStatsLogic {
	return &GetCatStatsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCatStatsLogic) GetCatStats() (*types.CatStatsResp, error) {
	// TODO: call cat gRPC service
	return &types.CatStatsResp{}, nil
}

// ──────────────────────────────────────────────
// GetHeatmap
// ──────────────────────────────────────────────

type GetHeatmapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHeatmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHeatmapLogic {
	return &GetHeatmapLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetHeatmapLogic) GetHeatmap(req *types.HeatmapReq) (*types.HeatmapResp, error) {
	// TODO: call cat gRPC service
	return &types.HeatmapResp{Points: []types.HeatmapPoint{}}, nil
}

// ──────────────────────────────────────────────
// Rescue records
// ──────────────────────────────────────────────

type ListRescueRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRescueRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRescueRecordsLogic {
	return &ListRescueRecordsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListRescueRecordsLogic) ListRescueRecords(catID int64) ([]types.RescueRecord, error) {
	// TODO: call cat gRPC service
	return []types.RescueRecord{}, nil
}

type AddRescueRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddRescueRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRescueRecordLogic {
	return &AddRescueRecordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddRescueRecordLogic) AddRescueRecord(catID, rescuerID int64, req *types.AddRescueRecordReq) (*types.RescueRecord, error) {
	if req.Status == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "救助状态不能为空")
	}
	// TODO: call cat gRPC service
	return &types.RescueRecord{CatId: catID, RescuerId: rescuerID, Status: req.Status, Description: req.Description}, nil
}

// ──────────────────────────────────────────────
// Health records
// ──────────────────────────────────────────────

type ListHealthRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListHealthRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListHealthRecordsLogic {
	return &ListHealthRecordsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ListHealthRecordsLogic) ListHealthRecords(catID int64) ([]types.HealthRecord, error) {
	// TODO: call cat gRPC service
	return []types.HealthRecord{}, nil
}

type AddHealthRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHealthRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHealthRecordLogic {
	return &AddHealthRecordLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddHealthRecordLogic) AddHealthRecord(catID int64, req *types.AddHealthRecordReq) (*types.HealthRecord, error) {
	if req.RecordType == "" {
		return nil, errorx.New(errorx.CodeBadRequest, "医疗记录类型不能为空")
	}
	// TODO: call cat gRPC service
	return &types.HealthRecord{CatId: catID, RecordType: req.RecordType, Content: req.Content, Operator: req.Operator}, nil
}

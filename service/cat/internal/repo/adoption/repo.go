package adoption

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ==================== Repository 接口定义 ====================

type Repository interface {
	// ===================== 领养申请表 (AdoptApplication) 相关 =====================

	// CreateApply 创建新的领养申请
	CreateApply(ctx context.Context, apply *AdoptApplication, tx ...*gorm.DB) error

	// GetApplyByID 查询申请详情
	GetApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*AdoptApplication, error)

	// GetApplyByCatAndUser 检查用户是否已对该猫提交过申请（防重复）
	GetApplyByCatAndUser(ctx context.Context, catID, userID uint64, excludeStatuses []string, tx ...*gorm.DB) (*AdoptApplication, error)

	// ListApplies 列表查询（志愿者查看我的申请 / 管理员查看所有申请）
	ListApplies(ctx context.Context, filter *AdoptApplicationFilter, tx ...*gorm.DB) ([]*AdoptApplication, int64, error)

	// UpdateApply 更新申请信息（一般用于后台修改备注等，状态变更推荐专用方法）
	UpdateApply(ctx context.Context, applyID uint64, values map[string]any, tx ...*gorm.DB) error

	// CancelApply 志愿者取消申请 (Pending -> Cancelled)
	CancelApply(ctx context.Context, applyID uint64, applicantUserID uint64, reason string, tx ...*gorm.DB) error

	// RejectApply 管理员驳回申请 (Pending -> Rejected)
	RejectApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error

	// ApproveApply 管理员通过申请 (Pending -> Approved)，并设置过期时间
	ApproveApply(ctx context.Context, applyID uint64, reviewerID uint64, expiryDays int, tx ...*gorm.DB) error

	// MarkExpired 标记申请过期 (Approved -> Expired)，通常由定时任务调用或手动触发
	MarkExpired(ctx context.Context, applyID uint64, updatedAt time.Time, tx ...*gorm.DB) error

	// UpdateStatus 通用状态更新
	UpdateStatus(ctx context.Context, applyID uint64, status string, additionalData map[string]any, tx ...*gorm.DB) error

	// DeleteApply 软删除申请（管理员收回权限时可能用到，需谨慎）
	DeleteApply(ctx context.Context, applyID uint64, tx ...*gorm.DB) error

	// ===================== 领养记录表 (Adoption) 相关 =====================

	// CreateAdoption 创建正式领养记录（仅在申请完成转化后调用）
	CreateAdoption(ctx context.Context, adoption *Adoption, tx ...*gorm.DB) error

	// GetAdoptionByID 获取领养详情
	GetAdoptionByID(ctx context.Context, adoptionID uint64, tx ...*gorm.DB) (*Adoption, error)

	// GetAdoptionByCatID 通过猫咪 ID 查找当前有效的领养记录
	GetAdoptionByCatID(ctx context.Context, catID uint64, tx ...*gorm.DB) (*Adoption, error)

	// ListAdoptions 领养记录列表（管理员查进度 / 志愿者查自己养的猫）
	ListAdoptions(ctx context.Context, filter *AdoptionFilter, tx ...*gorm.DB) ([]*Adoption, int64, error)

	// UpdateAdoption 更新领养基础信息
	UpdateAdoption(ctx context.Context, adoptionID uint64, values map[string]any, tx ...*gorm.DB) error

	// RecordHomeVisit 记录家访信息
	RecordHomeVisit(ctx context.Context, adoptionID uint64, homeVisitAt time.Time, visitorID uint64, photos, remark string, tx ...*gorm.DB) error

	// RecordFollowUpVisit 记录回访信息
	// visitType: 1=1 周，2=1 月，3=3 月，4=6 月
	RecordFollowUpVisit(ctx context.Context, adoptionID uint64, visitType int, visitorID uint64, visitedAt time.Time, photos, remark string, tx ...*gorm.DB) error

	// UpdateReturnStatus 更新退回状态
	UpdateReturnStatus(ctx context.Context, adoptionID uint64, returned bool, returnedToUserID uint64, returnReason, photos string, returnedAt time.Time, tx ...*gorm.DB) error

	// DeleteAdoption 软删除领养记录
	DeleteAdoption(ctx context.Context, adoptionID uint64, tx ...*gorm.DB) error
}

// ==================== 具体实现结构 ====================

type repository struct {
	db     *gorm.DB
	client *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) Repository {
	return &repository{
		db:     db,
		client: client,
	}
}

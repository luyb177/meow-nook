package cat

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ==================== 实现 ====================

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

// ==================== Repository 接口 ====================

type Repository interface {
	// ===================== 申请相关 =====================

	// CreateApply 创建申请
	CreateApply(ctx context.Context, apply *CatCreateApply, tx ...*gorm.DB) error

	// GetApplyByID 查询申请详情
	GetApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatCreateApply, error)

	// GetApplyByIDForUpdate 查询申请详情（悲观锁，审核时用）
	GetApplyByIDForUpdate(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatCreateApply, error)

	// ListApplies 申请列表（支持按申请人、状态、关键字筛选）
	ListApplies(ctx context.Context, filter ApplyListFilter, tx ...*gorm.DB) ([]*CatCreateApply, int64, error)

	// UpdateApply 更新申请（只允许 pending 状态）
	UpdateApply(ctx context.Context, applyID uint64, values map[string]any, tx ...*gorm.DB) error

	// CancelApply 志愿者取消申请（pending -> canceled）
	CancelApply(ctx context.Context, applyID uint64, applicantUserID uint64, reason string, tx ...*gorm.DB) error

	// RejectApply 管理员驳回申请（pending -> rejected）
	RejectApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error

	// ApproveApply 管理员同意申请（pending -> approved）
	// catID 是审核通过后生成的正式猫咪档案 ID
	ApproveApply(ctx context.Context, applyID uint64, reviewerID uint64, catID uint64, tx ...*gorm.DB) error

	// DeleteApply 删除申请（软删除）
	DeleteApply(ctx context.Context, applyID uint64, tx ...*gorm.DB) error

	// ===================== 正式猫咪档案相关 =====================

	CreateCat(ctx context.Context, cat *Cat, tx ...*gorm.DB) error
	GetCatByID(ctx context.Context, catID uint64, tx ...*gorm.DB) (*Cat, error)
	GetCatByCode(ctx context.Context, catCode string, tx ...*gorm.DB) (*Cat, error)
	ListCats(ctx context.Context, filter CatListFilter, tx ...*gorm.DB) ([]*Cat, int64, error)
	UpdateCat(ctx context.Context, catID uint64, values map[string]any, tx ...*gorm.DB) error
	UpdateCatAdoptionStatus(ctx context.Context, catID uint64, status string, adopterID uint64, tx ...*gorm.DB) error
	DeleteCat(ctx context.Context, catID uint64, tx ...*gorm.DB) error

	ListCatsWithNearby(ctx context.Context, filter CatListFilter, near *NearFilter, tx ...*gorm.DB) ([]*Cat, []float64, int64, error)
}

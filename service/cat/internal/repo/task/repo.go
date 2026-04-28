package task

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ==================== Repository 接口定义 ====================

type Repository interface {
	// ===================== 任务申请（CatTaskApply）====================

	// CreateTaskApply 志愿者申请创建任务
	CreateTaskApply(ctx context.Context, apply *CatTaskApply, tx ...*gorm.DB) error

	// GetTaskApplyByID 查询申请详情
	GetTaskApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatTaskApply, error)

	// GetTaskApplyByIDForUpdate 带锁查询申请（事务内使用）
	GetTaskApplyByIDForUpdate(ctx context.Context, applyID uint64, tx *gorm.DB) (*CatTaskApply, error)

	// ListTaskApplies 申请列表（志愿者查我的 / 管理员查待审核）
	ListTaskApplies(ctx context.Context, filter *TaskApplyFilter, tx ...*gorm.DB) ([]*CatTaskApply, int64, error)

	// CancelTaskApply 志愿者取消申请（Pending -> Cancelled）
	CancelTaskApply(ctx context.Context, applyID uint64, userID uint64, reason string, tx ...*gorm.DB) error

	// ApproveTaskApply 管理员审核通过（Pending -> Approved），同时关联taskID
	ApproveTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, urgencyLevel string, difficultyLevel, rewardPoints int32, taskID uint64, tx ...*gorm.DB) error

	// RejectTaskApply 管理员驳回申请（Pending -> Rejected）
	RejectTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error

	// ===================== 正式任务（CatTask）====================

	// CreateTask 创建正式任务
	CreateTask(ctx context.Context, task *CatTask, tx ...*gorm.DB) error

	// GetTaskByID 查询任务详情
	GetTaskByID(ctx context.Context, taskID uint64, tx ...*gorm.DB) (*CatTask, error)

	// GetTaskByIDForUpdate 带锁查询（事务内使用）
	GetTaskByIDForUpdate(ctx context.Context, taskID uint64, tx *gorm.DB) (*CatTask, error)

	// UpdateTask 更新任务基础信息
	UpdateTask(ctx context.Context, taskID uint64, values map[string]any, tx ...*gorm.DB) error

	// UpdateTaskStatus 更新任务状态
	UpdateTaskStatus(ctx context.Context, taskID uint64, status, remark string, tx ...*gorm.DB) error

	// ListTasks 任务列表（支持筛选）
	ListTasks(ctx context.Context, filter *TaskFilter, tx ...*gorm.DB) ([]*CatTask, int64, error)

	// ===================== 任务认领（TaskClaim）====================

	// ClaimTask 认领任务
	ClaimTask(ctx context.Context, taskID, userID uint64, tx ...*gorm.DB) error

	// AbandonTask 放弃任务
	AbandonTask(ctx context.Context, taskID, userID uint64, reason string, tx ...*gorm.DB) error

	// CompleteTask 完成任务
	CompleteTask(ctx context.Context, taskID, userID uint64, content string, imageURLs string, tx ...*gorm.DB) error

	// GetClaimByTaskAndUser 查询某用户对某任务的认领记录
	GetClaimByTaskAndUser(ctx context.Context, taskID, userID uint64, tx ...*gorm.DB) (*TaskClaim, error)

	// ListClaims 认领记录列表
	ListClaims(ctx context.Context, filter *TaskClaimFilter, tx ...*gorm.DB) ([]*TaskClaim, int64, error)

	// ListMyClaimedTaskIDs 查询用户认领的任务ID列表
	ListMyClaimedTaskIDs(ctx context.Context, userID uint64, status string, tx ...*gorm.DB) ([]uint64, error)

	// ListMyClaimedTasks 查询用户认领的任务（含认领详情）
	ListMyClaimedTasks(ctx context.Context, userID uint64, claimStatus string, page, pageSize int, tx ...*gorm.DB) ([]*CatTask, []*TaskClaim, int64, error)
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

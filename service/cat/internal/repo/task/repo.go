package task

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ==================== Repository 接口定义 ====================

type Repository interface {
	// ==================== 任务申请相关 ====================

	// CreateTaskApply 创建任务申请
	CreateTaskApply(ctx context.Context, apply *CatTaskApply, tx ...*gorm.DB) error

	// GetTaskApplyByID 获取任务申请详情
	GetTaskApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatTaskApply, error)

	// GetTaskApplyByCatAndUser 获取用户对某只猫的申请
	GetTaskApplyByCatAndUser(ctx context.Context, catID, userID uint64, tx ...*gorm.DB) (*CatTaskApply, error)

	// ListTaskApplies 任务申请列表
	ListTaskApplies(ctx context.Context, filter TaskApplyFilter, tx ...*gorm.DB) ([]*CatTaskApply, int64, error)

	// UpdateTaskApply 更新任务申请
	UpdateTaskApply(ctx context.Context, applyID uint64, values map[string]any, tx ...*gorm.DB) error

	// CancelTaskApply 取消任务申请
	CancelTaskApply(ctx context.Context, applyID uint64, userID uint64, reason string, tx ...*gorm.DB) error

	// ApproveTaskApply 审核通过任务申请，生成正式任务
	ApproveTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, tx ...*gorm.DB) (*CatTask, error)

	// RejectTaskApply 拒绝任务申请
	RejectTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error

	// DeleteTaskApply 软删除任务申请
	DeleteTaskApply(ctx context.Context, applyID uint64, tx ...*gorm.DB) error

	// ==================== 正式任务相关 ====================

	// CreateTask 创建正式任务
	CreateTask(ctx context.Context, task *CatTask, tx ...*gorm.DB) error

	// GetTaskByID 获取任务详情
	GetTaskByID(ctx context.Context, taskID uint64, tx ...*gorm.DB) (*CatTask, error)

	// ListTasks 任务列表
	ListTasks(ctx context.Context, filter TaskFilter, tx ...*gorm.DB) ([]*CatTask, int64, error)

	// UpdateTask 更新任务
	UpdateTask(ctx context.Context, taskID uint64, values map[string]any, tx ...*gorm.DB) error

	// CancelTask 取消任务
	CancelTask(ctx context.Context, taskID uint64, userID uint64, reason string, tx ...*gorm.DB) error

	// EscalateTaskUrgency 升级任务紧急度
	EscalateTaskUrgency(ctx context.Context, taskID uint64, tx ...*gorm.DB) error

	// DeleteTask 软删除任务
	DeleteTask(ctx context.Context, taskID uint64, tx ...*gorm.DB) error

	// ==================== 任务认领相关 ====================

	// CreateTaskClaim 创建任务认领
	CreateTaskClaim(ctx context.Context, claim *TaskClaim, tx ...*gorm.DB) error

	// GetTaskClaimByID 获取认领详情
	GetTaskClaimByID(ctx context.Context, claimID uint64, tx ...*gorm.DB) (*TaskClaim, error)

	// GetTaskClaimByTaskAndUser 获取用户对任务的认领记录
	GetTaskClaimByTaskAndUser(ctx context.Context, taskID, userID uint64, tx ...*gorm.DB) (*TaskClaim, error)

	// ListTaskClaims 认领列表
	ListTaskClaims(ctx context.Context, filter TaskClaimFilter, tx ...*gorm.DB) ([]*TaskClaim, int64, error)

	// UpdateTaskClaim 更新认领
	UpdateTaskClaim(ctx context.Context, claimID uint64, values map[string]any, tx ...*gorm.DB) error

	// CompleteTaskClaim 完成任务认领
	CompleteTaskClaim(ctx context.Context, claimID uint64, userID uint64, content, imageURLs string, tx ...*gorm.DB) error

	// AbandonTaskClaim 放弃任务认领
	AbandonTaskClaim(ctx context.Context, claimID uint64, userID uint64, reason string, tx ...*gorm.DB) error

	// DeleteTaskClaim 软删除认领
	DeleteTaskClaim(ctx context.Context, claimID uint64, tx ...*gorm.DB) error

	// ==================== 任务日志相关 ====================

	// CreateTaskLog 创建任务日志
	CreateTaskLog(ctx context.Context, log *TaskLog, tx ...*gorm.DB) error

	// ListTaskLogs 任务日志列表
	ListTaskLogs(ctx context.Context, taskID uint64, tx ...*gorm.DB) ([]*TaskLog, error)

	// ==================== 任务流转相关 ====================

	// CreateTaskFlow 创建任务流转记录
	CreateTaskFlow(ctx context.Context, flow *TaskFlow, tx ...*gorm.DB) error

	// ListTaskFlows 任务流转记录列表
	ListTaskFlows(ctx context.Context, taskID uint64, tx ...*gorm.DB) ([]*TaskFlow, error)
}

// ==================== 实现结构体 ====================

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

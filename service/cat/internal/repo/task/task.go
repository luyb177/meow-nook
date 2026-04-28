package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 状态常量 ====================

const (
	// 任务申请状态
	TaskApplyStatusPending   = "pending"   // 待审核
	TaskApplyStatusApproved  = "approved"  // 审核通过（已转为正式任务）
	TaskApplyStatusRejected  = "rejected"  // 审核拒绝
	TaskApplyStatusCancelled = "cancelled" // 申请人主动取消

	// 正式任务状态
	TaskStatusPending    = "pending"    // 待认领
	TaskStatusProcessing = "processing" // 进行中
	TaskStatusCompleted  = "completed"  // 已完成
	TaskStatusCancelled  = "cancelled"  // 已取消

	// 紧急程度
	UrgencyHigh   = "high"   // 高危（建议2小时内领取）
	UrgencyUrgent = "urgent" // 紧急（建议24小时内领取）
	UrgencyNormal = "normal" // 一般
)

// ==================== 错误定义 ====================

var (
	ErrTaskNotFound         = errors.New("任务不存在")
	ErrTaskApplyNotFound    = errors.New("任务申请不存在")
	ErrTaskApplyNotPending  = errors.New("申请状态不是待审核")
	ErrTaskApplyNotBelongTo = errors.New("该申请不属于当前用户")
	ErrTaskNotPending       = errors.New("任务状态不是待认领")
	ErrTaskFull             = errors.New("任务认领人数已满")
	ErrAlreadyClaimed       = errors.New("您已认领该任务")
	ErrClaimNotFound        = errors.New("认领记录不存在")
	ErrClaimNotBelongTo     = errors.New("该认领记录不属于当前用户")
	ErrTaskAlreadyClaimed   = errors.New("该任务已被认领")
	ErrCannotComplete       = errors.New("只有认领人才能完成任务")
	ErrTaskNotProcessing    = errors.New("任务不在进行中状态")
	ErrRejectReasonRequired = errors.New("拒绝原因不能为空")
)

// ==================== 数据模型 ====================

// CatTask 正式任务表
type CatTask struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	CatID     uint64 `gorm:"not null;index;comment:关联猫咪ID" json:"cat_id"`
	ApplyID   uint64 `gorm:"default:0;index;comment:来源申请ID，0表示管理员直接创建" json:"apply_id"`
	CreatorID uint64 `gorm:"not null;comment:创建人ID" json:"creator_id"`

	// 任务基础信息
	Title    string `gorm:"type:varchar(200);not null;comment:任务标题" json:"title"`
	TaskType string `gorm:"type:varchar(50);not null;comment:任务类型:feed/sterilize/vaccine/rescue/other" json:"task_type"`
	Summary  string `gorm:"type:varchar(500);comment:任务简介" json:"summary"`
	Detail   string `gorm:"type:text;comment:任务详细描述（支持富文本）" json:"detail"`

	// 任务优先级
	UrgencyLevel    string `gorm:"type:varchar(20);default:'normal';index;comment:紧急程度:high/urgent/normal" json:"urgency_level"`
	DifficultyLevel int32  `gorm:"default:1;comment:难度等级1-5" json:"difficulty_level"`
	RewardPoints    int32  `gorm:"default:0;comment:完成任务奖励积分" json:"reward_points"`

	// 任务限制
	MaxClaimers     int32 `gorm:"default:1;comment:最多认领人数" json:"max_claimers"`
	CurrentClaimers int32 `gorm:"default:0;comment:当前认领人数" json:"current_claimers"`

	// 任务状态
	Status string     `gorm:"type:varchar(30);default:'pending';index;comment:任务状态:pending/processing/completed/cancelled" json:"status"`
	Remark string     `gorm:"type:text;comment:管理员备注" json:"remark"`
	DeadlineAt *time.Time `gorm:"index;comment:任务截止时间" json:"deadline_at"`

	// 地区信息（冗余，方便筛选）
	Area string `gorm:"type:varchar(100);comment:所在区域" json:"area"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatTask) TableName() string {
	return "cat_tasks"
}

// ==================== 查询过滤器 ====================

type TaskFilter struct {
	CatID        uint64
	Status       string
	UrgencyLevel string
	DifficultyLevel int32
	Area         string

	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time
	DeadlineStart  *time.Time
	DeadlineEnd    *time.Time

	Page     int
	PageSize int
}

// ==================== 任务相关方法 ====================

func (r *repository) CreateTask(ctx context.Context, task *CatTask, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	return tx.Create(task).Error
}

func (r *repository) GetTaskByID(ctx context.Context, taskID uint64, txs ...*gorm.DB) (*CatTask, error) {
	tx := r.getDB(ctx, txs...)
	var task CatTask
	err := tx.First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *repository) GetTaskByIDForUpdate(ctx context.Context, taskID uint64, tx *gorm.DB) (*CatTask, error) {
	var task CatTask
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *repository) UpdateTask(ctx context.Context, taskID uint64, values map[string]any, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	return tx.Model(&CatTask{}).Where("id = ?", taskID).Updates(values).Error
}

func (r *repository) ListTasks(ctx context.Context, filter *TaskFilter, txs ...*gorm.DB) ([]*CatTask, int64, error) {
	tx := r.getDB(ctx, txs...)
	var total int64

	query := tx.Model(&CatTask{})

	if filter.CatID > 0 {
		query = query.Where("cat_id = ?", filter.CatID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.UrgencyLevel != "" {
		query = query.Where("urgency_level = ?", filter.UrgencyLevel)
	}
	if filter.DifficultyLevel > 0 {
		query = query.Where("difficulty_level = ?", filter.DifficultyLevel)
	}
	if filter.Area != "" {
		query = query.Where("area LIKE ?", "%"+filter.Area+"%")
	}
	if filter.CreatedAtStart != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		query = query.Where("created_at <= ?", *filter.CreatedAtEnd)
	}
	if filter.DeadlineStart != nil {
		query = query.Where("deadline_at >= ?", *filter.DeadlineStart)
	}
	if filter.DeadlineEnd != nil {
		query = query.Where("deadline_at <= ?", *filter.DeadlineEnd)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (filter.Page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	var tasks []*CatTask
	err := query.Order("urgency_level DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, total, err
}

func (r *repository) UpdateTaskStatus(ctx context.Context, taskID uint64, status, remark string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	return tx.Model(&CatTask{}).Where("id = ?", taskID).Updates(map[string]any{
		"status":     status,
		"remark":     remark,
		"updated_at": time.Now(),
	}).Error
}

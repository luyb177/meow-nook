package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 常量 ====================

const (

	// 正式任务状态
	TaskStatusPending    = "pending"    // 待认领
	TaskStatusProcessing = "processing" // 进行中
	TaskStatusCompleted  = "completed"  // 已完成
	TaskStatusCancelled  = "cancelled"  // 已取消

	// 紧急程度
	UrgencyLevelHigh   = "high"
	UrgencyLevelUrgent = "urgent"
	UrgencyLevelNormal = "normal"
)

// ==================== CatTask 正式任务表 ====================

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

	// 位置信息（适配地图展示）
	Location  string  `gorm:"type:varchar(255);comment:任务地点描述" json:"location"`
	Longitude float64 `gorm:"type:decimal(10,6);default:0;comment:经度" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal(10,6);default:0;comment:纬度" json:"latitude"`
	Area      string  `gorm:"type:varchar(100);comment:所在区域（冗余，方便筛选）" json:"area"`

	// 任务优先级
	UrgencyLevel      string `gorm:"type:varchar(20);default:'normal';index;comment:紧急程度:high/urgent/normal" json:"urgency_level"`
	DifficultyLevel   int32  `gorm:"default:1;comment:难度等级1-5" json:"difficulty_level"`
	RewardPoints      int32  `gorm:"default:0;comment:基础奖励积分" json:"reward_points"`
	FinalRewardPoints int32  `gorm:"default:0;comment:最终获得积分（可能因超时/质量调整）" json:"final_reward_points"`

	// 任务限制
	MaxClaimers     int32 `gorm:"default:1;comment:最多认领人数" json:"max_claimers"`
	CurrentClaimers int32 `gorm:"default:0;comment:当前认领人数" json:"current_claimers"`

	// 任务状态
	Status     string     `gorm:"type:varchar(30);default:'pending';index;comment:任务状态:pending/processing/completed/cancelled" json:"status"`
	Remark     string     `gorm:"type:text;comment:管理员备注" json:"remark"`
	DeadlineAt *time.Time `gorm:"index;comment:任务截止时间" json:"deadline_at"`

	// 自动升级紧急度相关
	LastEscalatedAt *time.Time `gorm:"comment:上次自动升级紧急度时间" json:"last_escalated_at"`
	EscalationCount int32      `gorm:"default:0;comment:已升级次数" json:"escalation_count"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatTask) TableName() string {
	return "cat_tasks"
}

type TaskFilter struct {
	CatID           uint64
	Status          string
	UrgencyLevel    string
	TaskType        string
	Area            string
	CreatorID       uint64
	IsOverdue       *bool
	DifficultyLevel int32

	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time
	DeadlineStart  *time.Time
	DeadlineEnd    *time.Time

	Page     int
	PageSize int
}

// ==================== 正式任务相关实现 ====================

func (r *repository) CreateTask(ctx context.Context, task *CatTask, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(task).Error
}

func (r *repository) GetTaskByID(ctx context.Context, taskID uint64, tx ...*gorm.DB) (*CatTask, error) {
	var task CatTask
	err := r.getDB(ctx, tx...).First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *repository) ListTasks(ctx context.Context, filter TaskFilter, tx ...*gorm.DB) ([]*CatTask, int64, error) {
	db := r.getDB(ctx, tx...).Model(&CatTask{})

	if filter.CatID > 0 {
		db = db.Where("cat_id = ?", filter.CatID)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.UrgencyLevel != "" {
		db = db.Where("urgency_level = ?", filter.UrgencyLevel)
	}
	if filter.TaskType != "" {
		db = db.Where("task_type = ?", filter.TaskType)
	}
	if filter.Area != "" {
		db = db.Where("area = ?", filter.Area)
	}
	if filter.CreatorID > 0 {
		db = db.Where("creator_id = ?", filter.CreatorID)
	}
	if filter.IsOverdue != nil {
		if *filter.IsOverdue {
			db = db.Where("deadline_at < ? AND status != ?", time.Now(), TaskStatusCompleted)
		} else {
			db = db.Where("deadline_at >= ? OR status = ?", time.Now(), TaskStatusCompleted)
		}
	}
	if filter.CreatedAtStart != nil {
		db = db.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		db = db.Where("created_at <= ?", *filter.CreatedAtEnd)
	}
	if filter.DeadlineStart != nil {
		db = db.Where("deadline_at >= ?", *filter.DeadlineStart)
	}
	if filter.DeadlineEnd != nil {
		db = db.Where("deadline_at <= ?", *filter.DeadlineEnd)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)
	var list []*CatTask
	err := db.Order("urgency_level DESC, deadline_at ASC, created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *repository) UpdateTask(ctx context.Context, taskID uint64, values map[string]any, tx ...*gorm.DB) error {
	res := r.getDB(ctx, tx...).Model(&CatTask{}).Where("id = ?", taskID).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}

func (r *repository) CancelTask(ctx context.Context, taskID uint64, userID uint64, reason string, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var task CatTask
		if err := tx.First(&task, taskID).Error; err != nil {
			return err
		}
		if task.CreatorID != userID {
			return errors.New("not authorized to cancel this task")
		}
		if task.Status != TaskStatusPending && task.Status != TaskStatusProcessing {
			return errors.New("only pending or processing task can be cancelled")
		}

		task.Status = TaskStatusCancelled
		task.Remark = reason
		return tx.Save(&task).Error
	})
}

func (r *repository) EscalateTaskUrgency(ctx context.Context, taskID uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var task CatTask
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			return err
		}
		if task.Status != TaskStatusPending {
			return errors.New("only pending task can be escalated")
		}

		newUrgency := task.UrgencyLevel
		switch task.UrgencyLevel {
		case UrgencyLevelNormal:
			newUrgency = UrgencyLevelUrgent
		case UrgencyLevelUrgent:
			newUrgency = UrgencyLevelHigh
		case UrgencyLevelHigh:
			return nil // 已经是最高级别
		}

		task.UrgencyLevel = newUrgency
		task.EscalationCount++
		task.LastEscalatedAt = ptr(time.Now())
		return tx.Save(&task).Error
	})
}

func (r *repository) DeleteTask(ctx context.Context, taskID uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&CatTask{}, taskID).Error
}

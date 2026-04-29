package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

const (
	// 认领状态
	ClaimStatusActive    = "active"
	ClaimStatusCompleted = "completed"
	ClaimStatusAbandoned = "abandoned"
)

// TaskClaim 任务认领记录表
type TaskClaim struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	TaskID uint64 `gorm:"not null;index;comment:任务ID" json:"task_id"`
	UserID uint64 `gorm:"not null;index;comment:认领人ID" json:"user_id"`

	// 状态
	Status string `gorm:"type:varchar(30);default:'active';index;comment:认领状态:active/completed/abandoned" json:"status"`

	// 完成/放弃信息
	CompletedAt   *time.Time `gorm:"comment:完成时间" json:"completed_at"`
	AbandonedAt   *time.Time `gorm:"comment:放弃时间" json:"abandoned_at"`
	AbandonReason string     `gorm:"type:text;comment:放弃原因" json:"abandon_reason"`

	// 进度内容（完成时的照片/备注）
	Content   string `gorm:"type:text;comment:完成备注" json:"content"`
	ImageURLs string `gorm:"type:text;comment:完成照片（JSON数组）" json:"image_urls"`

	// 信誉分相关
	IsOverdue           bool  `gorm:"default:false;comment:是否超时完成" json:"is_overdue"`
	ReputationDeduction int32 `gorm:"default:0;comment:信誉分扣减" json:"reputation_deduction"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:认领时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (TaskClaim) TableName() string {
	return "task_claims"
}

type TaskClaimFilter struct {
	TaskID    uint64
	UserID    uint64
	Status    string
	IsOverdue *bool

	Page     int
	PageSize int
}

// ==================== 任务认领相关实现 ====================

func (r *repository) CreateTaskClaim(ctx context.Context, claim *TaskClaim, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(claim).Error; err != nil {
			return err
		}

		// 更新任务当前认领人数
		return tx.Model(&CatTask{}).Where("id = ?", claim.TaskID).Update("current_claimers", gorm.Expr("current_claimers + 1")).Error
	})
}

func (r *repository) GetTaskClaimByID(ctx context.Context, claimID uint64, tx ...*gorm.DB) (*TaskClaim, error) {
	var claim TaskClaim
	err := r.getDB(ctx, tx...).First(&claim, claimID).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func (r *repository) GetTaskClaimByTaskAndUser(ctx context.Context, taskID, userID uint64, tx ...*gorm.DB) (*TaskClaim, error) {
	var claim TaskClaim
	err := r.getDB(ctx, tx...).
		Where("task_id = ? AND user_id = ? AND status = ?", taskID, userID, ClaimStatusActive).
		First(&claim).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &claim, nil
}

func (r *repository) ListTaskClaims(ctx context.Context, filter TaskClaimFilter, tx ...*gorm.DB) ([]*TaskClaim, int64, error) {
	db := r.getDB(ctx, tx...).Model(&TaskClaim{})

	if filter.TaskID > 0 {
		db = db.Where("task_id = ?", filter.TaskID)
	}
	if filter.UserID > 0 {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.IsOverdue != nil {
		db = db.Where("is_overdue = ?", *filter.IsOverdue)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)
	var list []*TaskClaim
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *repository) UpdateTaskClaim(ctx context.Context, claimID uint64, values map[string]any, tx ...*gorm.DB) error {
	res := r.getDB(ctx, tx...).Model(&TaskClaim{}).Where("id = ?", claimID).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("task claim not found")
	}
	return nil
}

func (r *repository) CompleteTaskClaim(ctx context.Context, claimID uint64, userID uint64, content, imageURLs string, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var claim TaskClaim
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&claim, claimID).Error; err != nil {
			return err
		}
		if claim.UserID != userID {
			return errors.New("not authorized to complete this claim")
		}
		if claim.Status != ClaimStatusActive {
			return errors.New("only active claim can be completed")
		}

		now := time.Now()
		isOverdue := false
		var task CatTask
		if err := tx.First(&task, claim.TaskID).Error; err == nil {
			isOverdue = task.DeadlineAt != nil && now.After(*task.DeadlineAt)
		}

		claim.Status = ClaimStatusCompleted
		claim.CompletedAt = ptr(now)
		claim.Content = content
		claim.ImageURLs = imageURLs
		claim.IsOverdue = isOverdue
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}

		// 更新任务状态为已完成（如果所有认领都已完成）
		var completedCount int64
		if err := tx.Model(&TaskClaim{}).Where("task_id = ? AND status = ?", claim.TaskID, ClaimStatusCompleted).Count(&completedCount).Error; err != nil {
			return err
		}
		var totalCount int64
		if err := tx.Model(&TaskClaim{}).Where("task_id = ?", claim.TaskID).Count(&totalCount).Error; err != nil {
			return err
		}
		if completedCount == totalCount {
			return tx.Model(&CatTask{}).Where("id = ?", claim.TaskID).Update("status", TaskStatusCompleted).Error
		}
		return nil
	})
}

func (r *repository) AbandonTaskClaim(ctx context.Context, claimID uint64, userID uint64, reason string, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var claim TaskClaim
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&claim, claimID).Error; err != nil {
			return err
		}
		if claim.UserID != userID {
			return errors.New("not authorized to abandon this claim")
		}
		if claim.Status != ClaimStatusActive {
			return errors.New("only active claim can be abandoned")
		}

		now := time.Now()
		claim.Status = ClaimStatusAbandoned
		claim.AbandonedAt = ptr(now)
		claim.AbandonReason = reason
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}

		// 更新任务当前认领人数
		return tx.Model(&CatTask{}).Where("id = ?", claim.TaskID).Update("current_claimers", gorm.Expr("current_claimers - 1")).Error
	})
}

func (r *repository) DeleteTaskClaim(ctx context.Context, claimID uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&TaskClaim{}, claimID).Error
}

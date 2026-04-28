package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 任务认领模型 ====================

const (
	ClaimStatusActive    = "active"    // 认领中
	ClaimStatusCompleted = "completed" // 已完成
	ClaimStatusAbandoned = "abandoned" // 已放弃
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
	CompletedAt  *time.Time `gorm:"comment:完成时间" json:"completed_at"`
	AbandonedAt  *time.Time `gorm:"comment:放弃时间" json:"abandoned_at"`
	AbandonReason string    `gorm:"type:text;comment:放弃原因" json:"abandon_reason"`

	// 进度内容（完成时的照片/备注）
	Content   string `gorm:"type:text;comment:完成备注" json:"content"`
	ImageURLs string `gorm:"type:text;comment:完成照片（JSON数组）" json:"image_urls"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:认领时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (TaskClaim) TableName() string {
	return "task_claims"
}

// ==================== 查询过滤器 ====================

type TaskClaimFilter struct {
	TaskID uint64
	UserID uint64
	Status string

	Page     int
	PageSize int
}

// ==================== 方法实现 ====================

func (r *repository) ClaimTask(ctx context.Context, taskID, userID uint64, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 锁定任务记录
		var task CatTask
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskNotFound
			}
			return err
		}

		if task.Status != TaskStatusPending {
			return ErrTaskNotPending
		}

		if task.CurrentClaimers >= task.MaxClaimers {
			return ErrTaskFull
		}

		// 检查用户是否已认领该任务
		var existClaim TaskClaim
		err := tx.Where("task_id = ? AND user_id = ? AND status = ?", taskID, userID, ClaimStatusActive).First(&existClaim).Error
		if err == nil {
			return ErrAlreadyClaimed
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 创建认领记录
		claim := &TaskClaim{
			TaskID: taskID,
			UserID: userID,
			Status: ClaimStatusActive,
		}
		if err := tx.Create(claim).Error; err != nil {
			return err
		}

		// 更新任务认领人数，若达到上限则改为processing
		newCount := task.CurrentClaimers + 1
		updateData := map[string]any{
			"current_claimers": newCount,
			"updated_at":       time.Now(),
		}
		if newCount >= task.MaxClaimers {
			updateData["status"] = TaskStatusProcessing
		}

		return tx.Model(&CatTask{}).Where("id = ?", taskID).Updates(updateData).Error
	})
}

func (r *repository) AbandonTask(ctx context.Context, taskID, userID uint64, reason string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 查找认领记录
		var claim TaskClaim
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("task_id = ? AND user_id = ? AND status = ?", taskID, userID, ClaimStatusActive).
			First(&claim).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrClaimNotFound
			}
			return err
		}

		// 更新认领记录
		now := time.Now()
		claim.Status = ClaimStatusAbandoned
		claim.AbandonedAt = &now
		claim.AbandonReason = reason
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}

		// 更新任务认领人数
		var task CatTask
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&task, taskID).Error; err != nil {
			return err
		}

		newCount := task.CurrentClaimers - 1
		if newCount < 0 {
			newCount = 0
		}
		updateData := map[string]any{
			"current_claimers": newCount,
			"updated_at":       time.Now(),
		}
		// 如果人数降到0，任务重回待认领
		if newCount == 0 && task.Status == TaskStatusProcessing {
			updateData["status"] = TaskStatusPending
		}

		return tx.Model(&CatTask{}).Where("id = ?", taskID).Updates(updateData).Error
	})
}

func (r *repository) CompleteTask(ctx context.Context, taskID, userID uint64, content string, imageURLs string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 查找认领记录
		var claim TaskClaim
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("task_id = ? AND user_id = ? AND status = ?", taskID, userID, ClaimStatusActive).
			First(&claim).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrClaimNotFound
			}
			return err
		}

		// 更新认领记录
		now := time.Now()
		claim.Status = ClaimStatusCompleted
		claim.CompletedAt = &now
		claim.Content = content
		claim.ImageURLs = imageURLs
		if err := tx.Save(&claim).Error; err != nil {
			return err
		}

		// 将任务标记为已完成
		return tx.Model(&CatTask{}).Where("id = ?", taskID).Updates(map[string]any{
			"status":     TaskStatusCompleted,
			"updated_at": time.Now(),
		}).Error
	})
}

func (r *repository) GetClaimByTaskAndUser(ctx context.Context, taskID, userID uint64, txs ...*gorm.DB) (*TaskClaim, error) {
	tx := r.getDB(ctx, txs...)
	var claim TaskClaim
	err := tx.Where("task_id = ? AND user_id = ?", taskID, userID).Order("id DESC").First(&claim).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClaimNotFound
		}
		return nil, err
	}
	return &claim, nil
}

func (r *repository) ListClaims(ctx context.Context, filter *TaskClaimFilter, txs ...*gorm.DB) ([]*TaskClaim, int64, error) {
	tx := r.getDB(ctx, txs...)
	var total int64

	query := tx.Model(&TaskClaim{})

	if filter.TaskID > 0 {
		query = query.Where("task_id = ?", filter.TaskID)
	}
	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
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

	var claims []*TaskClaim
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&claims).Error
	return claims, total, err
}

// ListMyClaimedTaskIDs returns task IDs claimed by a user (for joining with task table)
func (r *repository) ListMyClaimedTaskIDs(ctx context.Context, userID uint64, status string, txs ...*gorm.DB) ([]uint64, error) {
	tx := r.getDB(ctx, txs...)
	query := tx.Model(&TaskClaim{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var taskIDs []uint64
	err := query.Pluck("task_id", &taskIDs).Error
	return taskIDs, err
}

// ListMyClaimedTasksWithDetails returns tasks + claim info for a user
func (r *repository) ListMyClaimedTasks(ctx context.Context, userID uint64, claimStatus string, page, pageSize int, txs ...*gorm.DB) ([]*CatTask, []*TaskClaim, int64, error) {
	tx := r.getDB(ctx, txs...)

	// First get claims
	claimQuery := tx.Model(&TaskClaim{}).Where("user_id = ?", userID)
	if claimStatus != "" {
		claimQuery = claimQuery.Where("status = ?", claimStatus)
	}

	var total int64
	if err := claimQuery.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	var claims []*TaskClaim
	err := claimQuery.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&claims).Error
	if err != nil {
		return nil, nil, 0, err
	}

	if len(claims) == 0 {
		return nil, nil, 0, nil
	}

	// Get corresponding task IDs
	taskIDs := make([]uint64, 0, len(claims))
	for _, c := range claims {
		taskIDs = append(taskIDs, c.TaskID)
	}

	// Fetch tasks
	var tasks []*CatTask
	err = tx.Where("id IN ?", taskIDs).Find(&tasks).Error
	return tasks, claims, total, err
}

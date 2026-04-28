package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
)

// ==================== 任务申请模型 ====================

// CatTaskApply 志愿者申请创建任务表
type CatTaskApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	CatID           uint64 `gorm:"not null;index;comment:关联猫咪ID" json:"cat_id"`
	ApplicantUserID uint64 `gorm:"not null;index;comment:申请人ID" json:"applicant_user_id"`

	// 申请信息
	Title    string `gorm:"type:varchar(200);not null;comment:任务标题" json:"title"`
	TaskType string `gorm:"type:varchar(50);not null;comment:任务类型" json:"task_type"`
	Summary  string `gorm:"type:varchar(500);comment:任务简介" json:"summary"`
	Detail   string `gorm:"type:text;comment:任务详细描述" json:"detail"`

	DeadlineAt *time.Time `gorm:"comment:任务截止时间" json:"deadline_at"`

	// 审核信息
	Status       string     `gorm:"type:varchar(30);default:'pending';index;comment:申请状态:pending/approved/rejected/cancelled" json:"status"`
	ReviewerID   uint64     `gorm:"default:0;comment:审核人ID" json:"reviewer_id"`
	ReviewedAt   *time.Time `gorm:"comment:审核时间" json:"reviewed_at"`
	RejectReason string     `gorm:"type:text;comment:拒绝原因/取消原因" json:"reject_reason"`

	// 审核通过时填写的正式任务信息
	UrgencyLevel    string `gorm:"type:varchar(20);comment:紧急程度" json:"urgency_level"`
	DifficultyLevel int32  `gorm:"default:0;comment:难度等级" json:"difficulty_level"`
	RewardPoints    int32  `gorm:"default:0;comment:奖励积分" json:"reward_points"`

	// 审核通过后关联的正式任务ID
	TaskID uint64 `gorm:"default:0;index;comment:关联的正式任务ID" json:"task_id"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (CatTaskApply) TableName() string {
	return "cat_task_applies"
}

// ==================== 查询过滤器 ====================

type TaskApplyFilter struct {
	ApplicantUserID uint64
	Status          string
	CatID           uint64

	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time

	Page     int
	PageSize int
}

// ==================== 方法实现 ====================

func (r *repository) CreateTaskApply(ctx context.Context, apply *CatTaskApply, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	return tx.Create(apply).Error
}

func (r *repository) GetTaskApplyByID(ctx context.Context, applyID uint64, txs ...*gorm.DB) (*CatTaskApply, error) {
	tx := r.getDB(ctx, txs...)
	var apply CatTaskApply
	err := tx.First(&apply, applyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskApplyNotFound
		}
		return nil, err
	}
	return &apply, nil
}

func (r *repository) GetTaskApplyByIDForUpdate(ctx context.Context, applyID uint64, tx *gorm.DB) (*CatTaskApply, error) {
	var apply CatTaskApply
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskApplyNotFound
		}
		return nil, err
	}
	return &apply, nil
}

func (r *repository) ListTaskApplies(ctx context.Context, filter *TaskApplyFilter, txs ...*gorm.DB) ([]*CatTaskApply, int64, error) {
	tx := r.getDB(ctx, txs...)
	var total int64

	query := tx.Model(&CatTaskApply{})

	if filter.ApplicantUserID > 0 {
		query = query.Where("applicant_user_id = ?", filter.ApplicantUserID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CatID > 0 {
		query = query.Where("cat_id = ?", filter.CatID)
	}
	if filter.CreatedAtStart != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		query = query.Where("created_at <= ?", *filter.CreatedAtEnd)
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

	var applies []*CatTaskApply
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&applies).Error
	return applies, total, err
}

func (r *repository) CancelTaskApply(ctx context.Context, applyID uint64, userID uint64, reason string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskApplyNotFound
			}
			return err
		}

		if apply.ApplicantUserID != userID {
			return ErrTaskApplyNotBelongTo
		}

		if apply.Status != TaskApplyStatusPending {
			return ErrTaskApplyNotPending
		}

		now := time.Now()
		apply.Status = TaskApplyStatusCancelled
		apply.RejectReason = reason
		apply.ReviewedAt = &now

		return tx.Save(&apply).Error
	})
}

func (r *repository) ApproveTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, urgencyLevel string, difficultyLevel, rewardPoints int32, taskID uint64, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskApplyNotFound
			}
			return err
		}

		if apply.Status != TaskApplyStatusPending {
			return ErrTaskApplyNotPending
		}

		now := time.Now()
		apply.Status = TaskApplyStatusApproved
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = &now
		apply.UrgencyLevel = urgencyLevel
		apply.DifficultyLevel = difficultyLevel
		apply.RewardPoints = rewardPoints
		apply.TaskID = taskID

		return tx.Save(&apply).Error
	})
}

func (r *repository) RejectTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTaskApplyNotFound
			}
			return err
		}

		if apply.Status != TaskApplyStatusPending {
			return ErrTaskApplyNotPending
		}

		if reason == "" {
			return ErrRejectReasonRequired
		}

		now := time.Now()
		apply.Status = TaskApplyStatusRejected
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = &now
		apply.RejectReason = reason

		return tx.Save(&apply).Error
	})
}

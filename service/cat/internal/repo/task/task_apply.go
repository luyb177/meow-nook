package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	// 任务申请状态
	TaskApplyStatusPending   = "pending"
	TaskApplyStatusApproved  = "approved"
	TaskApplyStatusRejected  = "rejected"
	TaskApplyStatusCancelled = "cancelled"
)

// CatTaskApply 志愿者申请创建任务表
type CatTaskApply struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	CatID           uint64 `gorm:"not null;index;comment:关联猫咪ID" json:"cat_id"`
	ApplicantUserID uint64 `gorm:"not null;index;comment:申请人ID" json:"applicant_user_id"`

	// 申请信息
	Title    string `gorm:"type:varchar(200);not null;comment:任务标题" json:"title"`
	TaskType string `gorm:"type:varchar(50);not null;comment:任务类型:feed/sterilize/vaccine/rescue/other" json:"task_type"`
	Summary  string `gorm:"type:varchar(500);comment:任务简介" json:"summary"`
	Detail   string `gorm:"type:text;comment:任务详细描述" json:"detail"`

	// 位置信息（适配地图展示）
	Location  string  `gorm:"type:varchar(255);comment:发现地点描述" json:"location"`
	Longitude float64 `gorm:"type:decimal(10,6);default:0;comment:经度" json:"longitude"`
	Latitude  float64 `gorm:"type:decimal(10,6);default:0;comment:纬度" json:"latitude"`

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

type TaskApplyFilter struct {
	ApplicantUserID uint64
	CatID           uint64
	Status          string
	TaskType        string

	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time

	Page     int
	PageSize int
}

// ==================== 任务申请相关实现 ====================

func (r *repository) CreateTaskApply(ctx context.Context, apply *CatTaskApply, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(apply).Error
}

func (r *repository) GetTaskApplyByID(ctx context.Context, applyID uint64, tx ...*gorm.DB) (*CatTaskApply, error) {
	var apply CatTaskApply
	err := r.getDB(ctx, tx...).First(&apply, applyID).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repository) GetTaskApplyByCatAndUser(ctx context.Context, catID, userID uint64, tx ...*gorm.DB) (*CatTaskApply, error) {
	var apply CatTaskApply
	err := r.getDB(ctx, tx...).
		Where("cat_id = ? AND applicant_user_id = ?", catID, userID).
		Where("status IN (?)", []string{TaskApplyStatusPending, TaskApplyStatusApproved}).
		First(&apply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apply, nil
}

func (r *repository) ListTaskApplies(ctx context.Context, filter TaskApplyFilter, tx ...*gorm.DB) ([]*CatTaskApply, int64, error) {
	db := r.getDB(ctx, tx...).Model(&CatTaskApply{})

	if filter.ApplicantUserID > 0 {
		db = db.Where("applicant_user_id = ?", filter.ApplicantUserID)
	}
	if filter.CatID > 0 {
		db = db.Where("cat_id = ?", filter.CatID)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.TaskType != "" {
		db = db.Where("task_type = ?", filter.TaskType)
	}
	if filter.CreatedAtStart != nil {
		db = db.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil {
		db = db.Where("created_at <= ?", *filter.CreatedAtEnd)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit, offset := normalizePage(filter.Page, filter.PageSize)
	var list []*CatTaskApply
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *repository) UpdateTaskApply(ctx context.Context, applyID uint64, values map[string]any, tx ...*gorm.DB) error {
	res := r.getDB(ctx, tx...).Model(&CatTaskApply{}).Where("id = ?", applyID).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("task apply not found")
	}
	return nil
}

func (r *repository) CancelTaskApply(ctx context.Context, applyID uint64, userID uint64, reason string, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.First(&apply, applyID).Error; err != nil {
			return err
		}
		if apply.ApplicantUserID != userID {
			return errors.New("not authorized to cancel this apply")
		}
		if apply.Status != TaskApplyStatusPending {
			return errors.New("only pending apply can be cancelled")
		}

		apply.Status = TaskApplyStatusCancelled
		apply.RejectReason = reason
		return tx.Save(&apply).Error
	})
}

func (r *repository) ApproveTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, tx ...*gorm.DB) (*CatTask, error) {
	var task *CatTask
	err := r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.First(&apply, applyID).Error; err != nil {
			return err
		}
		if apply.Status != TaskApplyStatusPending {
			return errors.New("only pending apply can be approved")
		}

		task = &CatTask{
			CatID:           apply.CatID,
			ApplyID:         apply.ID,
			CreatorID:       reviewerID,
			Title:           apply.Title,
			TaskType:        apply.TaskType,
			Summary:         apply.Summary,
			Detail:          apply.Detail,
			Location:        apply.Location,
			Longitude:       apply.Longitude,
			Latitude:        apply.Latitude,
			UrgencyLevel:    apply.UrgencyLevel,
			DifficultyLevel: apply.DifficultyLevel,
			RewardPoints:    apply.RewardPoints,
			DeadlineAt:      apply.DeadlineAt,
		}
		if err := tx.Create(task).Error; err != nil {
			return err
		}

		apply.Status = TaskApplyStatusApproved
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = ptr(time.Now())
		apply.TaskID = task.ID
		return tx.Save(&apply).Error
	})
	return task, err
}

func (r *repository) RejectTaskApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Transaction(func(tx *gorm.DB) error {
		var apply CatTaskApply
		if err := tx.First(&apply, applyID).Error; err != nil {
			return err
		}
		if apply.Status != TaskApplyStatusPending {
			return errors.New("only pending apply can be rejected")
		}

		apply.Status = TaskApplyStatusRejected
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = ptr(time.Now())
		apply.RejectReason = reason
		return tx.Save(&apply).Error
	})
}

func (r *repository) DeleteTaskApply(ctx context.Context, applyID uint64, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Delete(&CatTaskApply{}, applyID).Error
}

package adoption

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
	AdoptApplyStatusPending   = "pending"   // 待审核
	AdoptApplyStatusApproved  = "approved"  // 审核通过（待完成领养手续）
	AdoptApplyStatusRejected  = "rejected"  // 审核拒绝
	AdoptApplyStatusCancelled = "cancelled" // 申请人主动取消
	AdoptApplyStatusCompleted = "completed" // 领养手续完成（已转入 adoptions 表）
	AdoptApplyStatusExpired   = "expired"   // 审核通过后超时未完成手续
)

// ==================== 错误定义 ====================

var (
	ErrApplyNotFound            = errors.New("领养申请不存在")
	ErrApplyNotPending          = errors.New("申请状态不是待审核，无法操作")
	ErrApplyNotBelongToUser     = errors.New("该申请不属于当前用户")
	ErrCatNotAvailable          = errors.New("该猫咪当前不可领养")
	ErrDuplicateApply           = errors.New("您已申请过这只猫咪，请勿重复申请")
	ErrMaxApplyLimitReached     = errors.New("该猫咪申请人数已满，请稍后再试")
	ErrInsufficientCreditScore  = errors.New("积分不足，无法申请领养")
	ErrRejectReasonRequired     = errors.New("拒绝原因不能为空")
	ErrApplyAlreadyReviewed     = errors.New("该申请已被审核")
	ErrApplyCannotCancel        = errors.New("当前状态无法取消申请")
	ErrApplyExpiredCannotReview = errors.New("该申请已过期")
)

// ==================== 数据模型 ====================

type AdoptApplication struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// ==================== 关联信息 ====================
	CatID       uint64 `gorm:"not null;index;comment:猫咪ID" json:"cat_id"`
	ApplicantID uint64 `gorm:"not null;index;comment:申请人ID（志愿者）" json:"applicant_id"`

	// ==================== 申请信息 ====================
	ApplyReason   string `gorm:"type:text;comment:申请理由" json:"apply_reason"`
	ContactPhone  string `gorm:"type:varchar(20);comment:联系电话" json:"contact_phone"`
	ContactWechat string `gorm:"type:varchar(64);comment:微信号" json:"contact_wechat"`

	// 申请时积分快照（防止后续积分变化影响历史记录）
	ApplicantCreditScore int `gorm:"default:0;comment:申请人积分快照" json:"applicant_credit_score"`

	// ==================== 审核信息 ====================
	Status       string     `gorm:"type:varchar(30);default:'pending';index;comment:申请状态" json:"status"`
	ReviewerID   uint64     `gorm:"default:0;index;comment:审核人ID" json:"reviewer_id"`
	ReviewedAt   *time.Time `gorm:"comment:审核时间" json:"reviewed_at"`
	RejectReason string     `gorm:"type:text;comment:拒绝原因" json:"reject_reason"`

	// ==================== 审核通过后的时间控制 ====================
	// 审核通过后，申请人需要在 7 天内完成：签订协议 + 家访
	// 如果超时未完成，状态自动变为 expired
	ApprovedAt *time.Time `gorm:"comment:审核通过时间" json:"approved_at"`
	ExpiresAt  *time.Time `gorm:"index;comment:过期时间（通过后7天）" json:"expires_at"`

	// ==================== 完成信息 ====================
	// 如果领养完成，关联到 adoptions 表
	AdoptionID uint64 `gorm:"default:0;index;comment:关联的领养记录ID，0表示未完成" json:"adoption_id"`

	// ==================== 通用字段 ====================
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (AdoptApplication) TableName() string {
	return "adopt_applications"
}

// ==================== 查询过滤器 ====================

type AdoptApplicationFilter struct {
	CatID       uint64 // 按猫咪筛选（管理员用）
	ApplicantID uint64 // 按申请人筛选（我的申请）
	ReviewerID  uint64 // 按审核人筛选

	Status string // 按状态筛选：pending/approved/rejected/cancelled/completed/expired

	// 时间范围
	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time

	Page     int
	PageSize int

	SortBy    string // created_at/approved_at
	SortOrder string // asc/desc
}

func (r *repository) CreateApply(ctx context.Context, apply *AdoptApplication, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	err := tx.Create(apply).Error
	return err
}

func (r *repository) GetApplyByID(ctx context.Context, applyID uint64, txs ...*gorm.DB) (*AdoptApplication, error) {
	tx := r.getDB(ctx, txs...)
	var apply AdoptApplication
	err := tx.First(&apply, applyID).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repository) GetApplyByCatAndUser(ctx context.Context, catID, userID uint64, excludeStatuses []string, txs ...*gorm.DB) (*AdoptApplication, error) {
	tx := r.getDB(ctx, txs...)
	query := tx.Where("cat_id = ? AND applicant_id = ?", catID, userID)

	if len(excludeStatuses) > 0 {
		query = query.Where("status NOT IN ?", excludeStatuses)
	}

	var apply AdoptApplication
	err := query.Order("id DESC").First(&apply).Error
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func (r *repository) ListApplies(ctx context.Context, filter *AdoptApplicationFilter, txs ...*gorm.DB) ([]*AdoptApplication, int64, error) {
	tx := r.getDB(ctx, txs...)
	var total int64

	query := tx.Model(&AdoptApplication{})

	if filter.CatID > 0 {
		query = query.Where("cat_id = ?", filter.CatID)
	}
	if filter.ApplicantID > 0 {
		query = query.Where("applicant_id = ?", filter.ApplicantID)
	}
	if filter.ReviewerID > 0 {
		query = query.Where("reviewer_id = ?", filter.ReviewerID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
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

	var applies []*AdoptApplication
	sortField := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "approved_at":
			sortField = "approved_at"
		case "updated_at":
			sortField = "updated_at"
		}
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" || filter.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	err := query.Order(sortField + " " + sortOrder).Offset(offset).Limit(pageSize).Find(&applies).Error
	return applies, total, err
}

func (r *repository) UpdateApply(ctx context.Context, applyID uint64, values map[string]any, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	err := tx.Model(&AdoptApplication{}).Where("id = ?", applyID).Updates(values).Error
	return err
}

func (r *repository) CancelApply(ctx context.Context, applyID uint64, applicantUserID uint64, reason string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	// 使用事务确保数据一致性
	return tx.Transaction(func(tx *gorm.DB) error {
		// 锁住申请记录
		var apply AdoptApplication
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			return err
		}

		// 校验归属权
		if apply.ApplicantID != applicantUserID {
			return ErrApplyNotBelongToUser
		}

		// 校验状态
		if apply.Status != AdoptApplyStatusPending {
			return ErrApplyCannotCancel
		}

		now := time.Now()
		apply.Status = AdoptApplyStatusCancelled
		apply.RejectReason = reason // 复用 reject_reason 字段存储取消原因
		apply.ReviewedAt = &now

		return tx.Save(&apply).Error
	})
}

func (r *repository) RejectApply(ctx context.Context, applyID uint64, reviewerID uint64, reason string, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 加锁防止并发审核
		var apply AdoptApplication
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			return err
		}

		if apply.Status != AdoptApplyStatusPending {
			return ErrApplyAlreadyReviewed
		}

		now := time.Now()
		apply.Status = AdoptApplyStatusRejected
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = &now
		apply.RejectReason = reason

		return tx.Save(&apply).Error
	})
}

func (r *repository) ApproveApply(ctx context.Context, applyID uint64, reviewerID uint64, expiryDays int, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	return tx.Transaction(func(tx *gorm.DB) error {
		var apply AdoptApplication
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&apply, applyID).Error; err != nil {
			return err
		}

		if apply.Status != AdoptApplyStatusPending {
			return ErrApplyAlreadyReviewed
		}

		now := time.Now()
		expiresAt := now.AddDate(0, 0, expiryDays)

		apply.Status = AdoptApplyStatusApproved
		apply.ReviewerID = reviewerID
		apply.ReviewedAt = &now
		apply.ApprovedAt = &now
		apply.ExpiresAt = &expiresAt

		return tx.Save(&apply).Error
	})
}

func (r *repository) MarkExpired(ctx context.Context, applyID uint64, updatedAt time.Time, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	result := tx.Model(&AdoptApplication{}).Where("id = ? AND status = ?", applyID, AdoptApplyStatusApproved).Updates(map[string]interface{}{
		"status":     AdoptApplyStatusExpired,
		"updated_at": updatedAt,
	})

	if result.RowsAffected == 0 {
		return errors.New("申请不存在或状态不是 approved")
	}
	return result.Error
}

func (r *repository) UpdateStatus(ctx context.Context, applyID uint64, status string, additionalData map[string]any, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)

	data := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	for k, v := range additionalData {
		data[k] = v
	}

	return tx.Model(&AdoptApplication{}).Where("id = ?", applyID).Updates(data).Error
}

func (r *repository) DeleteApply(ctx context.Context, applyID uint64, txs ...*gorm.DB) error {
	tx := r.getDB(ctx, txs...)
	err := tx.Delete(&AdoptApplication{}, applyID).Error
	return err
}

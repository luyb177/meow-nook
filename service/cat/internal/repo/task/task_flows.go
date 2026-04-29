package task

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// TaskFlow 任务流转日志表（记录任务全生命周期操作）
type TaskFlow struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	TaskID uint64 `gorm:"not null;index;comment:关联任务ID" json:"task_id"`
	UserID uint64 `gorm:"not null;index;comment:操作人ID" json:"user_id"`

	// 操作信息
	Action     string `gorm:"type:varchar(50);not null;comment:操作类型:create/approve/reject/claim/abandon/complete/cancel/escalate" json:"action"`
	FromStatus string `gorm:"type:varchar(30);comment:原状态" json:"from_status"`
	ToStatus   string `gorm:"type:varchar(30);comment:新状态" json:"to_status"`
	Remark     string `gorm:"type:text;comment:操作备注" json:"remark"`

	// 扩展信息
	OldValues string `gorm:"type:text;comment:变更前值（JSON）" json:"old_values"`
	NewValues string `gorm:"type:text;comment:变更后值（JSON）" json:"new_values"`

	// 通用字段
	CreatedAt time.Time             `gorm:"autoCreateTime;comment:操作时间" json:"created_at"`
	UpdatedAt time.Time             `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `gorm:"index;softDelete:nano" json:"deleted_at,omitempty"`
}

func (TaskFlow) TableName() string {
	return "task_flows"
}

// ==================== 任务流转相关实现 ====================

func (r *repository) CreateTaskFlow(ctx context.Context, flow *TaskFlow, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(flow).Error
}

func (r *repository) ListTaskFlows(ctx context.Context, taskID uint64, tx ...*gorm.DB) ([]*TaskFlow, error) {
	var flows []*TaskFlow
	err := r.getDB(ctx, tx...).Where("task_id = ?", taskID).Order("created_at DESC").Find(&flows).Error
	return flows, err
}

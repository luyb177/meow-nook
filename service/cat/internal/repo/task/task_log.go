package task

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TaskLog struct {
	ID uint64 ` gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	TaskID uint64 `gorm:"not null;index" json:"task_id"`
	UserID uint64 `gorm:"not null;index" json:"user_id"`

	// 操作信息
	Action string `gorm:"type:varchar(50);not null;index" json:"action"`
	// action: created, claimed, completed, abandoned, approved, rejected, cancelled, status_changed
	Content   string `gorm:"type:text" json:"content"`
	ExtraData string `gorm:"type:text" json:"extra_data"` // JSON 格式，存储额外信息

	// 通用字段
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TaskLog) TableName() string {
	return "task_logs"
}

// ==================== 任务日志相关实现 ====================

func (r *repository) CreateTaskLog(ctx context.Context, log *TaskLog, tx ...*gorm.DB) error {
	return r.getDB(ctx, tx...).Create(log).Error
}

func (r *repository) ListTaskLogs(ctx context.Context, taskID uint64, tx ...*gorm.DB) ([]*TaskLog, error) {
	var logs []*TaskLog
	err := r.getDB(ctx, tx...).Where("task_id = ?", taskID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

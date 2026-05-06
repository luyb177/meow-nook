package user

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:254;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"password_hash"`
	Username     string    `gorm:"size:64" json:"username"`
	Avatar       string    `gorm:"size:512" json:"avatar"`
	Phone        string    `gorm:"size:32" json:"phone"`
	Area         string    `gorm:"size:128" json:"area"`
	Gender       string    `gorm:"size:16" json:"gender"`
	ServiceTypes []string  `gorm:"type:json;serializer:json" json:"service_types"`
	Points       int32     `gorm:"default:200" json:"points"`
	Role         string    `gorm:"size:32;default:'user'" json:"role"` // user / volunteer / admin
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

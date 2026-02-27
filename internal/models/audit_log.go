package models

import "time"

type AuditLog struct {
	ID uint `gorm:"primaryKey"`

	UserID *string `gorm:"type:uuid;index"`
	Event  string  `gorm:"index;not null"`

	Email string `gorm:"index"`
	IP    string
	UA    string

	CreatedAt time.Time
}

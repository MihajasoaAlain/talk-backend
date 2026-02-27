package models

import "time"

type RefreshToken struct {
	ID uint `gorm:"primaryKey"`

	UserID string `gorm:"type:uuid;index;not null"`
	User   User   `gorm:"constraint:OnDelete:CASCADE;"`

	TokenHash string `gorm:"uniqueIndex;not null"`

	ExpiresAt time.Time `gorm:"index;not null"`
	RevokedAt *time.Time

	ReplacedByID *uint

	CreatedAt time.Time
}

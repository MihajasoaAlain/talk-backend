package models

import "time"

type Conversation struct {
	ID        uint `gorm:"primaryKey"`
	IsGroup   bool `gorm:"not null;default:false"`
	Title     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`

	Password string `json:"password" gorm:"not null"`

	FailedLoginAttempts int        `json:"-" gorm:"not null;default:0"`
	LockedUntil         *time.Time `json:"-"`

	LastLoginAt *time.Time `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

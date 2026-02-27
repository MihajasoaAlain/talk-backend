package models

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"type:uuid;primaryKey"`
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
	if u.ID == "" {
		id, idErr := newUUIDv4()
		if idErr != nil {
			return idErr
		}
		u.ID = id
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func newUUIDv4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	), nil
}

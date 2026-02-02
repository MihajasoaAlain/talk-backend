package repository

import (
	"errors"
	"time"

	"talk-backend/internal/models"

	"gorm.io/gorm"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository interface {
	Create(rt *models.RefreshToken) error
	FindValidByHash(hash string) (*models.RefreshToken, error)
	Revoke(rt *models.RefreshToken, when time.Time) error
	Update(rt *models.RefreshToken) error
}

type refreshTokenRepository struct{ db *gorm.DB }

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(rt *models.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *refreshTokenRepository) FindValidByHash(hash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Preload("User").
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > NOW()", hash).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepository) Revoke(rt *models.RefreshToken, when time.Time) error {
	rt.RevokedAt = &when
	return r.db.Save(rt).Error
}

func (r *refreshTokenRepository) Update(rt *models.RefreshToken) error {
	return r.db.Save(rt).Error
}

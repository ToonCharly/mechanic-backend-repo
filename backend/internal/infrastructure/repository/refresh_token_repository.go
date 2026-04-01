package repository

import (
	"mechanic-backend/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a new refresh token in the database
func (r *RefreshTokenRepository) Create(token *entities.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindByToken retrieves a refresh token by its value
func (r *RefreshTokenRepository) FindByToken(token string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	return &refreshToken, err
}

// FindByUserID retrieves all active refresh tokens for a user
func (r *RefreshTokenRepository) FindByUserID(userID uuid.UUID) ([]entities.RefreshToken, error) {
	var tokens []entities.RefreshToken
	err := r.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Find(&tokens).Error
	return tokens, err
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(tokenID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entities.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", now).Error
}

// RevokeAllForUser revokes all active tokens for a user (logout from all devices)
func (r *RefreshTokenRepository) RevokeAllForUser(userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&entities.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

// DeleteExpired removes expired tokens from the database (cleanup job)
func (r *RefreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&entities.RefreshToken{}).Error
}

package usecases

import (
	"errors"
	"mechanic-backend/internal/application/dto"
	"mechanic-backend/internal/config"
	"mechanic-backend/internal/domain/entities"
	"mechanic-backend/internal/domain/repositories"
	"mechanic-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo RefreshTokenRepository
	cfg              *config.Config
}

type RefreshTokenRepository interface {
	Create(token *entities.RefreshToken) error
	FindByToken(token string) (*entities.RefreshToken, error)
	FindByUserID(userID uuid.UUID) ([]entities.RefreshToken, error)
	Revoke(tokenID uuid.UUID) error
	RevokeAllForUser(userID uuid.UUID) error
	DeleteExpired() error
}

func NewAuthUseCase(userRepo repositories.UserRepository, refreshTokenRepo RefreshTokenRepository, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
	}
}

func (uc *AuthUseCase) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := uc.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token pair
	return uc.generateTokenPair(user)
}

func (uc *AuthUseCase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate token pair
	return uc.generateTokenPair(user)
}

func (uc *AuthUseCase) RefreshAccessToken(req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Find refresh token in database
	storedToken, err := uc.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Validate token
	if !storedToken.IsValid() {
		return nil, errors.New("refresh token expired or revoked")
	}

	// Find user
	user, err := uc.userRepo.FindByID(storedToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token (15 minutes)
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Name,
		string(user.Role),
		uc.cfg.JWT.Secret,
	)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	plainRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	if err := uc.refreshTokenRepo.Revoke(storedToken.ID); err != nil {
		return nil, err
	}

	// Store new refresh token (no hashing - already secure random token)
	newRefreshToken := &entities.RefreshToken{
		UserID:    user.ID,
		Token:     plainRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := uc.refreshTokenRepo.Create(newRefreshToken); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: plainRefreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

func (uc *AuthUseCase) Logout(userID uuid.UUID, refreshToken string) error {
	// Find and revoke the specific refresh token
	storedToken, err := uc.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return nil // Token already invalid, silently succeed
	}

	if storedToken.UserID != userID {
		return errors.New("unauthorized")
	}

	return uc.refreshTokenRepo.Revoke(storedToken.ID)
}

func (uc *AuthUseCase) LogoutAllDevices(userID uuid.UUID) error {
	return uc.refreshTokenRepo.RevokeAllForUser(userID)
}

// generateTokenPair creates both access and refresh tokens
func (uc *AuthUseCase) generateTokenPair(user *entities.User) (*dto.AuthResponse, error) {
	// Generate access token (15 minutes)
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Name,
		string(user.Role),
		uc.cfg.JWT.Secret,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (7 days)
	plainRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token in database (no hashing - already secure random token)
	refreshToken := &entities.RefreshToken{
		UserID:    user.ID,
		Token:     plainRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := uc.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: plainRefreshToken,
		ExpiresIn:    900, // 15 minutes in seconds
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (uc *AuthUseCase) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return uc.userRepo.FindByID(id)
}

package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"messenger-vkedu/internal/user"
)

type UserRepository interface {
	Create(ctx context.Context, username, passwordHash string) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	users         UserRepository
	refreshTokens RefreshTokenRepository
	tokenManager  *TokenManager
}

func NewService(users UserRepository, refreshTokens RefreshTokenRepository, tokenManager *TokenManager) *Service {
	return &Service{
		users:         users,
		refreshTokens: refreshTokens,
		tokenManager:  tokenManager,
	}
}

func (service *Service) Register(ctx context.Context, username, password string) (*user.User, error) {
	passwordHash, err := GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	return service.users.Create(ctx, username, passwordHash)
}

func (service *Service) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	foundUser, err := service.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrWrongLoginOrPassword
		}
		return nil, err
	}

	if err := CheckPassword(foundUser.PasswordHash, password); err != nil {
		return nil, err
	}

	return service.issueTokenPair(ctx, foundUser.ID)
}

func (service *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	storedToken, err := service.refreshTokens.GetByHash(ctx, HashRefreshToken(refreshToken))
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	if storedToken.Revoked || time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	if err := service.refreshTokens.Revoke(ctx, storedToken.ID); err != nil {
		if errors.Is(err, ErrTokenRevoked) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	return service.issueTokenPair(ctx, storedToken.UserID)
}

func (service *Service) Logout(ctx context.Context, refreshToken string) error {
	storedToken, err := service.refreshTokens.GetByHash(ctx, HashRefreshToken(refreshToken))
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil
		}
		return err
	}

	err = service.refreshTokens.Revoke(ctx, storedToken.ID)
	if err != nil && !errors.Is(err, ErrTokenRevoked) {
		return err
	}
	return nil
}

func (service *Service) issueTokenPair(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	accessToken, err := service.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, tokenHash, expiresAt := service.tokenManager.GenerateRefreshToken()
	if err := service.refreshTokens.Create(ctx, userID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const refreshTokenSize = 32

type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (manager *TokenManager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(manager.accessTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(manager.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signedToken, nil
}

func (manager *TokenManager) ParseAccessToken(accessToken string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(
		accessToken,
		&claims,
		manager.signingKey,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidAccessToken, err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid subject", ErrInvalidAccessToken)
	}
	return userID, nil
}

func (manager *TokenManager) signingKey(token *jwt.Token) (any, error) {
	return manager.secret, nil
}

func (manager *TokenManager) GenerateRefreshToken() (refreshToken, tokenHash string, expiresAt time.Time) {
	randomBytes := make([]byte, refreshTokenSize)
	rand.Read(randomBytes)

	refreshToken = base64.RawURLEncoding.EncodeToString(randomBytes)
	return refreshToken, HashRefreshToken(refreshToken), time.Now().Add(manager.refreshTTL)
}

func HashRefreshToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}

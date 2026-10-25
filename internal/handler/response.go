package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"messenger-vkedu/internal/auth"
	"messenger-vkedu/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type errorResponse struct {
	Error string `json:"error"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type tokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newUserResponse(createdUser *user.User) userResponse {
	return userResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		CreatedAt: createdUser.CreatedAt,
	}
}

func newTokenPairResponse(tokenPair *auth.TokenPair) tokenPairResponse {
	return tokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}

func respondError(ginCtx *gin.Context, statusCode int, message string) {
	ginCtx.JSON(statusCode, errorResponse{Error: message})
}

func handleServiceError(ginCtx *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrUsernameTaken):
		respondError(ginCtx, http.StatusConflict, "username already taken")
	case errors.Is(err, auth.ErrPasswordTooLong):
		respondError(ginCtx, http.StatusBadRequest, "password must be at most 72 bytes")
	case errors.Is(err, auth.ErrWrongLoginOrPassword):
		respondError(ginCtx, http.StatusUnauthorized, "wrong login or password")
	case errors.Is(err, auth.ErrInvalidRefreshToken):
		respondError(ginCtx, http.StatusUnauthorized, "invalid refresh token")
	default:
		log.Printf("LOG - internal error: %v", err)
		respondError(ginCtx, http.StatusInternalServerError, "internal server error")
	}
}

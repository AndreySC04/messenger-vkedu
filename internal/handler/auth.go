package handler

import (
	"context"
	"net/http"

	"messenger-vkedu/internal/auth"
	"messenger-vkedu/internal/middleware"
	"messenger-vkedu/internal/user"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (*user.User, error)
	Login(ctx context.Context, username, password string) (*auth.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*auth.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (handler *AuthHandler) Register(ginCtx *gin.Context) {
	var request registerRequest
	if err := ginCtx.ShouldBindJSON(&request); err != nil {
		respondError(ginCtx, http.StatusBadRequest, "username must be 3-50 characters, password at least 8 characters")
		return
	}

	createdUser, err := handler.authService.Register(ginCtx.Request.Context(), request.Username, request.Password)
	if err != nil {
		handleServiceError(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, newUserResponse(createdUser))
}

func (handler *AuthHandler) Login(ginCtx *gin.Context) {
	var request loginRequest
	if err := ginCtx.ShouldBindJSON(&request); err != nil {
		respondError(ginCtx, http.StatusBadRequest, "username and password are required")
		return
	}

	tokenPair, err := handler.authService.Login(ginCtx.Request.Context(), request.Username, request.Password)
	if err != nil {
		handleServiceError(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, newTokenPairResponse(tokenPair))
}

func (handler *AuthHandler) Refresh(ginCtx *gin.Context) {
	var request refreshTokenRequest
	if err := ginCtx.ShouldBindJSON(&request); err != nil {
		respondError(ginCtx, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokenPair, err := handler.authService.Refresh(ginCtx.Request.Context(), request.RefreshToken)
	if err != nil {
		handleServiceError(ginCtx, err)
		return
	}

	ginCtx.JSON(http.StatusOK, newTokenPairResponse(tokenPair))
}

func (handler *AuthHandler) Logout(ginCtx *gin.Context) {
	var request refreshTokenRequest
	if err := ginCtx.ShouldBindJSON(&request); err != nil {
		respondError(ginCtx, http.StatusBadRequest, "refresh_token is required")
		return
	}

	if err := handler.authService.Logout(ginCtx.Request.Context(), request.RefreshToken); err != nil {
		handleServiceError(ginCtx, err)
		return
	}

	ginCtx.Status(http.StatusNoContent)
}

func (handler *AuthHandler) Me(ginCtx *gin.Context) {
	userID, ok := middleware.UserIDFromContext(ginCtx)
	if !ok {
		respondError(ginCtx, http.StatusUnauthorized, "unauthorized")
		return
	}

	ginCtx.JSON(http.StatusOK, gin.H{"user_id": userID})
}

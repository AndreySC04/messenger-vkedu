package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	userIDContextKey = "user_id"
	bearerPrefix     = "Bearer "
)

type AccessTokenParser interface {
	ParseAccessToken(accessToken string) (uuid.UUID, error)
}

func RequireAuth(tokenParser AccessTokenParser) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		authorizationHeader := ginCtx.GetHeader("Authorization")
		if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
			abortUnauthorized(ginCtx, "missing or malformed authorization header")
			return
		}

		accessToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)
		userID, err := tokenParser.ParseAccessToken(accessToken)
		if err != nil {
			abortUnauthorized(ginCtx, "invalid or expired access token")
			return
		}

		ginCtx.Set(userIDContextKey, userID)
		ginCtx.Next()
	}
}

func UserIDFromContext(ginCtx *gin.Context) (uuid.UUID, bool) {
	value, exists := ginCtx.Get(userIDContextKey)
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	return userID, ok
}

func abortUnauthorized(ginCtx *gin.Context, message string) {
	ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}

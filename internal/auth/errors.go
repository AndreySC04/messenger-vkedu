package auth

import "errors"

var (
	ErrTokenNotFound        = errors.New("refresh token not found")
	ErrTokenRevoked         = errors.New("refresh token already revoked")
	ErrPasswordTooLong      = errors.New("password too long")
	ErrWrongLoginOrPassword = errors.New("wrong login or password")
	ErrHashTooShort         = errors.New("hash is broken")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

package user

import "errors"

var (
	ErrUserNotFound  = errors.New("ERR - user not found")
	ErrUsernameTaken = errors.New("ERR - username already taken")
)

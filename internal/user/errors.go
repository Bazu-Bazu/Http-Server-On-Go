package user

import "errors"

var (
	ErrInvalidInput  = errors.New("Invalid input")
	ErrNotFound      = errors.New("User not found")
	ErrNicknameTaken = errors.New("Nickname already exists")
)

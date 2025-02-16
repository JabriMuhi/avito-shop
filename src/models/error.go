package models

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with that name already exists")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrEmptyName         = errors.New("empty name")
	ErrNoSuchUser        = errors.New("no such user")
)

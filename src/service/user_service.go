//go:generate mockgen -source=./user_service.go -destination=./mocks/mockgen_user_storage.go -package=mock_service
package service

import (
	"avito-shop/src/models"
	"context"
	"errors"
)

type userGetter interface {
	GetUserByName(ctx context.Context, name string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

type UserService struct {
	userStorage userGetter
}

func NewUserService(userStorage userGetter) *UserService {
	return &UserService{userStorage}
}

func (s *UserService) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	if name == "" {
		return nil, models.ErrEmptyName
	}
	return s.userStorage.GetUserByName(ctx, name)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("id is empty")
	}
	return s.userStorage.GetUserByID(ctx, id)
}

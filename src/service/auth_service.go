package service

import (
	"avito-shop/avito"
	"avito-shop/src/auth"
	"avito-shop/src/models"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userStorage interface {
	CreateUser(ctx context.Context, name, password string) (*models.User, error)
	GetUserIDAndPasswordByName(ctx context.Context, name string) (int64, string, error)
	IsUserExist(ctx context.Context, name string) (bool, error)
}

type AuthServiceImpl struct {
	userRepo userStorage
}

func NewAuthService(userRepo userStorage) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo}
}

func (s *AuthServiceImpl) Authenticate(ctx context.Context, req *avito.AuthRequest) (*avito.AuthResponse, error) {
	exist, err := s.userRepo.IsUserExist(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !exist {
		_, err := s.userRepo.CreateUser(ctx, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, models.ErrUserAlreadyExists) {
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}

			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	userID, hashedPassword, err := s.userRepo.GetUserIDAndPasswordByName(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
	}

	token, err := auth.GenerateJWT(req.Username, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &avito.AuthResponse{Token: token}, nil
}

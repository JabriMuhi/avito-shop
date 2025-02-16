package service

import (
	"avito-shop/src/models"
	"avito-shop/src/service/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_service.NewMockuserGetter(ctrl)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := int64(1)
	expectedUser := &models.User{ID: userID, Name: "testuser"}

	t.Run("Successful retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)

		user, err := userService.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("Error during retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, errors.New("database error"))

		_, err := userService.GetUserByID(ctx, userID)
		assert.Error(t, err)
	})
}

func TestUserService_GetUserByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_service.NewMockuserGetter(ctrl)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	name := "testuser"
	expectedUser := &models.User{ID: 1, Name: name}

	t.Run("Successful retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByName(ctx, name).Return(expectedUser, nil)

		user, err := userService.GetUserByName(ctx, name)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("Error during retrieval", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByName(ctx, name).Return(nil, errors.New("database error"))

		_, err := userService.GetUserByName(ctx, name)
		assert.Error(t, err)
	})

	t.Run("Retrieval with empty name", func(t *testing.T) {
		_, err := userService.GetUserByName(ctx, "")
		assert.Error(t, err)
	})
}

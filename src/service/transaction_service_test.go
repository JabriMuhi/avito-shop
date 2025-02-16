package service

import (
	"avito-shop/src/service/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_TransferCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_service.NewMocktransactionStorage(ctrl)
	transactionService := NewTransactionService(mockTransactionRepo)

	ctx := context.Background()
	senderID := int64(1)
	receiverID := int64(2)
	amount := int64(100)

	t.Run("Successful transfer", func(t *testing.T) {
		mockTransactionRepo.EXPECT().TransferCoins(ctx, senderID, receiverID, amount).Return(nil)

		err := transactionService.TransferCoins(ctx, senderID, receiverID, amount)
		assert.NoError(t, err)
	})

	t.Run("Error during transfer", func(t *testing.T) {
		mockTransactionRepo.EXPECT().TransferCoins(ctx, senderID, receiverID, amount).Return(errors.New("database error"))

		err := transactionService.TransferCoins(ctx, senderID, receiverID, amount)
		assert.Error(t, err)
	})

	t.Run("Transfer with zero amount", func(t *testing.T) {
		err := transactionService.TransferCoins(ctx, senderID, receiverID, 0)
		assert.Error(t, err)
	})
}

func TestTransactionService_GetCoinTransactionsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_service.NewMocktransactionStorage(ctrl)
	transactionService := NewTransactionService(mockTransactionRepo)

	ctx := context.Background()
	userID := int64(1)
	expectedSent := map[string]int64{"user2": 100}
	expectedReceived := map[string]int64{"user3": 50}

	t.Run("Successful retrieval", func(t *testing.T) {
		mockTransactionRepo.EXPECT().GetCoinTransactionsByUserID(ctx, userID).Return(expectedSent, expectedReceived, nil)

		sent, received, err := transactionService.GetCoinTransactionsByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSent, sent)
		assert.Equal(t, expectedReceived, received)
	})

	t.Run("Error during retrieval", func(t *testing.T) {
		mockTransactionRepo.EXPECT().GetCoinTransactionsByUserID(ctx, userID).Return(nil, nil, errors.New("database error"))

		_, _, err := transactionService.GetCoinTransactionsByUserID(ctx, userID)
		assert.Error(t, err)
	})

	t.Run("Retrieval with invalid user ID", func(t *testing.T) {
		mockTransactionRepo.EXPECT().GetCoinTransactionsByUserID(ctx, int64(0)).Return(nil, nil, errors.New("invalid user ID"))

		_, _, err := transactionService.GetCoinTransactionsByUserID(ctx, int64(0))
		assert.Error(t, err)
	})
}

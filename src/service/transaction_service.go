//go:generate mockgen -source=./transaction_service.go -destination=./mocks/mockgen_transaction_storage.go -package=mock_service
package service

import (
	"avito-shop/src/models"
	"context"
)

type transactionStorage interface {
	TransferCoins(ctx context.Context, senderID, receiverID, amount int64) error
	GetCoinTransactionsByUserID(ctx context.Context, userID int64) (map[string]int64, map[string]int64, error)
}

type TransactionService struct {
	transactionStorage transactionStorage
}

func NewTransactionService(ts transactionStorage) *TransactionService {
	return &TransactionService{ts}
}

func (s *TransactionService) TransferCoins(ctx context.Context, senderID, receiverID, amount int64) error {
	if amount < 0 || amount == 0 {
		return models.ErrInvalidAmount
	}
	return s.transactionStorage.TransferCoins(ctx, senderID, receiverID, amount)
}

func (s *TransactionService) GetCoinTransactionsByUserID(ctx context.Context, userID int64) (map[string]int64, map[string]int64, error) {
	return s.transactionStorage.GetCoinTransactionsByUserID(ctx, userID)
}

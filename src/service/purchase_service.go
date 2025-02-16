//go:generate mockgen -source=./purchase_service.go -destination=./mocks/mockgen_purchase_storage.go -package=mock_service
package service

import (
	"avito-shop/avito"
	"context"
	"errors"
)

type purchaseStorage interface {
	BuyMerch(ctx context.Context, userID int64, merchName string) error
	GetPurchasedMerchByUserID(ctx context.Context, userID int64) ([]*avito.InventoryItem, error)
}

type PurchaseService struct {
	purchaseService purchaseStorage
}

func NewPurchaseService(purchaseService purchaseStorage) *PurchaseService {
	return &PurchaseService{purchaseService: purchaseService}
}

func (s *PurchaseService) BuyMerch(ctx context.Context, userID int64, merchName string) error {
	if merchName == "" {
		return errors.New("merchName is empty")
	}
	return s.purchaseService.BuyMerch(ctx, userID, merchName)
}

func (s *PurchaseService) GetPurchasedMerchByUserID(ctx context.Context, userID int64) ([]*avito.InventoryItem, error) {
	return s.purchaseService.GetPurchasedMerchByUserID(ctx, userID)
}

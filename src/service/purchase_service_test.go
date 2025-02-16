package service

import (
	"avito-shop/avito"
	"avito-shop/src/service/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseService_BuyMerch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPurchaseRepo := mock_service.NewMockpurchaseStorage(ctrl)
	purchaseService := NewPurchaseService(mockPurchaseRepo)

	ctx := context.Background()
	userID := int64(1)
	merchName := "t-shirt"

	t.Run("Successful purchase", func(t *testing.T) {
		mockPurchaseRepo.EXPECT().BuyMerch(ctx, userID, merchName).Return(nil)

		err := purchaseService.BuyMerch(ctx, userID, merchName)
		assert.NoError(t, err)
	})

	t.Run("Error during purchase", func(t *testing.T) {
		mockPurchaseRepo.EXPECT().BuyMerch(ctx, userID, merchName).Return(errors.New("database error"))

		err := purchaseService.BuyMerch(ctx, userID, merchName)
		assert.Error(t, err)
	})

	t.Run("Purchase with empty merch name", func(t *testing.T) {
		err := purchaseService.BuyMerch(ctx, userID, "")
		assert.Error(t, err)
	})
}

func TestPurchaseService_GetPurchasedMerchByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPurchaseRepo := mock_service.NewMockpurchaseStorage(ctrl)
	purchaseService := NewPurchaseService(mockPurchaseRepo)

	ctx := context.Background()
	userID := int64(1)
	expectedItems := []*avito.InventoryItem{
		{Type: "t-shirt", Quantity: 1},
	}

	t.Run("Successful retrieval", func(t *testing.T) {
		mockPurchaseRepo.EXPECT().GetPurchasedMerchByUserID(ctx, userID).Return(expectedItems, nil)

		items, err := purchaseService.GetPurchasedMerchByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedItems, items)
	})

	t.Run("Error during retrieval", func(t *testing.T) {
		mockPurchaseRepo.EXPECT().GetPurchasedMerchByUserID(ctx, userID).Return(nil, errors.New("database error"))

		_, err := purchaseService.GetPurchasedMerchByUserID(ctx, userID)
		assert.Error(t, err)
	})

	t.Run("Retrieval with invalid user ID", func(t *testing.T) {
		mockPurchaseRepo.EXPECT().GetPurchasedMerchByUserID(ctx, int64(0)).Return(nil, errors.New("invalid user ID"))

		_, err := purchaseService.GetPurchasedMerchByUserID(ctx, int64(0))
		assert.Error(t, err)
	})
}

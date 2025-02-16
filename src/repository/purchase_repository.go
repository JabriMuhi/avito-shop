package repository

import (
	"avito-shop/avito"
	"avito-shop/src/models"
	"avito-shop/src/repository/pg"
	"context"
	"database/sql"
	"errors"
)

type PurchaseRepo struct {
	txProvider txProvider
	merch      map[string]models.Merch
}

func NewPurchaseRepo(txProvider txProvider, merch map[string]models.Merch) *PurchaseRepo {
	return &PurchaseRepo{txProvider, merch}
}

func (r *PurchaseRepo) BuyMerch(ctx context.Context, userID int64, merchName string) error {
	return r.txProvider.InWriteTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		user, err := pg.GetAndLockUserByID(ctx, tx, userID)
		if err != nil {
			return err
		}

		merch, exists := r.merch[merchName]
		if !exists {
			return errors.New("merch not found")
		}

		if user.Balance < merch.Price {
			return errors.New("insufficient balance")
		}

		newBalance := user.Balance - merch.Price
		err = pg.UpdateUserBalance(ctx, tx, userID, newBalance)
		if err != nil {
			return err
		}

		return pg.BuyMerch(ctx, tx, userID, merchName)
	})
}

func (r *PurchaseRepo) GetPurchasedMerchByUserID(ctx context.Context, userID int64) ([]*avito.InventoryItem, error) {
	var items []*avito.InventoryItem
	var err error

	err = r.txProvider.InReadTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		items, err = pg.GetPurchasedMerchByUserID(ctx, tx, userID)

		return err
	})

	return items, err
}

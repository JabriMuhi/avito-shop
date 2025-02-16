package repository

import (
	"avito-shop/src/models"
	"avito-shop/src/repository/pg"
	"context"
	"database/sql"
)

type txProvider interface {
	InWriteTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error
	InReadTx(ctx context.Context, f func(ctx context.Context, tx *sql.Tx) error) error
}

type MerchRepository struct {
	txProvider txProvider
}

func NewMerchRepository(txProvider txProvider) *MerchRepository {
	return &MerchRepository{txProvider: txProvider}
}

func (r *MerchRepository) GetAllMerch() ([]models.Merch, error) {
	var merch []models.Merch
	var err error

	err = r.txProvider.InReadTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		merch, err = pg.GetAllMerch(ctx, tx)

		return err
	})

	return merch, err
}

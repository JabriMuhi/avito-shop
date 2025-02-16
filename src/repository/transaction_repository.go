package repository

import (
	"avito-shop/src/repository/pg"
	"context"
	"database/sql"
	"errors"
)

type TransactionRepo struct {
	txProvider txProvider
}

func NewTransactionRepo(txProvider txProvider) *TransactionRepo {
	return &TransactionRepo{txProvider}
}

func (t *TransactionRepo) TransferCoins(ctx context.Context, senderID, receiverID, amount int64) error {
	return t.txProvider.InWriteTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := pg.LockUsersForTransaction(ctx, tx, senderID, receiverID)
		if err != nil {
			return err
		}

		senderUser, err := pg.GetUserByID(ctx, tx, senderID)
		if err != nil {
			return err
		}

		receiverUser, err := pg.GetUserByID(ctx, tx, receiverID)
		if err != nil {
			return err
		}

		if senderUser.Balance < amount {
			return errors.New("insufficient funds")
		}

		err = pg.UpdateUserBalance(ctx, tx, senderID, senderUser.Balance-amount)
		if err != nil {
			return err
		}

		err = pg.UpdateUserBalance(ctx, tx, receiverID, receiverUser.Balance+amount)
		if err != nil {
			return err
		}

		return pg.AddTransaction(ctx, tx, senderID, receiverID, amount)
	})
}

func (t *TransactionRepo) GetCoinTransactionsByUserID(ctx context.Context, userID int64) (map[string]int64, map[string]int64, error) {
	var userSend map[string]int64
	var receivedToUser map[string]int64
	var err error

	err = t.txProvider.InReadTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		userSend, receivedToUser, err = pg.GetCoinTransactionsByUserID(ctx, tx, userID)

		return err
	})

	return userSend, receivedToUser, nil
}

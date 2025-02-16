package repository

import (
	"avito-shop/src/models"
	"avito-shop/src/repository/pg"
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

const defaultUserBalance = 1000

type UserRepo struct {
	txProvider txProvider
}

func NewUserRepo(txProvider txProvider) *UserRepo {
	return &UserRepo{txProvider}
}

func (r *UserRepo) CreateUser(ctx context.Context, name, password string) (*models.User, error) {
	var user *models.User

	// Хэшшшшшшшшшшшшшшшшшшшшшшшшшшшшш хэшулямбус
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = r.txProvider.InWriteTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		user, err = pg.CreateUser(ctx, tx, name, hashedPassword, defaultUserBalance)

		return err
	})

	return user, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	var user *models.User
	var err error

	err = r.txProvider.InReadTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		user, err = pg.GetUserByName(ctx, tx, name)

		return err
	})

	return user, err
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user *models.User
	var err error

	err = r.txProvider.InReadTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		user, err = pg.GetUserByID(ctx, tx, id)

		return err
	})

	return user, err
}

func (r *UserRepo) GetUserIDAndPasswordByName(ctx context.Context, name string) (int64, string, error) {
	var userID int64
	var hashedPassword string
	var err error

	err = r.txProvider.InWriteTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		userID, hashedPassword, err = pg.GetUserIDAndPasswordByName(ctx, tx, name)

		return err
	})

	return userID, hashedPassword, nil
}

func (r *UserRepo) IsUserExist(ctx context.Context, name string) (bool, error) {
	var exist bool
	var err error

	err = r.txProvider.InWriteTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		exist, err = pg.IsUserExist(ctx, tx, name)

		return err
	})

	return exist, err

}

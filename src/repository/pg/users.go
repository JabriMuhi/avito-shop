package pg

import (
	"avito-shop/src/models"
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateUser(_ context.Context, tx *sql.Tx, name string, hashedPassword []byte, defaultUserBalance int) (*models.User, error) {
	var user models.User
	var exists bool

	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.ErrUserAlreadyExists
	}

	err = tx.QueryRow(`INSERT INTO users (name, password, balance)
		 VALUES ($1, $2, $3) RETURNING id, name, balance`,
		name, hashedPassword, defaultUserBalance).Scan(&user.ID, &user.Name, &user.Balance)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByName(_ context.Context, tx *sql.Tx, name string) (*models.User, error) {
	var user models.User

	err := tx.QueryRow("SELECT id, name, balance FROM users WHERE name=$1", name).Scan(&user.ID, &user.Name, &user.Balance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoSuchUser
		}

		return nil, err
	}

	return &user, nil
}

func GetAndLockUserByID(_ context.Context, tx *sql.Tx, id int64) (*models.User, error) {
	var user models.User

	err := tx.QueryRow("SELECT id, name, balance FROM users WHERE id=$1 FOR UPDATE", id).Scan(&user.ID, &user.Name, &user.Balance)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(_ context.Context, tx *sql.Tx, id int64) (*models.User, error) {
	var user models.User

	err := tx.QueryRow("SELECT id, name, balance FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Balance)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserBalance(_ context.Context, tx *sql.Tx, userID, newBalance int64) error {
	_, err := tx.Exec("UPDATE users SET balance = $1 WHERE id = $2", newBalance, userID)
	return err
}

func LockUsersForTransaction(_ context.Context, tx *sql.Tx, SenderID, ReceiverID int64) error {
	_, err := tx.Exec("SELECT * FROM users WHERE id = $1 OR id = $2 FOR UPDATE", ReceiverID, SenderID)
	return err
}

func GetUserIDAndPasswordByName(_ context.Context, tx *sql.Tx, name string) (int64, string, error) {
	var userID int64
	var hashedPassword string

	err := tx.QueryRow("SELECT id, password FROM users WHERE name=$1", name).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", status.Errorf(codes.Unauthenticated, "Invalid username or password")
		}
		return 0, "", status.Errorf(codes.Internal, "Internal server error")
	}

	return userID, hashedPassword, nil
}

func IsUserExist(_ context.Context, tx *sql.Tx, name string) (bool, error) {
	var exist bool

	err := tx.QueryRow("SELECT EXISTS(SELECT FROM users WHERE name=$1)", name).Scan(&exist)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	return exist, nil

}

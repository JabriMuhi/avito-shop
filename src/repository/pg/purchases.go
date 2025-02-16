package pg

import (
	"avito-shop/avito"
	"context"
	"database/sql"
	"time"
)

func BuyMerch(_ context.Context, tx *sql.Tx, userID int64, merchName string) error {
	_, err := tx.Exec("INSERT INTO purchases (user_id, merch_name, timestamp) VALUES ($1, $2, $3)", userID, merchName, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func GetPurchasedMerchByUserID(_ context.Context, tx *sql.Tx, userID int64) ([]*avito.InventoryItem, error) {
	rows, err := tx.Query("SELECT merch_name, COUNT(*) FROM purchases WHERE user_id=$1 GROUP BY merch_name", userID)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var inventory []*avito.InventoryItem

	for rows.Next() {
		var name string
		var count int32

		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		inventory = append(inventory, &avito.InventoryItem{Type: name, Quantity: count})
	}

	return inventory, nil
}

package pg

import (
	"context"
	"database/sql"
)

func AddTransaction(_ context.Context, tx *sql.Tx, senderID, receiverID, amount int64) error {
	_, err := tx.Exec("INSERT INTO transactions (sender_id, receiver_id, amount, timestamp) VALUES ($1, $2, $3, now())",
		senderID, receiverID, amount)

	if err != nil {
		return err
	}

	return nil
}

func GetCoinTransactionsByUserID(_ context.Context, tx *sql.Tx, userID int64) (map[string]int64, map[string]int64, error) {
	userSend := make(map[string]int64)
	receivedToUser := make(map[string]int64)

	query := `SELECT users.name, amount FROM transactions
	JOIN users ON users.id = transactions.sender_id
	WHERE receiver_id = $1`

	rowsReceivedToUserHistory, err := tx.Query(query, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rowsReceivedToUserHistory.Close()

	for rowsReceivedToUserHistory.Next() {
		var senderName string
		var amount int64

		if err := rowsReceivedToUserHistory.Scan(&senderName, &amount); err != nil {
			return nil, nil, err
		}

		receivedToUser[senderName] += amount
	}
	if err = rowsReceivedToUserHistory.Err(); err != nil {
		return nil, nil, err
	}

	query = `SELECT users.name, amount FROM transactions
	JOIN users ON users.id = transactions.receiver_id
	WHERE sender_id = $1`

	rowsSendUserHistory, err := tx.Query(query, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rowsSendUserHistory.Close()

	for rowsSendUserHistory.Next() {
		var receiverName string
		var amount int64

		if err := rowsSendUserHistory.Scan(&receiverName, &amount); err != nil {
			return nil, nil, err
		}

		userSend[receiverName] += amount
	}
	if err = rowsSendUserHistory.Err(); err != nil {
		return nil, nil, err
	}

	return userSend, receivedToUser, nil
}

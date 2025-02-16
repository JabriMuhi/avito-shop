package pg

import (
	"avito-shop/src/models"
	"context"
	"database/sql"
)

func GetAllMerch(_ context.Context, tx *sql.Tx) ([]models.Merch, error) {
	rows, err := tx.Query("SELECT * FROM merch")

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var merch []models.Merch
	for rows.Next() {
		var tmpMerch models.Merch
		if err := rows.Scan(&tmpMerch.ID, &tmpMerch.Name, &tmpMerch.Price); err != nil {
			return nil, err
		}
		merch = append(merch, tmpMerch)
	}

	return merch, nil
}

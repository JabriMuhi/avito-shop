package models

import "time"

type Purchase struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MerchName string    `json:"merch_name"`
	Timestamp time.Time `json:"timestamp"`
}

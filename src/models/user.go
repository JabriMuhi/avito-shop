package models

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Balance  int64  `json:"balance"`
}

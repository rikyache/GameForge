package models

import "time"

type Player struct {
	ID        int64     `db:"id"		json:"id"`
	Name      string    `db:"name"	json:"name"`
	Balance   int       `db:"balance"	json:"balance"`
	CreatedAt time.Time `db:"created_at"	json:"created_at"`
}

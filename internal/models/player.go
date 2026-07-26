package models

type Player struct {
	Id     int64  `db:"id"		json:"id"`
	Name   string `db:"name"	json:"name"`
	Health int    `db:"health"	json:"health"`
}

type PlayerResponse struct {
	Name   string `json:"name"`
	Health int    `json:"health"`
}

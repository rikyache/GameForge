package models

type PurchaseRequest struct {
	PlayerID int64 `json:"player_id"`
	GameID   int64 `json:"game_id"`
}

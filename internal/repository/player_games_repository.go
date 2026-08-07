package repository

import (
	"database/sql"
	"errors"
	"testsmth/internal/models"
)

type PlayerGameRepository struct {
	DB *sql.DB
}

func NewPlayerGameRepository(db *sql.DB) *PlayerGameRepository {
	return &PlayerGameRepository{
		DB: db,
	}
}

func (r *PlayerGameRepository) BuyGame(playerID int64, gameID int64) error {
	query := `
    INSERT INTO player_games(player_id, game_id)
    VALUES ($1, $2)
	`

	_, err := r.DB.Exec(query, playerID, gameID)

	return err
}

func (r *PlayerGameRepository) GetPlayerGames(playerID int64) ([]models.OwnedGame, error) {

	query := `
	SELECT
		g.id,
		g.name,
		g.genre,
		pg.bought_at
	
	FROM player_games pg
	JOIN games g
	ON g.id = pg.game_id
	
	WHERE pg.player_id = $1
	`

	rows, err := r.DB.Query(query, playerID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []models.OwnedGame

	for rows.Next() {

		var game models.OwnedGame

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Genre,
			&game.BoughtAt,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *PlayerGameRepository) RemoveGame(playerID int64, gameID int64) error {
	query := `
	DELETE FROM player_games
    WHERE player_id = $1
	AND game_id = $2
`

	result, err := r.DB.Exec(query, playerID, gameID)

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("game not found")
	}

	return nil
}

func (r *PlayerGameRepository) Exists(playerID int64, gameID int64) (bool, error) {
	query := `
    SELECT EXISTS (
		SELECT 1
		FROM player_games
		WHERE player_id = $1
		AND game_id = $2
)
`
	var exists bool

	err := r.DB.QueryRow(query, playerID, gameID).Scan(&exists)
	return exists, err
}

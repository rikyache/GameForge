package repository

import (
	"database/sql"
	"errors"
	"testsmth/internal/models"
)

type PlayerGameRepository struct {
	DB         *sql.DB
	PlayerRepo *PlayerRepository
	GameRepo   *GameRepository
}

func NewPlayerGameRepository(
	db *sql.DB,
	playerRepo *PlayerRepository,
	gameRepo *GameRepository,
) *PlayerGameRepository {
	return &PlayerGameRepository{
		DB:         db,
		PlayerRepo: playerRepo,
		GameRepo:   gameRepo,
	}
}

func (r *PlayerGameRepository) AddGame(playerID int64, gameID int64) error {
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

func (r *PlayerGameRepository) ExistsTx(tx *sql.Tx, playerID int64, gameID int64) (bool, error) {
	var exists bool

	err := tx.QueryRow(`
		SELECT EXISTS (
		    SELECT 1
		    FROM player_games
		    WHERE player_id = $1
		    AND game_id = $2
		)
	`, playerID, gameID).Scan(&exists)

	return exists, err
}

func (r *PlayerGameRepository) AddGameTx(tx *sql.Tx, playerID int64, gameID int64) error {
	_, err := tx.Exec(`
		INSERT INTO player_games(player_id, game_id)
		VALUES ($1, $2)
	`, playerID, gameID)

	return err
}

func (r *PlayerGameRepository) BuyGame(playerID int64, gameID int64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balance, err := r.PlayerRepo.GetBalance(tx, playerID)
	if err != nil {
		return err
	}

	price, err := r.GameRepo.GetPrice(tx, gameID)
	if err != nil {
		return err
	}

	if balance < price {
		return errors.New("not enough balance to buy")
	}

	exists, err := r.ExistsTx(tx, playerID, gameID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("game already owned")
	}

	err = r.PlayerRepo.UpdateBalance(tx, playerID, balance-price)
	if err != nil {
		return err
	}

	err = r.AddGameTx(tx, playerID, gameID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

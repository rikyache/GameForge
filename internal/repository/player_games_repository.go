package repository

import (
	"context"
	"database/sql"
	"errors"
	"testsmth/internal/apperrors"
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

func (r *PlayerGameRepository) GetPlayerGames(ctx context.Context, playerID int64) ([]models.OwnedGame, error) {

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

	rows, err := r.DB.QueryContext(ctx, query, playerID)
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

func (r *PlayerGameRepository) RemoveGame(ctx context.Context, playerID int64, gameID int64) error {
	query := `
	DELETE FROM player_games
    WHERE player_id = $1
	AND game_id = $2
`

	result, err := r.DB.ExecContext(ctx, query, playerID, gameID)
	if err != nil {
		return err
	}

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

func (r *PlayerGameRepository) ExistsTx(ctx context.Context, tx *sql.Tx, playerID int64, gameID int64) (bool, error) {
	var exists bool

	err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (
		    SELECT 1
		    FROM player_games
		    WHERE player_id = $1
		    AND game_id = $2
		)
	`, playerID, gameID).Scan(&exists)

	return exists, err
}

func (r *PlayerGameRepository) AddGameTx(ctx context.Context, tx *sql.Tx, playerID int64, gameID int64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO player_games(player_id, game_id)
		VALUES ($1, $2)
	`, playerID, gameID)

	return err
}

func (r *PlayerGameRepository) BuyGame(ctx context.Context, playerID int64, gameID int64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	balance, err := r.PlayerRepo.GetBalance(ctx, tx, playerID)
	if err != nil {
		return err
	}

	price, err := r.GameRepo.GetPrice(ctx, tx, gameID)
	if err != nil {
		return err
	}

	if balance < price {
		return errors.New("not enough balance to buy")
	}

	exists, err := r.ExistsTx(ctx, tx, playerID, gameID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("game already owned")
	}

	err = r.PlayerRepo.UpdateBalance(ctx, tx, playerID, balance-price)
	if err != nil {
		return err
	}

	err = r.AddGameTx(ctx, tx, playerID, gameID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (r *PlayerGameRepository) Refund(ctx context.Context, playerID int64, gameID int64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var price int

	err = tx.QueryRowContext(ctx, `
		SELECT g.price
		FROM player_games pg
		JOIN games g ON g.id = pg.game_id
		WHERE pg.player_id = $1
		  AND pg.game_id = $2
	`, playerID, gameID).Scan(&price)

	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.ErrGameNotFound
	}

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE players
		SET balance = balance + $1
		WHERE id = $2
	`, price, playerID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM player_games
		WHERE player_id = $1
		  AND game_id = $2
	`, playerID, gameID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

package repository

import (
	"database/sql"
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
)

type PlayerRepository struct {
	DB *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		DB: db,
	}
}

func (r *PlayerRepository) AddPlayer(name string) error {
	query := `
	INSERT INTO players (name)
	VALUES ($1)
`

	_, err := r.DB.Exec(query, name)

	return err
}

func (r *PlayerRepository) GetPlayer(id int64) (*models.Player, error) {
	query := `
		SELECT id, name, balance, created_at
		FROM players
		WHERE id = $1
	`

	var player models.Player

	err := r.DB.QueryRow(query, id).Scan(
		&player.ID,
		&player.Name,
		&player.Balance,
		&player.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *PlayerRepository) ListPlayers() ([]models.Player, error) {
	query := `
	SELECT id, name, balance, created_at 
	FROM players
	`

	var Players []models.Player

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player models.Player

		err := rows.Scan(
			&player.ID,
			&player.Name,
			&player.Balance,
			&player.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		Players = append(Players, player)
	}

	return Players, nil
}

func (r *PlayerRepository) RemovePlayer(id int64) error {

	query := `
		DELETE FROM players
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *PlayerRepository) GetBalance(tx *sql.Tx, playerID int64) (int, error) {
	var balance int

	err := tx.QueryRow(`
	SELECT balance 
	FROM players
	WHERE id = $1
	`, playerID).Scan(&balance)

	return balance, err
}

func (r *PlayerRepository) UpdateBalance(tx *sql.Tx, playerID int64, balance int) error {
	_, err := tx.Exec(`
	UPDATE players 
	SET balance = $1
	WHERE id = $2
	`, balance, playerID)

	return err
}

func (r *PlayerRepository) Deposit(playerID int64, amount int) error {
	query := `
    UPDATE players
    SET balance = balance + $1
    WHERE id = $2
    RETURNING id, balance;
`

	result, err := r.DB.Exec(query, amount, playerID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.ErrPlayerNotFound
	}

	return nil
}

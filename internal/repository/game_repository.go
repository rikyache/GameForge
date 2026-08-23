package repository

import (
	"database/sql"
	"testsmth/internal/models"
)

type GameRepository struct {
	DB *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{
		DB: db,
	}
}

func (r *GameRepository) AddGame(name string, genre string, price int) error {
	query := `
	INSERT INTO games (name, genre, price)
	VALUES ($1, $2, $3);
	`

	_, err := r.DB.Exec(query, name, genre, price)

	return err
}

func (r *GameRepository) GetGame(id int64) (*models.Game, error) {
	query := `
		SELECT id, name, genre, price
		FROM games 
		WHERE id = $1
`

	game := models.Game{}

	err := r.DB.QueryRow(query, id).Scan(
		&game.ID,
		&game.Name,
		&game.Genre,
		&game.Price,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *GameRepository) ListGames() ([]models.Game, error) {
	query := `
    SELECT id, name, genre, price
    FROM games
`
	var games []models.Game
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var game models.Game

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Genre,
			&game.Price,
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

func (r *GameRepository) RemoveGame(id int64) error {
	query := `
    DELETE FROM games
    WHERE id = $1
`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *GameRepository) GetPrice(tx *sql.Tx, gameID int64) (int, error) {
	var price int

	err := tx.QueryRow(`
		SELECT price
		FROM games
		WHERE id = $1
	`, gameID).Scan(&price)

	return price, err
}

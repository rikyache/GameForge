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

func (r *GameRepository) AddGame(name string, genre string) error {
	query := `
	INSERT INTO games (name, genre)
	VALUES ($1, $2);
	`

	_, err := r.DB.Exec(query, name, genre)

	return err
}

func (r *GameRepository) GetGame(id int) (*models.Game, error) {
	query := `
		SELECT id, name, genre 
		FROM games 
		WHERE id = $1
`

	game := models.Game{}

	err := r.DB.QueryRow(query, id).Scan(
		&game.ID,
		&game.Name,
		&game.Genre,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *GameRepository) ListGames() ([]models.Game, error) {
	query := `
    SELECT id, name, genre
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

func (r *GameRepository) RemoveGame(id int) error {
	query := `
    DELETE FROM games
    WHERE id = $1
`
	_, err := r.DB.Exec(query, id)
	return err
}

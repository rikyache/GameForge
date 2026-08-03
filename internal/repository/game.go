package repository

import (
	"database/sql"
	"testsmth/internal/models"
)

func AddGame(db *sql.DB, name string, genre string) error {
	query := `
	INSERT INTO games (name, genre)
	VALUES ($1, $2);
	`

	_, err := db.Exec(query, name, genre)

	return err
}

func GetGame(db *sql.DB, id int) (*models.Game, error) {
	query := `
		SELECT id, name, genre 
		FROM games 
		WHERE id = $1
`

	game := models.Game{}

	err := db.QueryRow(query, id).Scan(
		&game.ID,
		&game.Name,
		&game.Genre,
	)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func ListGames(db *sql.DB) ([]models.Game, error) {
	query := `
    SELECT id, name, genre
    FROM games
`
	var games []models.Game
	rows, err := db.Query(query)
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

func RemoveGame(db *sql.DB, id int) error {
	query := `
    DELETE FROM games
    WHERE id = $1
`
	_, err := db.Exec(query, id)
	return err
}

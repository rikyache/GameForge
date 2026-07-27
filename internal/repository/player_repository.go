package repository

import (
	"database/sql"
	"testsmth/internal/models"
)

func CreateTable(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS players (
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    health INT
	);
	`

	_, err := db.Exec(query)

	return err
}

func AddPlayer(db *sql.DB, name string, health int) error {
	query := `
	INSERT INTO players (name, health)
	VALUES ($1, $2)
`

	_, err := db.Exec(query, name, health)

	return err
}

func GetPlayer(db *sql.DB, id int) (models.Player, error) {
	query := `
	SELECT id, name, health FROM players
	WHERE id = $1
`

	var player models.Player

	err := db.QueryRow(query, id).Scan(
		&player.Id,
		&player.Name,
		&player.Health,
	)

	return player, err

}

func ListPlayers(db *sql.DB) ([]models.Player, error) {
	query := `
	SELECT id, name, health FROM players
	`

	var Players []models.Player

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player models.Player

		err := rows.Scan(
			&player.Id,
			&player.Name,
			&player.Health,
		)

		if err != nil {
			return nil, err
		}

		Players = append(Players, player)
	}

	return Players, nil
}

package repository

import (
	"database/sql"
	"testsmth/internal/models"
)

func AddPlayer(db *sql.DB, name string) error {
	query := `
	INSERT INTO players (name)
	VALUES ($1)
`

	_, err := db.Exec(query, name)

	return err
}

func GetPlayer(db *sql.DB, id int) (models.Player, error) {
	query := `
	SELECT id, name FROM players
	WHERE id = $1
`

	var player models.Player

	err := db.QueryRow(query, id).Scan(
		&player.Id,
		&player.Name,
	)

	return player, err

}

func ListPlayers(db *sql.DB) ([]models.Player, error) {
	query := `
	SELECT id, name FROM players
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
		)

		if err != nil {
			return nil, err
		}

		Players = append(Players, player)
	}

	return Players, nil
}

func RemovePlayer(db *sql.DB, id int) error {

	query := `
		DELETE FROM players
		WHERE id = $1
	`

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

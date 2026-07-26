package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := `
	host=postgres
	port=5432
	user=kirill
	password=12345
	dbname=practice
	sslmode=disable
	`

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

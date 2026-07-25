package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Player struct {
	id     int64  `db:"id"`
	name   string `db:"name"`
	health int    `db:"health"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Docker")
	})

	connStr := `
		host=postgres
		port=5432
		user=kirill
		password=12345
		dbname=practice
		sslmode=disable
		`

	//открыли и СРАЗУ ставим в очередь на закрытие
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
	CREATE TABLE IF NOT EXISTS players (
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    health INT
	    );
`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query15 := `
	INSERT INTO players (name, health)
	VALUES ('Kirill', 100)
`

	_, err = db.Exec(query15)

	query2 := `
	SELECT * FROM players
`

	rows, err := db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		var health int

		err := rows.Scan(&id, &name, &health)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Id: %v, name - %s, health - %v\n", id, name, health)
	}

	http.ListenAndServe(":8080", nil)
}

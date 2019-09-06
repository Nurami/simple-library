package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectToDB() {
	connStr := "user=postgres password=postgres dbname=library sslmode=disable host=localhost port=5432"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	db = database
}

func getNotes() []note {
	rows, err := db.Query("SELECT * FROM note")
	if err != nil {
		log.Panic(err)
	}
	result := make([]note, 0)
	for rows.Next() {
		tmp := note{}
		err = rows.Scan(&tmp)
		if err != nil {
			log.Panic(err)
		}
		result = append(result, tmp)
	}
	return result
}

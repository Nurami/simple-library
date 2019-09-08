package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

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
		err = rows.Scan(&tmp.ID, &tmp.BookName, &tmp.EmergenceDate, &tmp.BuyingDate, &tmp.ReadingEndDate, &tmp.Status)
		if err != nil {
			log.Panic(err)
		}
		result = append(result, tmp)
	}
	return result
}

func getNoteByID(ID string) note {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		log.Panic(err)
	}
	row := db.QueryRow("SELECT * FROM note WHERE id = $1", intID)
	tmp := note{}
	err = row.Scan(&tmp.ID, &tmp.BookName, &tmp.EmergenceDate, &tmp.BuyingDate, &tmp.ReadingEndDate, &tmp.Status)
	if err != nil {
		log.Panic(err)
	}
	return tmp
}

func addNote(note note) error {
	_, err := db.Exec("INSERT INTO note(book_name, emergence_date, buying_date, reading_end_date, status) VALUES ($1, $2, $3, $4, $5)", note.BookName, note.EmergenceDate, note.BuyingDate, note.ReadingEndDate, note.Status)
	return err
}

func changeDateAndStatusOfNote(note note) error {
	if note.Status != "want" {
		note.ReadingEndDate = time.Now()
		note.Status = "read"
		_, err := db.Exec("UPDATE note SET reading_end_date = $1, status = $2 WHERE id = $3", note.ReadingEndDate, note.Status, note.ID)
		return err
	}
	note.BuyingDate = time.Now()
	note.Status = "bought"
	_, err := db.Exec("UPDATE note SET buying_date = $1, status = $2 WHERE id = $3", note.BuyingDate, note.Status, note.ID)
	return err

}

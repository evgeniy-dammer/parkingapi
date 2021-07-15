package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var database *sql.DB

//DatabaseConnect returns database connection
func DatabaseConnect() {
	conStr := "host=127.0.0.1 port=5433 user=parking password=parking dbname=parking sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}
	database = db
}

//GetDB returns database object
func GetDB() *sql.DB {
	return database
}

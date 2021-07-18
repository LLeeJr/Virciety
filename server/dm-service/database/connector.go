package database

import (
	"database/sql"
	"log"
)

const dbsource = "postgresql://user:pass@localhost:5433/db?sslmode=disable"

func Connect() *sql.DB {
	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", dbsource)
		if err != nil {
			log.Fatal("can't connect to db", err)
		}
	}

	return db
}
package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Open() {
	var err error

	DB, err = sql.Open("sqlite", "scalar.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

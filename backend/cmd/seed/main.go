package main

import "github.com/MystiaFin/scalar/backend/internal/db"

func main() {
	db.Open()
	defer db.DB.Close()

	db.Migrate()
	db.Seed()
}

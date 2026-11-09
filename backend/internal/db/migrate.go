package db

import (
	"embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrations embed.FS

// sqlite's user_version remembers the last migration that ran
func Migrate() {
	var current int
	if err := DB.QueryRow(`PRAGMA user_version`).Scan(&current); err != nil {
		log.Fatalf("failed to read database version: %v", err)
	}

	// embed lists files sorted by name
	files, err := migrations.ReadDir("migrations")
	if err != nil {
		log.Fatalf("failed to read migrations: %v", err)
	}

	for _, file := range files {
		version := readVersion(file.Name())
		if version <= current {
			continue
		}

		runMigration(file.Name(), version)

		log.Printf("migrated: %s", file.Name())
	}
}

// 001_create_transactions.sql is version 1
func readVersion(name string) int {
	prefix, _, _ := strings.Cut(name, "_")

	version, err := strconv.Atoi(prefix)
	if err != nil {
		log.Fatalf("failed to read version from %s: %v", name, err)
	}

	return version
}

func runMigration(name string, version int) {
	query, err := migrations.ReadFile("migrations/" + name)
	if err != nil {
		log.Fatalf("failed to read %s: %v", name, err)
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Fatalf("failed to start %s: %v", name, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(query)); err != nil {
		log.Fatalf("failed to run %s: %v", name, err)
	}

	// pragma can't take a ? parameter
	setVersion := fmt.Sprintf(`PRAGMA user_version = %d`, version)
	if _, err := tx.Exec(setVersion); err != nil {
		log.Fatalf("failed to save version of %s: %v", name, err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed to commit %s: %v", name, err)
	}
}

package db

import "testing"

// every test gets its own database in a temp folder, never the real scalar.db
func openTemp(t *testing.T) {
	t.Chdir(t.TempDir())

	Open()
	t.Cleanup(func() { DB.Close() })
}

func countTransactions(t *testing.T) int {
	var count int
	if err := DB.QueryRow(`SELECT COUNT(1) FROM transactions`).Scan(&count); err != nil {
		t.Fatalf("failed to count transactions: %v", err)
	}

	return count
}

func TestMigrateCreatesEmptyTable(t *testing.T) {
	openTemp(t)

	Migrate()

	if count := countTransactions(t); count != 0 {
		t.Fatalf("expected 0 transactions, got %d", count)
	}
}

func TestMigrateTwiceChangesNothing(t *testing.T) {
	openTemp(t)

	Migrate()
	Migrate()

	var version int
	if err := DB.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		t.Fatalf("failed to read version: %v", err)
	}

	if version != 1 {
		t.Fatalf("expected version 1, got %d", version)
	}
}

func TestSeedTwiceAddsRowsOnce(t *testing.T) {
	openTemp(t)

	Migrate()
	Seed()
	Seed()

	if count := countTransactions(t); count != len(seedRows) {
		t.Fatalf("expected %d transactions, got %d", len(seedRows), count)
	}
}

func TestSeedHasTransactionToday(t *testing.T) {
	openTemp(t)

	Migrate()
	Seed()

	var count int
	query := `SELECT COUNT(1) FROM transactions WHERE date = date('now', 'localtime')`
	if err := DB.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("failed to count today's transactions: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 transaction today, got %d", count)
	}
}

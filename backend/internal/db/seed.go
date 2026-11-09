package db

import (
	"fmt"
	"log"
	"time"
)

type seedRow struct {
	merchant  string
	amount    int
	category  string
	isExpense bool
	daysAgo   int
	confirmed bool
}

// made up rows, dates follow today so the current month is never empty
var seedRows = []seedRow{
	{"Starbucks", 60000, "food & drink", true, 0, true},
	{"Grab", 25000, "transport", true, 1, true},
	{"Indomaret", 85000, "food & drink", true, 2, false},
	{"Spotify", 55000, "entertainment", true, 3, true},
	{"Tokopedia", 250000, "shopping", true, 4, false},
	{"PLN", 300000, "utilities", true, 5, true},
	{"Gojek", 30000, "transport", true, 6, false},
	{"GoPay", 100000, "transfer", false, 7, true},
	{"Kimia Farma", 75000, "health", true, 8, false},
	{"Netflix", 120000, "entertainment", true, 10, true},
	{"Shopee", 180000, "shopping", true, 12, false},
	{"Telkomsel", 100000, "utilities", true, 14, true},
	{"Alfamart", 45000, "food & drink", true, 16, false},
	{"Halodoc", 150000, "health", true, 20, true},
	{"KAI", 200000, "transport", true, 24, false},
	{"Flip", 500000, "transfer", false, 28, true},
	{"CGV", 50000, "entertainment", true, 32, false},
	{"Uniqlo", 400000, "shopping", true, 36, true},
	{"McDonald's", 70000, "food & drink", true, 40, false},
	{"Pajak Daerah", 200000, "other", true, 45, false},
}

// running this again adds nothing, the email_hash values stay the same
func Seed() {
	today := time.Now()

	for i, row := range seedRows {
		emailHash := fmt.Sprintf("seed-%02d", i+1)
		date := today.AddDate(0, 0, -row.daysAgo).Format("2006-01-02")

		_, err := DB.Exec(`
			INSERT OR IGNORE INTO transactions
				(email_hash, merchant, amount, category, is_expense, date, confirmed)
			VALUES
				(?, ?, ?, ?, ?, ?, ?)
		`, emailHash, row.merchant, row.amount, row.category, row.isExpense, date, row.confirmed)
		if err != nil {
			log.Fatalf("failed to insert %s: %v", row.merchant, err)
		}
	}

	log.Printf("seeded: %d transactions", len(seedRows))
}

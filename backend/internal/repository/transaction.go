package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/MystiaFin/scalar/backend/internal/db"
	"github.com/MystiaFin/scalar/backend/internal/models"
)

const columns = `id, email_hash, merchant, amount, category, is_expense, date, description, confirmed, updated_at`

type scanner interface {
	Scan(dest ...any) error
}

func scanTransaction(s scanner) (models.Transaction, error) {
	var tx models.Transaction
	err := s.Scan(
		&tx.ID, &tx.EmailHash, &tx.Merchant, &tx.Amount, &tx.Category,
		&tx.IsExpense, &tx.Date, &tx.Description, &tx.Confirmed, &tx.UpdatedAt,
	)
	return tx, err
}

func InsertTransaction(tx models.Transaction) error {
	_, err := db.DB.Exec(`
        INSERT OR IGNORE INTO transactions
            (email_hash, merchant, amount, category, is_expense, date, description)
        VALUES
            (?, ?, ?, ?, ?, ?, ?)
    `, tx.EmailHash, tx.Merchant, tx.Amount, tx.Category, tx.IsExpense, tx.Date, tx.Description)
	return err
}

func HasTransactions() (bool, error) {
	var count int
	err := db.DB.QueryRow(`SELECT COUNT(1) FROM transactions`).Scan(&count)
	return count > 0, err
}

func ExistsByHash(hash string) (bool, error) {
	var count int
	err := db.DB.QueryRow(
		`SELECT COUNT(1) FROM transactions WHERE email_hash = ?`, hash,
	).Scan(&count)
	return count > 0, err
}

func extractWords(s string) []string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return ' '
	}, s)
	return strings.Fields(clean)
}

func IsDuplicateTransaction(merchant string, amount float64, date string) (bool, error) {
	newMerchant := strings.ToLower(strings.TrimSpace(merchant))
	if newMerchant == "" {
		return false, nil
	}

	rows, err := db.DB.Query(
		`SELECT merchant FROM transactions WHERE amount = ? AND date = ?`, amount, date,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	newWords := extractWords(newMerchant)

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing); err != nil {
			return false, err
		}

		existingMerchant := strings.ToLower(existing)
		if existingMerchant == "" {
			continue
		}

		if strings.Contains(existingMerchant, newMerchant) || strings.Contains(newMerchant, existingMerchant) {
			return true, nil
		}

		existingWords := extractWords(existingMerchant)
		for _, w1 := range newWords {
			if len(w1) < 4 {
				continue
			}
			for _, w2 := range existingWords {
				if w1 == w2 {
					return true, nil
				}
			}
		}
	}

	return false, rows.Err()
}

func GetAll() ([]models.Transaction, error) {
	rows, err := db.DB.Query(`SELECT ` + columns + ` FROM transactions ORDER BY date DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []models.Transaction{}
	for rows.Next() {
		tx, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, tx)
	}

	return results, rows.Err()
}

func GetByID(id int64) (*models.Transaction, error) {
	row := db.DB.QueryRow(`SELECT `+columns+` FROM transactions WHERE id = ?`, id)

	tx, err := scanTransaction(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func GetConfirmedExamples(limit int) ([]models.Transaction, error) {
	rows, err := db.DB.Query(`
        SELECT merchant, amount, category, is_expense
        FROM transactions
        WHERE id IN (
            SELECT MAX(id)
            FROM transactions
            WHERE confirmed = 1
            GROUP BY merchant COLLATE NOCASE
        )
        ORDER BY updated_at DESC, id DESC
        LIMIT ?
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []models.Transaction{}
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.Merchant, &tx.Amount, &tx.Category, &tx.IsExpense); err != nil {
			return nil, err
		}
		results = append(results, tx)
	}

	return results, rows.Err()
}

func UpdateByID(id int64, category string, isExpense bool, description string) error {
	_, err := db.DB.Exec(`
        UPDATE transactions
        SET category = ?, is_expense = ?, description = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `, category, isExpense, description, id)
	return err
}

func UpdateByMerchant(merchant, category string, isExpense bool) (int64, error) {
	result, err := db.DB.Exec(`
        UPDATE transactions
        SET category = ?, is_expense = ?, confirmed = 1, updated_at = CURRENT_TIMESTAMP
        WHERE merchant = ? COLLATE NOCASE
    `, category, isExpense, merchant)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

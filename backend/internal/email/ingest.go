package email

import (
    "errors"
    "log"
    "strings"
    "sync"

    "github.com/MystiaFin/scalar/backend/internal/models"
    "github.com/MystiaFin/scalar/backend/internal/repository"
)

type Extraction struct {
    Merchant  string
    Amount    float64
    Category  string
    IsExpense bool
}

type Extractor func(subject, body string) (Extraction, error)

type Summary struct {
    Fetched int
    Saved   int
    Skipped int
    Failed  int
}

var ErrSyncBusy = errors.New("sync is still running")

var (
    syncMu        sync.Mutex
    skippedMu     sync.Mutex
    skippedHashes = map[string]bool{}
)

func isSkipped(hash string) bool {
    skippedMu.Lock()
    defer skippedMu.Unlock()

    return skippedHashes[hash]
}

func markSkipped(hash string) {
    skippedMu.Lock()
    defer skippedMu.Unlock()

    skippedHashes[hash] = true
}

func Sync(days int, extract Extractor) (Summary, error) {
    if extract == nil {
        return Summary{}, errors.New("extractor belum diatur")
    }

    if !syncMu.TryLock() {
        return Summary{}, ErrSyncBusy
    }
    defer syncMu.Unlock()

    emails, err := FetchReceipts(days, func(hash string) bool {
        if isSkipped(hash) {
            return true
        }

        exists, err := repository.ExistsByHash(hash)
        if err != nil {
            log.Printf("failed to check hash: %v", err)
            return true
        }

        return exists
    })
    if err != nil {
        return Summary{}, err
    }

    summary := Summary{Fetched: len(emails)}

    for _, raw := range emails {
        result, err := extract(raw.Subject, raw.Body)
        if err != nil {
            log.Printf("extract failed (%s): %v", raw.Subject, err)
            summary.Failed++
            continue
        }

        merchant := strings.TrimSpace(result.Merchant)
        if merchant == "" || result.Amount <= 0 {
            log.Printf("incomplete extraction, skipped: %s", raw.Subject)
            summary.Failed++
            continue
        }

        category := result.Category
        if !models.IsValidCategory(category) {
            category = "other"
        }

        date := raw.Date.Local().Format("2006-01-02")

        duplicate, err := repository.IsDuplicateTransaction(merchant, result.Amount, date)
        if err != nil {
            log.Printf("failed to check duplicate: %v", err)
            summary.Failed++
            continue
        }
        if duplicate {
            markSkipped(raw.Hash)
            summary.Skipped++
            continue
        }

        tx := models.Transaction{
            EmailHash: raw.Hash,
            Merchant:  merchant,
            Amount:    result.Amount,
            Category:  category,
            IsExpense: result.IsExpense,
            Date:      date,
        }

        if err := repository.InsertTransaction(tx); err != nil {
            log.Printf("failed to insert: %v", err)
            summary.Failed++
            continue
        }

        summary.Saved++
    }

    return summary, nil
}

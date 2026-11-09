package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MystiaFin/scalar/backend/internal/db"
	"github.com/MystiaFin/scalar/backend/internal/models"
	"github.com/MystiaFin/scalar/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func newRouter(t *testing.T) *gin.Engine {
	t.Chdir(t.TempDir())

	db.Open()
	t.Cleanup(func() { db.DB.Close() })

	db.Migrate()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	Register(r)

	return r
}

func addTransaction(t *testing.T, hash string, merchant string) {
	tx := models.Transaction{
		EmailHash: hash,
		Merchant:  merchant,
		Amount:    50000,
		Category:  "shopping",
		IsExpense: true,
		Date:      "2026-10-01",
	}

	if err := repository.InsertTransaction(tx); err != nil {
		t.Fatalf("failed to insert transaction: %v", err)
	}
}

func send(r *gin.Engine, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	return rec
}

func readTransactions(t *testing.T, r *gin.Engine) []models.Transaction {
	rec := send(r, http.MethodGet, "/api/transactions", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(rec.Body.Bytes(), &transactions); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	return transactions
}

func TestEmptyDatabaseReturnsEmptyList(t *testing.T) {
	r := newRouter(t)

	rec := send(r, http.MethodGet, "/api/transactions", "")

	if rec.Body.String() != "[]" {
		t.Fatalf("expected [], got %s", rec.Body)
	}
}

func TestListReturnsSavedTransactions(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")

	transactions := readTransactions(t, r)

	if len(transactions) != 1 || transactions[0].Merchant != "Tokopedia" {
		t.Fatalf("unexpected transactions: %+v", transactions)
	}
}

func TestGetOneTransaction(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")

	id := readTransactions(t, r)[0].ID

	rec := send(r, http.MethodGet, fmt.Sprintf("/api/transactions/%d", id), "")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
	}

	var tx models.Transaction
	if err := json.Unmarshal(rec.Body.Bytes(), &tx); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if tx.Merchant != "Tokopedia" {
		t.Fatalf("unexpected transaction: %+v", tx)
	}

	rec = send(r, http.MethodGet, "/api/transactions/999", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body)
	}
}

func TestPatchUpdatesOnlyThatTransaction(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")
	addTransaction(t, "hash-2", "Tokopedia")

	transactions := readTransactions(t, r)
	path := fmt.Sprintf("/api/transactions/%d", transactions[0].ID)

	rec := send(r, http.MethodPatch, path, `{"category":"other","is_expense":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
	}

	changed := 0
	for _, tx := range readTransactions(t, r) {
		if tx.Category == "other" {
			changed++
		}
	}

	if changed != 1 {
		t.Fatalf("expected 1 changed transaction, got %d", changed)
	}
}

func TestPatchKeepsOmittedFields(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")

	id := readTransactions(t, r)[0].ID
	path := fmt.Sprintf("/api/transactions/%d", id)

	rec := send(r, http.MethodPatch, path, `{"description":"beli headset"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
	}

	tx := readTransactions(t, r)[0]

	if tx.Description != "beli headset" || tx.Category != "shopping" || !tx.IsExpense {
		t.Fatalf("omitted fields were overwritten: %+v", tx)
	}
}

func TestPatchConfirmedUpdatesWholeMerchant(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")
	addTransaction(t, "hash-2", "Tokopedia")

	transactions := readTransactions(t, r)
	path := fmt.Sprintf("/api/transactions/%d", transactions[0].ID)

	rec := send(r, http.MethodPatch, path, `{"category":"other","is_expense":true,"description":"x","confirmed":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body)
	}

	described := 0
	for _, tx := range readTransactions(t, r) {
		if tx.Category != "other" || !tx.Confirmed {
			t.Fatalf("expected every entry updated and confirmed: %+v", tx)
		}
		if tx.Description == "x" {
			described++
		}
	}

	if described != 1 {
		t.Fatalf("expected description on 1 entry, got %d", described)
	}
}

func TestPatchRejectsUnknownCategory(t *testing.T) {
	r := newRouter(t)

	addTransaction(t, "hash-1", "Tokopedia")

	id := readTransactions(t, r)[0].ID
	path := fmt.Sprintf("/api/transactions/%d", id)

	rec := send(r, http.MethodPatch, path, `{"category":"xyz"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body)
	}
}

func TestPatchUnknownTransactionReturnsNotFound(t *testing.T) {
	r := newRouter(t)

	rec := send(r, http.MethodPatch, "/api/transactions/999", `{"category":"other"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body)
	}
}

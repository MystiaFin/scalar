package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MystiaFin/scalar/backend/internal/models"
	"github.com/MystiaFin/scalar/backend/internal/repository"
)

type TransactionResult struct {
	Merchant  string  `json:"merchant"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	IsExpense bool    `json:"is_expense"`
}

const (
	Model             = "llama3:8b"
	ollamaGenerateURL = "http://localhost:11434/api/generate"
	maxBodyChars      = 3000
)

var ollamaClient = &http.Client{Timeout: 3 * time.Minute}

type ollamaRequest struct {
	Model  string          `json:"model"`
	Prompt string          `json:"prompt"`
	Format json.RawMessage `json:"format"`
	Stream bool            `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

var extractSchema = json.RawMessage(`{
    "type": "object",
    "properties": {
        "merchant":   { "type": "string" },
        "amount":     { "type": "number" },
        "category": {
            "type": "string",
            "enum": ["food & drink", "transport", "utilities", "entertainment", "shopping", "health", "transfer", "other"]
        },
        "is_expense": { "type": "boolean" }
    },
    "required": ["merchant", "amount", "category", "is_expense"]
}`)

func truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}

	return string(runes[:limit])
}

func buildFewShotBlock(examples []models.Transaction) string {
	if len(examples) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Here are past transactions that were confirmed correct — use these as reference:\n")
	for _, ex := range examples {
		sb.WriteString(fmt.Sprintf(
			"  merchant: %s → category: %s, is_expense: %v\n",
			ex.Merchant, ex.Category, ex.IsExpense,
		))
	}
	sb.WriteString("\n")
	return sb.String()
}

func AskOllama(subject, body string) (TransactionResult, error) {
	examples, err := repository.GetConfirmedExamples(10)
	if err != nil {
		log.Printf("could not load confirmed examples: %v", err)
		examples = nil
	}

	fewShot := buildFewShotBlock(examples)

	prompt := fmt.Sprintf(`You are parsing a financial transaction receipt or email.
Extract the merchant name, amount, category, and whether it is an expense.

Category guide:
- food & drink  : restaurants, cafes, groceries, snacks, coffee shops, delivery
- health        : pharmacy, clinic, hospital, doctor, supplements, medical
- transport     : ride-sharing (Uber, Lyft), public transit, toll, parking, fuel, flights
- utilities     : electricity, water, internet, mobile phone plans, monthly subscriptions
- entertainment : cinema, streaming (Netflix, Spotify), gaming, events, hobbies
- shopping      : clothing, electronics, books, household items, retail, marketplaces (Amazon)
- transfer      : bank transfers, digital wallet top-ups (PayPal, Venmo, Apple Pay)
- other         : anything that clearly does not fit the above

%sSubject: %s
Body: %s

Rules:
- merchant: extract ONLY the core primary brand name. Strip out payment gateways (e.g., Stripe, PayPal, Square), locations/branches, corporate suffixes (e.g., LLC, Inc, Ltd), and ignore receipt/order IDs.
- amount: total transaction value as a plain number. Strip ALL currency symbols ($, €, £, Rp, USD, etc.). Use 0 if not found.
- category: pick the single best fit from the list above.
- is_expense: true if money is going out, false if money is coming in.

Respond only with the JSON object.`, fewShot, subject, truncate(body, maxBodyChars))

	reqBody, err := json.Marshal(ollamaRequest{
		Model:  Model,
		Prompt: prompt,
		Format: extractSchema,
		Stream: false,
	})
	if err != nil {
		return TransactionResult{}, fmt.Errorf("marshal: %w", err)
	}

	resp, err := ollamaClient.Post(ollamaGenerateURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return TransactionResult{}, fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return TransactionResult{}, fmt.Errorf("ollama returned %s: %s", resp.Status, strings.TrimSpace(string(detail)))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return TransactionResult{}, fmt.Errorf("decode response: %w", err)
	}

	var result TransactionResult
	if err := json.Unmarshal([]byte(ollamaResp.Response), &result); err != nil {
		return TransactionResult{}, fmt.Errorf("parse json: %w", err)
	}

	return result, nil
}

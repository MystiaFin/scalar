package email

import (
    "encoding/json"
    "fmt"
    "strings"

    "github.com/MystiaFin/scalar/backend/internal/models"
)

const maxBodyChars = 3000

const promptTemplate = `You are parsing a financial transaction receipt or email.
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

Respond only with the JSON object.`

// the enum comes from models.Categories so the model can only answer with a known category
var extractSchema = buildSchema(models.Categories)

func buildSchema(categories []string) json.RawMessage {
    schema := map[string]any{
        "type": "object",
        "properties": map[string]any{
            "merchant": map[string]any{"type": "string"},
            "amount":   map[string]any{"type": "number"},
            "category": map[string]any{
                "type": "string",
                "enum": categories,
            },
            "is_expense": map[string]any{"type": "boolean"},
        },
        "required": []string{"merchant", "amount", "category", "is_expense"},
    }

    data, err := json.Marshal(schema)
    if err != nil {
        panic(fmt.Sprintf("failed to build extract schema: %v", err))
    }

    return data
}

func buildPrompt(subject, body string, examples []models.Transaction) string {
    fewShot := buildFewShotBlock(examples)

    return fmt.Sprintf(promptTemplate, fewShot, subject, truncate(body, maxBodyChars))
}

func buildFewShotBlock(examples []models.Transaction) string {
    if len(examples) == 0 {
        return ""
    }

    var sb strings.Builder

    sb.WriteString("Here are past transactions that were confirmed correct — use these as reference:\n")

    for _, ex := range examples {
        fmt.Fprintf(&sb, "  merchant: %s → category: %s, is_expense: %v\n", ex.Merchant, ex.Category, ex.IsExpense)
    }

    sb.WriteString("\n")

    return sb.String()
}

func truncate(s string, limit int) string {
    runes := []rune(s)
    if len(runes) <= limit {
        return s
    }

    return string(runes[:limit])
}

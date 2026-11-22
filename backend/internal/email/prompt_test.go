package email

import (
    "encoding/json"
    "slices"
    "strings"
    "testing"

    "github.com/MystiaFin/scalar/backend/internal/models"
)

func TestSchemaEnumMatchesCategories(t *testing.T) {
    var schema struct {
        Properties struct {
            Category struct {
                Enum []string `json:"enum"`
            } `json:"category"`
        } `json:"properties"`
        Required []string `json:"required"`
    }

    if err := json.Unmarshal(extractSchema, &schema); err != nil {
        t.Fatalf("schema is not valid json: %v", err)
    }

    if !slices.Equal(schema.Properties.Category.Enum, models.Categories) {
        t.Errorf("enum = %v, want %v", schema.Properties.Category.Enum, models.Categories)
    }

    if len(schema.Required) != 4 {
        t.Errorf("required = %v", schema.Required)
    }
}

func TestPromptIncludesSubjectAndBody(t *testing.T) {
    prompt := buildPrompt("Your Grab receipt", "Total Rp 25.000", nil)

    if !strings.Contains(prompt, "Subject: Your Grab receipt") {
        t.Error("prompt is missing the subject")
    }

    if !strings.Contains(prompt, "Body: Total Rp 25.000") {
        t.Error("prompt is missing the body")
    }

    if strings.Contains(prompt, "confirmed correct") {
        t.Error("prompt has a few-shot block without examples")
    }
}

func TestPromptIncludesConfirmedExamples(t *testing.T) {
    examples := []models.Transaction{
        {Merchant: "Tokopedia", Category: "shopping", IsExpense: true},
    }

    prompt := buildPrompt("subject", "body", examples)

    if !strings.Contains(prompt, "merchant: Tokopedia → category: shopping, is_expense: true") {
        t.Errorf("prompt is missing the example:\n%s", prompt)
    }
}

func TestPromptTruncatesLongBody(t *testing.T) {
    body := strings.Repeat("é", maxBodyChars+500)

    prompt := buildPrompt("subject", body, nil)

    if strings.Count(prompt, "é") != maxBodyChars {
        t.Errorf("body was not truncated to %d characters", maxBodyChars)
    }
}

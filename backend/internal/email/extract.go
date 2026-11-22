package email

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/MystiaFin/scalar/backend/internal/ollama"
    "github.com/MystiaFin/scalar/backend/internal/repository"
)

const (
    defaultModel    = "llama3:8b"
    fewShotLimit    = 10
    generateTimeout = 3 * time.Minute
)

type extractResponse struct {
    Merchant  string  `json:"merchant"`
    Amount    float64 `json:"amount"`
    Category  string  `json:"category"`
    IsExpense bool    `json:"is_expense"`
}

func Extract(subject, body string) (Extraction, error) {
    examples, err := repository.GetConfirmedExamples(fewShotLimit)
    if err != nil {
        log.Printf("could not load confirmed examples: %v", err)
        examples = nil
    }

    prompt := buildPrompt(subject, body, examples)

    ctx, cancel := context.WithTimeout(context.Background(), generateTimeout)
    defer cancel()

    return extractWith(ctx, ollama.FromEnv(), model(), prompt)
}

func extractWith(ctx context.Context, client *ollama.Client, model, prompt string) (Extraction, error) {
    raw, err := client.Generate(ctx, ollama.GenerateRequest{
        Model:  model,
        Prompt: prompt,
        Format: extractSchema,
    })
    if err != nil {
        return Extraction{}, err
    }

    var resp extractResponse
    if err := json.Unmarshal([]byte(raw), &resp); err != nil {
        return Extraction{}, fmt.Errorf("parse extraction: %w", err)
    }

    return Extraction(resp), nil
}

func model() string {
    if name := os.Getenv("OLLAMA_MODEL"); name != "" {
        return name
    }

    return defaultModel
}

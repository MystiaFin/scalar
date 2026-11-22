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

func AskOllama(subject, body string) (TransactionResult, error) {
    examples, err := repository.GetConfirmedExamples(10)
    if err != nil {
        log.Printf("could not load confirmed examples: %v", err)
        examples = nil
    }

    prompt := buildPrompt(subject, body, examples)

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

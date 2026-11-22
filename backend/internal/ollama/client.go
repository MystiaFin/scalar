package ollama

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
)

const DefaultURL = "http://localhost:11434"

var httpClient = &http.Client{}

type Client struct {
    baseURL string
}

type StatusError struct {
    Status string
    Detail string
}

func (e *StatusError) Error() string {
    return fmt.Sprintf("ollama returned %s: %s", e.Status, e.Detail)
}

type GenerateRequest struct {
    Model  string          `json:"model"`
    Prompt string          `json:"prompt"`
    Format json.RawMessage `json:"format,omitempty"`
    Stream bool            `json:"stream"`
}

type generateResponse struct {
    Response string `json:"response"`
}

func New(baseURL string) *Client {
    return &Client{baseURL: strings.TrimRight(baseURL, "/")}
}

// reads OLLAMA_URL at call time so values from .env are picked up
func FromEnv() *Client {
    if url := os.Getenv("OLLAMA_URL"); url != "" {
        return New(url)
    }

    return New(DefaultURL)
}

func (c *Client) Generate(ctx context.Context, req GenerateRequest) (string, error) {
    payload, err := json.Marshal(req)
    if err != nil {
        return "", fmt.Errorf("marshal request: %w", err)
    }

    body, err := c.send(ctx, http.MethodPost, "/api/generate", bytes.NewReader(payload))
    if err != nil {
        return "", err
    }

    var resp generateResponse
    if err := json.Unmarshal(body, &resp); err != nil {
        return "", fmt.Errorf("decode response: %w", err)
    }

    return resp.Response, nil
}

func (c *Client) Tags(ctx context.Context) (json.RawMessage, error) {
    body, err := c.send(ctx, http.MethodGet, "/api/tags", nil)
    if err != nil {
        return nil, err
    }

    if !json.Valid(body) {
        return nil, fmt.Errorf("decode response: invalid json")
    }

    return json.RawMessage(body), nil
}

func (c *Client) send(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
    if err != nil {
        return nil, fmt.Errorf("build request: %w", err)
    }

    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("ollama request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        detail, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
        return nil, &StatusError{Status: resp.Status, Detail: strings.TrimSpace(string(detail))}
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    return data, nil
}

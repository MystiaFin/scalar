package ollama

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGenerateSendsRequestAndReturnsResponse(t *testing.T) {
    var got GenerateRequest

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost || r.URL.Path != "/api/generate" {
            t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
        }

        if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
            t.Errorf("failed to decode request: %v", err)
        }

        w.Write([]byte(`{"response": "{\"ok\": true}"}`))
    }))
    defer server.Close()

    out, err := New(server.URL + "/").Generate(context.Background(), GenerateRequest{
        Model:  "test-model",
        Prompt: "hello",
        Format: json.RawMessage(`{"type":"object"}`),
    })
    if err != nil {
        t.Fatalf("generate failed: %v", err)
    }

    if out != `{"ok": true}` {
        t.Errorf("response = %q", out)
    }

    if got.Model != "test-model" || got.Prompt != "hello" || got.Stream {
        t.Errorf("request = %+v", got)
    }
}

func TestNonOKStatusReturnsStatusError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "model not found", http.StatusNotFound)
    }))
    defer server.Close()

    _, err := New(server.URL).Generate(context.Background(), GenerateRequest{Model: "missing"})

    var statusErr *StatusError
    if !errors.As(err, &statusErr) {
        t.Fatalf("expected StatusError, got %v", err)
    }

    if statusErr.Detail != "model not found" {
        t.Errorf("detail = %q", statusErr.Detail)
    }
}

func TestTagsReturnsRawJSON(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/api/tags" {
            t.Errorf("unexpected path: %s", r.URL.Path)
        }

        w.Write([]byte(`{"models":[{"name":"llama3:8b"}]}`))
    }))
    defer server.Close()

    tags, err := New(server.URL).Tags(context.Background())
    if err != nil {
        t.Fatalf("tags failed: %v", err)
    }

    if string(tags) != `{"models":[{"name":"llama3:8b"}]}` {
        t.Errorf("tags = %s", tags)
    }
}

func TestTagsRejectsInvalidJSON(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`not json`))
    }))
    defer server.Close()

    if _, err := New(server.URL).Tags(context.Background()); err == nil {
        t.Fatal("expected error for invalid json")
    }
}

func TestFromEnvUsesOllamaURL(t *testing.T) {
    t.Setenv("OLLAMA_URL", "http://example.test:1234/")

    if got := FromEnv().baseURL; got != "http://example.test:1234" {
        t.Errorf("baseURL = %q", got)
    }

    t.Setenv("OLLAMA_URL", "")

    if got := FromEnv().baseURL; got != DefaultURL {
        t.Errorf("baseURL = %q", got)
    }
}

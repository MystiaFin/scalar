package routes

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestModelsPassesThroughOllamaTags(t *testing.T) {
    ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"models":[{"name":"llama3:8b"}]}`))
    }))
    defer ollama.Close()

    t.Setenv("OLLAMA_URL", ollama.URL)

    r := newRouter(t)

    rec := send(r, http.MethodGet, "/api/models", "")
    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
    }

    if rec.Body.String() != `{"models":[{"name":"llama3:8b"}]}` {
        t.Errorf("body = %s", rec.Body.String())
    }
}

func TestModelsReturnsBadGatewayWhenOllamaIsDown(t *testing.T) {
    ollama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
    ollama.Close()

    t.Setenv("OLLAMA_URL", ollama.URL)

    r := newRouter(t)

    rec := send(r, http.MethodGet, "/api/models", "")
    if rec.Code != http.StatusBadGateway {
        t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
    }
}

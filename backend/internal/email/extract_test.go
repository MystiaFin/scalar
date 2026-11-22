package email

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/MystiaFin/scalar/backend/internal/ollama"
)

func fakeOllama(t *testing.T, response string, gotModel *string) *ollama.Client {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var req ollama.GenerateRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            t.Errorf("failed to decode request: %v", err)
        }

        if gotModel != nil {
            *gotModel = req.Model
        }

        if len(req.Format) == 0 {
            t.Error("request has no schema")
        }

        out, _ := json.Marshal(map[string]string{"response": response})
        w.Write(out)
    }))
    t.Cleanup(server.Close)

    return ollama.New(server.URL)
}

func TestExtractParsesModelAnswer(t *testing.T) {
    var gotModel string

    client := fakeOllama(t, `{"merchant":"Grab","amount":25000,"category":"transport","is_expense":true}`, &gotModel)

    got, err := extractWith(context.Background(), client, "test-model", "prompt")
    if err != nil {
        t.Fatalf("extract failed: %v", err)
    }

    want := Extraction{Merchant: "Grab", Amount: 25000, Category: "transport", IsExpense: true}
    if got != want {
        t.Errorf("extraction = %+v, want %+v", got, want)
    }

    if gotModel != "test-model" {
        t.Errorf("model = %q", gotModel)
    }
}

func TestExtractRejectsNonJSONAnswer(t *testing.T) {
    client := fakeOllama(t, "sorry, I cannot help with that", nil)

    if _, err := extractWith(context.Background(), client, "test-model", "prompt"); err == nil {
        t.Fatal("expected error for non-json answer")
    }
}

func TestModelUsesOllamaModelEnv(t *testing.T) {
    t.Setenv("OLLAMA_MODEL", "qwen2.5:7b")

    if got := model(); got != "qwen2.5:7b" {
        t.Errorf("model = %q", got)
    }

    t.Setenv("OLLAMA_MODEL", "")

    if got := model(); got != defaultModel {
        t.Errorf("model = %q", got)
    }
}

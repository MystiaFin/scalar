package handlers

import (
    "context"
    "errors"
    "net/http"
    "time"

    "github.com/MystiaFin/scalar/backend/internal/ollama"

    "github.com/gin-gonic/gin"
)

const tagsTimeout = 5 * time.Second

func GetAvailableModels(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), tagsTimeout)
    defer cancel()

    tags, err := ollama.FromEnv().Tags(ctx)

    var statusErr *ollama.StatusError
    if errors.As(err, &statusErr) {
        c.JSON(http.StatusBadGateway, gin.H{"error": "Ollama returned " + statusErr.Status})
        return
    }

    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{
            "error":  "could not get models from Ollama, is it running?",
            "detail": err.Error(),
        })
        return
    }

    c.Data(http.StatusOK, "application/json; charset=utf-8", tags)
}

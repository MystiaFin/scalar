package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

const ollamaTagsURL = "http://localhost:11434/api/tags"

var ollamaClient = &http.Client{Timeout: 5 * time.Second}

func GetAvailableModels(c *gin.Context) {
    resp, err := ollamaClient.Get(ollamaTagsURL)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{
            "error":  "could not connect to local Ollama, is it running?",
            "detail": err.Error(),
        })
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.JSON(http.StatusBadGateway, gin.H{"error": "Ollama returned " + resp.Status})
        return
    }

    var result map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "failed to parse Ollama response"})
        return
    }

    c.JSON(http.StatusOK, result)
}

package routes

import (
    "github.com/MystiaFin/scalar/backend/internal/handlers"

    "github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
    api := r.Group("/api")
    {
        api.GET("/transactions", handlers.GetTransactions)
        api.GET("/transactions/:id", handlers.GetTransaction)
        api.PATCH("/transactions/:id", handlers.UpdateTransaction)
        api.GET("/models", handlers.GetAvailableModels)
    }
}

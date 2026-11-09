package handlers

import (
	"net/http"
	"strconv"

	"github.com/MystiaFin/scalar/backend/internal/models"
	"github.com/MystiaFin/scalar/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type updateRequest struct {
	Category    *string `json:"category"`
	IsExpense   *bool   `json:"is_expense"`
	Description *string `json:"description"`
	Confirmed   *bool   `json:"confirmed"`
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}

	return id, true
}

func GetTransactions(c *gin.Context) {
	transactions, err := repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func GetTransaction(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	tx, err := repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func UpdateTransaction(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var body updateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if body.Category != nil && !models.IsValidCategory(*body.Category) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "invalid category",
			"categories": models.Categories,
		})
		return
	}

	tx, err := repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	category := tx.Category
	if body.Category != nil {
		category = *body.Category
	}

	isExpense := tx.IsExpense
	if body.IsExpense != nil {
		isExpense = *body.IsExpense
	}

	description := tx.Description
	if body.Description != nil {
		description = *body.Description
	}

	if err := repository.UpdateByID(id, category, isExpense, description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	affected := int64(1)

	if body.Confirmed != nil && *body.Confirmed {
		affected, err = repository.UpdateByMerchant(tx.Merchant, category, isExpense)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "updated",
		"affected": affected,
	})
}

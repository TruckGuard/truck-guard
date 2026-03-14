package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/repository"
)

func HandleListSystemSettings(c *gin.Context) {
	settings, err := repository.ListSystemSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func HandleUpsertSystemSetting(c *gin.Context) {
	var input struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.UpsertSystemSetting(c.Request.Context(), input.Key, input.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

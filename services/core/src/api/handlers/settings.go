package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

// System Settings Handlers

func HandleListSettings(c *gin.Context) {
	var settings []models.SystemSetting
	if err := repository.DB.WithContext(c.Request.Context()).Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func HandleUpdateSetting(c *gin.Context) {
	var input models.SystemSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.SystemSetting
	if err := repository.DB.WithContext(c.Request.Context()).Where("key = ?", input.Key).First(&setting).Error; err != nil {
		// Create if not exists
		setting = input
		if err := repository.DB.WithContext(c.Request.Context()).Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create setting"})
			return
		}
	} else {
		setting.Value = input.Value
		if err := repository.DB.WithContext(c.Request.Context()).Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
			return
		}
	}

	c.JSON(http.StatusOK, setting)
}


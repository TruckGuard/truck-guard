package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
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

	authID := c.GetHeader("X-User-ID")
	var setting models.SystemSetting
	var action string
	var oldValue string

	if err := repository.DB.WithContext(c.Request.Context()).Where("key = ?", input.Key).First(&setting).Error; err != nil {
		// Create if not exists
		setting = input
		action = "create"
		if err := repository.DB.WithContext(c.Request.Context()).Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create setting"})
			return
		}
	} else {
		oldValue = setting.Value
		setting.Value = input.Value
		action = "update"
		if err := repository.DB.WithContext(c.Request.Context()).Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
			return
		}
	}

	// Log audit
	repository.LogSystemAudit(c.Request.Context(), authID, action, "setting:"+setting.Key, map[string]string{
		"key":  setting.Key,
		"old":  oldValue,
		"new":  setting.Value,
		"name": setting.Name,
	}, "Зміна системних налаштувань")

	c.JSON(http.StatusOK, setting)
}

func HandleListSystemAuditEvents(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	var audits []models.SystemAudit
	var total int64

	q := repository.DB.WithContext(c.Request.Context()).Model(&models.SystemAudit{})

	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count audits"})
		return
	}

	err := q.
		Preload("User").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audits"})
		return
	}

	utils.SendPaginatedResponse(c, audits, total, page, limit)
}

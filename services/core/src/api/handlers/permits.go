package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
)

func HandleGetPermits(c *gin.Context) {
	slog.Info("HandleGetPermits: Start")
	limit, offset, page := utils.GetPagination(c)
	plate := c.Query("plate")

	permits, total, err := repository.GetPermits(c.Request.Context(), limit, offset, plate)
	if err != nil {
		slog.Error("HandleGetPermits: error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permits"})
		return
	}
	slog.Info("HandleGetPermits: done", "count", len(permits), "total", total)

	utils.SendPaginatedResponse(c, permits, total, page, limit)
}

func HandleCreatePermit(c *gin.Context) {
	var input models.Permit
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID := c.GetHeader("X-User-ID")
	if err := repository.CreatePermit(c.Request.Context(), &input, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permit"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func HandleUpdatePermit(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID := c.GetHeader("X-User-ID")
	permit, err := repository.UpdatePermit(c.Request.Context(), id, input, authID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permit"})
		return
	}

	c.JSON(http.StatusOK, permit)
}

func HandleGetPermitByID(c *gin.Context) {
	id := c.Param("id")
	permit, err := repository.GetPermitByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permit not found"})
		return
	}

	c.JSON(http.StatusOK, permit)
}

func HandleValidatePermit(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetHeader("X-User-ID")
	permit, err := repository.ValidatePermit(c.Request.Context(), id, authID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate permit"})
		return
	}
	c.JSON(http.StatusOK, permit)
}

func HandleGetPermitAuditEvents(c *gin.Context) {
	id := c.Param("id")
	var audits []models.PermitAudit

	// Ensure permit exists? Or just query audits directly.
	if err := repository.DB.WithContext(c.Request.Context()).Preload("User").Where("permit_id = ?", id).Order("created_at desc").Find(&audits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit events"})
		return
	}

	c.JSON(http.StatusOK, audits)
}

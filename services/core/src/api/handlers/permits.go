package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/pkg/notify"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
)

func HandleGetPermits(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)

	// Parse is_closed filter
	slog.Debug("Full url", "url", c.Request.URL.String())
	var isClosed *bool
	if val := c.Query("is_closed"); val != "" {
		b := val == "true"
		isClosed = &b
	}

	// Sanitise sort_order
	sortOrder := strings.ToLower(c.Query("sort_order"))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	params := repository.PermitQueryParams{
		Plate:             c.Query("plate"),
		IsClosed:          isClosed,
		SortField:         c.Query("sort_field"),
		SortOrder:         sortOrder,
		FilterFrom:        c.Query("filter_from"),
		FilterTo:          c.Query("filter_to"),
		FilterPostID:      c.Query("filter_post_id"),
		FilterVehicleType: c.Query("filter_vehicle_type"),
		FilterPaymentType: c.Query("filter_payment_type"),
		FilterPayer:       c.Query("filter_payer"),
		Search:            c.Query("search"),
		CustomFilters:     c.Query("custom_filters"),
	}

	permits, total, err := repository.GetPermits(c.Request.Context(), limit, offset, params)
	if err != nil {
		slog.Error("HandleGetPermits: error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permits"})
		return
	}

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

	if input.CustomsPostID != nil {
		plate := input.PlateFront
		if plate == "" {
			plate = input.PlateBack
		}
		notify.Global.PublishPermit(*input.CustomsPostID, input.ID, input.Code, plate)
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
	audits, err := repository.GetPermitAuditEvents(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit events"})
		return
	}
	c.JSON(http.StatusOK, audits)
}

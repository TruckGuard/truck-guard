package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/pkg/notify"
	"github.com/truckguard/core/src/repository"
	"github.com/truckguard/core/src/utils"
	"gorm.io/gorm"
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

	var isVoid *bool
	if val := c.Query("is_void"); val != "" {
		b := val == "true"
		isVoid = &b
		slog.Info("HandleGetPermits: void filter applied", "val", val, "b", b)
	} else {
		slog.Info("HandleGetPermits: void filter NOT applied (showing all)")
	}

	// Sanitise sort_order
	sortOrder := strings.ToLower(c.Query("sort_order"))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	params := repository.PermitQueryParams{
		Plate:             c.Query("plate"),
		IsClosed:          isClosed,
		IsVoid:            isVoid,
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
	var input struct {
		models.Permit
		CameraEventID *uint `json:"camera_event_id,omitempty"`
		ScaleEventID  *uint `json:"scale_event_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID := c.GetHeader("X-User-ID")
	if err := repository.CreatePermit(c.Request.Context(), &input.Permit, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permit"})
		return
	}

	// Link events if provided
	if input.CameraEventID != nil {
		var event models.PlateEvent
		if err := repository.DB.First(&event, *input.CameraEventID).Error; err == nil {
			event.PermitID = &input.Permit.ID
			repository.DB.Save(&event)
			// Also fill plate from event if permit is empty
			if input.Permit.PlateFront == "" {
				repository.DB.Model(&input.Permit).Update("plate_front", event.Plate)
				input.Permit.PlateFront = event.Plate
			}
		}
	}
	if input.ScaleEventID != nil {
		var event models.WeightEvent
		if err := repository.DB.First(&event, *input.ScaleEventID).Error; err == nil {
			event.PermitID = &input.Permit.ID
			repository.DB.Save(&event)
			// Also fill weight if permit is empty
			if input.Permit.TotalWeight == 0 {
				repository.DB.Model(&input.Permit).Update("total_weight", event.Weight)
				input.Permit.TotalWeight = event.Weight
			}
		}
	}

	if input.Permit.CustomsPostID != nil {
		plate := input.Permit.PlateFront
		if plate == "" {
			plate = input.Permit.PlateBack
		}
		notify.Global.PublishPermit(*input.Permit.CustomsPostID, input.Permit.ID, input.Permit.Code, plate)
	}

	c.JSON(http.StatusCreated, input.Permit)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Перепустку не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при завантаженні перепустки"})
		}
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

func HandleListAuditEvents(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)

	params := repository.AuditQueryParams{
		Action:     c.Query("action"),
		UserID:     c.Query("user_id"),
		FilterFrom: c.Query("filter_from"),
		FilterTo:   c.Query("filter_to"),
	}

	audits, total, err := repository.GetAllAuditEvents(c.Request.Context(), limit, offset, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit events"})
		return
	}

	utils.SendPaginatedResponse(c, audits, total, page, limit)
}

func HandleRestorePermit(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetHeader("X-User-ID")
	if err := repository.RestorePermit(c.Request.Context(), id, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore permit"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "restored"})
}

func HandleVoidPermit(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetHeader("X-User-ID")
	if err := repository.VoidPermit(c.Request.Context(), id, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to void permit"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "voided"})
}

func HandleDeletePermit(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetHeader("X-User-ID")
	if err := repository.DeletePermit(c.Request.Context(), id, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete permit"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func HandleLinkPermit(c *gin.Context) {
	var input struct {
		PermitID  uint   `json:"permit_id" binding:"required"`
		EventID   uint   `json:"event_id" binding:"required"`
		EventType string `json:"event_type" binding:"required"` // plate or weight
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID := c.GetHeader("X-User-ID")
	if err := repository.LinkPermitEvent(c.Request.Context(), input.PermitID, input.EventID, input.EventType, authID); err != nil {
		slog.Error("HandleLinkPermit: error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link permit: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "linked"})
}

func HandleUnlinkPermit(c *gin.Context) {
	var input struct {
		PermitID  uint   `json:"permit_id" binding:"required"`
		EventID   uint   `json:"event_id" binding:"required"`
		EventType string `json:"event_type" binding:"required"` // plate or weight
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID := c.GetHeader("X-User-ID")
	if err := repository.UnlinkPermitEvent(c.Request.Context(), input.PermitID, input.EventID, input.EventType, authID); err != nil {
		slog.Error("HandleUnlinkPermit: error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlink permit: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unlinked"})
}

func HandleLogPrintPermit(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetHeader("X-User-ID")
	if err := repository.LogPrintPermit(c.Request.Context(), id, authID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log print event: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "printed"})
}

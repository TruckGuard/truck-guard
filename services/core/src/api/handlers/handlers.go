package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/truckguard/core/src/logic"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	datarepo "github.com/truckguard/core/src/repository/data"
	"github.com/truckguard/core/src/utils"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func HandlePlateEvent(c *gin.Context) {
	var event models.PlateEvent
	slog.Info("Plate event received", "event", event)

	if err := c.ShouldBindBodyWith(&event, binding.JSON); err != nil {
		slog.Error("Failed to bind plate event", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sysEventID, exists := c.Get("system_event_id"); exists {
		event.SystemEventID = sysEventID.(uint)
	}

	// Populate camera name if missing
	if event.CameraSourceName == "" && event.CameraID != "" {
		if cam, err := repository.GetCameraBySourceID(c.Request.Context(), event.CameraID); err == nil {
			event.CameraSourceName = cam.Name
		}
	}

	if err := repository.DB.WithContext(c.Request.Context()).Create(&event).Error; err != nil {
		slog.Error("Failed to save plate event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save event"})
		return
	}
	slog.Info("Plate event saved", "id", event.ID, "plate", event.Plate, "camera_id", event.CameraID)

	// Check if plate is excluded
	if excluded, comment := datarepo.IsPlateExcluded(c.Request.Context(), event.Plate); excluded {
		slog.Info("Plate is ignored, skipping matching logic", "plate", event.Plate, "comment", comment)
		c.JSON(http.StatusAccepted, gin.H{"status": "ignored", "id": event.ID})
		return
	}

	detachCtx := trace.ContextWithSpan(context.Background(), trace.SpanFromContext(c.Request.Context()))
	go logic.MatchPlateEvent(detachCtx, &event)

	c.JSON(http.StatusAccepted, gin.H{"status": "processing", "id": event.ID})
}

func HandleGetPlateEvents(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	plate := c.Query("plate")
	from := c.Query("from")
	to := c.Query("to")
	onlyUnlinked := c.Query("unlinked") == "true"

	events, total, err := repository.GetPlateEvents(c.Request.Context(), limit, offset, plate, from, to, onlyUnlinked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	utils.SendPaginatedResponse(c, events, total, page, limit)
}

func HandleGetPlateEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := repository.GetPlateEventByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Подію не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні події"})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

func HandleWeightEvent(c *gin.Context) {
	var event models.WeightEvent

	if err := c.ShouldBindBodyWith(&event, binding.JSON); err != nil {
		slog.Error("Failed to bind weight event", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sysEventID, exists := c.Get("system_event_id"); exists {
		event.SystemEventID = sysEventID.(uint)
	}

	// Populate scale name if missing
	if event.ScaleSourceName == "" && event.ScaleID != "" {
		if scale, err := repository.GetScaleBySourceID(c.Request.Context(), event.ScaleID); err == nil {
			event.ScaleSourceName = scale.Name
		}
	}

	if err := repository.DB.WithContext(c.Request.Context()).Create(&event).Error; err != nil {
		slog.Error("Failed to record weight", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record weight"})
		return
	}
	slog.Info("Weight event saved", "id", event.ID, "weight", event.Weight, "scale_id", event.ScaleID)

	detachCtx := trace.ContextWithSpan(context.Background(), trace.SpanFromContext(c.Request.Context()))
	go logic.MatchWeightEvent(detachCtx, &event)

	c.JSON(http.StatusAccepted, gin.H{"status": "weight_recorded", "id": event.ID})
}

func HandleGetWeightEvents(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	from := c.Query("from")
	to := c.Query("to")
	onlyUnlinked := c.Query("unlinked") == "true"

	events, total, err := repository.GetWeightEvents(c.Request.Context(), limit, offset, from, to, onlyUnlinked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weight events"})
		return
	}
	utils.SendPaginatedResponse(c, events, total, page, limit)
}

func HandleGetWeightEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := repository.GetWeightEventByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Подію зважування не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні події зважування"})
		}
		return
	}
	c.JSON(http.StatusOK, event)
}

func HandleGetSystemEvents(c *gin.Context) {
	limit, offset, page := utils.GetPagination(c)
	eventType := c.Query("type")
	from := c.Query("from")
	to := c.Query("to")

	events, total, err := repository.GetSystemEvents(c.Request.Context(), limit, offset, eventType, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch system events"})
		return
	}
	utils.SendPaginatedResponse(c, events, total, page, limit)
}

func HandleGetSystemEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := repository.GetSystemEventByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Системну подію не знайдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка при отриманні системної події"})
		}
		return
	}
	c.JSON(http.StatusOK, event)
}

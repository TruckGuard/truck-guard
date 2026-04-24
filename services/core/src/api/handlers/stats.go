package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

type DashboardStats struct {
	PermitsInZone     int64   `json:"permits_in_zone"`
	PermitsToday      int64   `json:"permits_today"`
	PermitsTodayClosed int64  `json:"permits_today_closed"`
	PlateEventsToday  int64   `json:"plate_events_today"`
	WeightEventsToday int64   `json:"weight_events_today"`
	PermitsVoidToday  int64   `json:"permits_void_today"`
	AvgDaysInZone     float64 `json:"avg_days_in_zone"`
}

func HandleGetStats(c *gin.Context) {
	ctx := c.Request.Context()
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var stats DashboardStats

	// Permits currently in zone (not closed, not void)
	repository.DB.WithContext(ctx).
		Model(&models.Permit{}).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Where("is_closed = ? AND is_void = ?", false, false).
		Count(&stats.PermitsInZone)

	// Permits created today
	repository.DB.WithContext(ctx).
		Model(&models.Permit{}).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Where("is_void = ? AND created_at >= ?", false, todayStart).
		Count(&stats.PermitsToday)

	// Permits closed today
	repository.DB.WithContext(ctx).
		Model(&models.Permit{}).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Where("is_closed = ? AND is_void = ? AND exit_time >= ?", true, false, todayStart).
		Count(&stats.PermitsTodayClosed)

	// Permits voided today
	repository.DB.WithContext(ctx).
		Model(&models.Permit{}).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Where("is_void = ? AND updated_at >= ?", true, todayStart).
		Count(&stats.PermitsVoidToday)

	// Plate events today
	repository.DB.WithContext(ctx).
		Model(&models.PlateEvent{}).
		Scopes(models.ScopeByPost(ctx, models.PlateEvent{}, "read")).
		Where("created_at >= ?", todayStart).
		Count(&stats.PlateEventsToday)

	// Weight events today
	repository.DB.WithContext(ctx).
		Model(&models.WeightEvent{}).
		Scopes(models.ScopeByPost(ctx, models.WeightEvent{}, "read")).
		Where("created_at >= ?", todayStart).
		Count(&stats.WeightEventsToday)

	// Average days in zone for permits currently in zone
	var avgResult struct {
		Avg *float64
	}
	repository.DB.WithContext(ctx).
		Model(&models.Permit{}).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Where("is_closed = ? AND is_void = ?", false, false).
		Select("AVG(COALESCE(days_in_zone, 0)) as avg").
		Scan(&avgResult)

	if avgResult.Avg != nil {
		stats.AvgDaysInZone = *avgResult.Avg
	}

	c.JSON(http.StatusOK, stats)
}

package logic

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer = otel.Tracer("matchmaker")
	meter  = otel.Meter("matchmaker")

	eventsMatchedCounter, _ = meter.Int64Counter("events_matched_total",
		metric.WithDescription("Number of processed plate/weight events"),
	)
	permitsCreatedCounter, _ = meter.Int64Counter("permits_created_total",
		metric.WithDescription("Number of new permits created"),
	)
)

const (
	defaultPermitTimeout = 60 * time.Second
)

// getPermitTimeout fetches the timeout from SystemSettings or returns default
func getPermitTimeout(ctx context.Context) time.Duration {
	var setting models.SystemSetting
	if err := repository.DB.WithContext(ctx).Where("key = ?", "permit_timeout_seconds").First(&setting).Error; err != nil {
		return defaultPermitTimeout
	}
	sec, err := strconv.Atoi(setting.Value)
	if err != nil {
		return defaultPermitTimeout
	}
	slog.Debug("Permit timeout", "timeout", sec)
	return time.Duration(sec) * time.Second
}

func GetOrCreatePermit(ctx context.Context, customsPostID uint, currentSourceID string) *models.Permit {
	db := repository.DB.WithContext(ctx)
	key := fmt.Sprintf("active_permit:customs_post:%d", customsPostID)
	lockKey := fmt.Sprintf("lock:active_permit:%d", customsPostID)

	// Try to get lock for 5 seconds to prevent race conditions
	lockOk, err := repository.RDB.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
	if err != nil || !lockOk {
		slog.Debug("Waiting for permit lock", "post_id", customsPostID)
		time.Sleep(200 * time.Millisecond)
		return GetOrCreatePermit(ctx, customsPostID, currentSourceID)
	}
	defer repository.RDB.Del(ctx, lockKey)

	permitID, _ := repository.RDB.Get(ctx, key).Uint64()
	var permit models.Permit

	if permitID > 0 {
		if err := db.Preload("PlateEvents").Preload("WeightEvents").First(&permit, permitID).Error; err == nil {
			timeout := getPermitTimeout(ctx)
			if !permit.IsClosed && !permit.IsVoid && permit.VerifiedBy == nil && time.Since(permit.LastActivityAt) < timeout {

				alreadySentSources := make(map[string]bool)
				for _, pe := range permit.PlateEvents {
					alreadySentSources[pe.CameraID] = true
				}
				for _, we := range permit.WeightEvents {
					alreadySentSources[we.ScaleID] = true
				}

				var triggerCameras []models.CameraConfig
				db.Where("customs_post_id = ? AND trigger_permit_creation = ?", customsPostID, true).Find(&triggerCameras)
				var triggerScales []models.ScaleConfig
				db.Where("customs_post_id = ? AND trigger_permit_creation = ?", customsPostID, true).Find(&triggerScales)

				allTriggered := true
				for _, cam := range triggerCameras {
					if !alreadySentSources[cam.SourceID] {
						allTriggered = false
						break
					}
				}
				if allTriggered {
					for _, scale := range triggerScales {
						if !alreadySentSources[scale.SourceID] {
							allTriggered = false
							break
						}
					}
				}

				isTriggerSource := false
				for _, cam := range triggerCameras {
					if cam.SourceID == currentSourceID {
						isTriggerSource = true
						break
					}
				}
				if !isTriggerSource {
					for _, scale := range triggerScales {
						if scale.SourceID == currentSourceID {
							isTriggerSource = true
							break
						}
					}
				}

				// Equipment counts
				var totalCams int64
				db.Model(&models.CameraConfig{}).Where("customs_post_id = ?", customsPostID).Count(&totalCams)
				var totalScales int64
				db.Model(&models.ScaleConfig{}).Where("customs_post_id = ?", customsPostID).Count(&totalScales)
				totalEvents := int64(len(permit.PlateEvents) + len(permit.WeightEvents))

				slog.Debug("Checking permit reuse",
					"permit_id", permit.ID,
					"source", currentSourceID,
					"already_sent", alreadySentSources[currentSourceID],
					"all_triggered", allTriggered,
					"is_trigger_source", isTriggerSource,
					"count", totalEvents,
					"limit", totalCams+totalScales,
				)

				// Decision Logic:
				// 1. If SAME SOURCE sent an event again -> New Permit
				// 2. If permit is ALREADY TRIGGERED and a NEW TRIGGER arrives -> New Permit
				// 3. If PERMIT IS FULL (all equipment reported) -> New Permit

				shouldCreateNew := false
				if alreadySentSources[currentSourceID] {
					shouldCreateNew = true
					slog.Info("Same source sent duplicate event, starting new permit", "permit_id", permit.ID, "source", currentSourceID)
				} else if isTriggerSource && allTriggered {
					shouldCreateNew = true
					slog.Info("All triggers met and new trigger event arrived, starting new permit", "permit_id", permit.ID, "source", currentSourceID)
				} else if totalEvents >= (totalCams + totalScales) {
					shouldCreateNew = true
					slog.Info("Permit is physically full, starting new", "permit_id", permit.ID, "events", totalEvents, "limit", totalCams+totalScales)
				}

				if !shouldCreateNew {
					slog.Debug("Reusing active permit", "permit_id", permit.ID)
					return &permit
				}
			}
		}
	}

	permit = models.Permit{
		CustomsPostID:  &customsPostID,
		Code:           fmt.Sprintf("PRM-%d-%d", customsPostID, time.Now().Unix()),
		EntryTime:      time.Now(),
		LastActivityAt: time.Now(),
	}

	if err := db.Create(&permit).Error; err != nil {
		slog.Error("Failed to create permit", "error", err)
		return nil
	}

	// Set in Redis with timeout
	timeout := getPermitTimeout(ctx)
	repository.RDB.Set(ctx, key, permit.ID, timeout)

	permitsCreatedCounter.Add(ctx, 1)
	slog.Info("Created new permit", "permit_id", permit.ID, "code", permit.Code)

	return &permit
}

func MatchPlateEvent(ctx context.Context, event *models.PlateEvent) {
	ctx, span := tracer.Start(ctx, "MatchPlateEvent",
		trace.WithAttributes(attribute.String("plate", event.Plate)))
	defer span.End()

	db := repository.DB
	if err := db.WithContext(ctx).Preload("Camera").First(event).Error; err != nil {
		slog.Error("Failed to load plate event", "error", err)
		span.RecordError(err)
		return
	}

	if event.Camera.CustomsPostID == nil {
		slog.Warn("Camera not assigned to a CustomsPost", "camera_id", event.CameraID)
		return
	}
	customsPostID := *event.Camera.CustomsPostID

	processEvent(ctx, customsPostID, event.CameraID, event, func(permit *models.Permit) {
		plate := event.Plate
		if event.PlateCorrected != "" {
			plate = event.PlateCorrected
		}

		switch event.Camera.Type {
		case "front":
			permit.PlateFront = plate
		case "back":
			permit.PlateBack = plate
		default:
			slog.Warn("Unknown camera type", "camera_type", event.Camera.Type)
			if permit.PlateFront == "" {
				permit.PlateFront = plate
			} else if permit.PlateBack == "" && permit.PlateFront != plate {
				permit.PlateBack = plate
			}
		}
	})

	eventsMatchedCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("type", "plate"),
	))
}

func MatchWeightEvent(ctx context.Context, event *models.WeightEvent) {
	ctx, span := tracer.Start(ctx, "MatchWeightEvent",
		trace.WithAttributes(attribute.Float64("weight", event.Weight)))
	defer span.End()

	db := repository.DB
	if err := db.WithContext(ctx).Preload("Scale").First(event).Error; err != nil {
		slog.Error("Failed to load weight event", "error", err)
		span.RecordError(err)
		return
	}

	if event.Scale.CustomsPostID == nil {
		slog.Warn("Scale not assigned to a CustomsPost", "scale_id", event.ScaleID)
		return
	}
	customsPostID := *event.Scale.CustomsPostID

	processEvent(ctx, customsPostID, event.ScaleID, event, func(permit *models.Permit) {
		if event.Weight > 0 {
			permit.TotalWeight = event.Weight
		}
	})

	eventsMatchedCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("type", "weight"),
	))
}

// processEvent handles the common logic: find active permit or create new, then apply updates
func processEvent(ctx context.Context, customsPostID uint, sourceID string, event interface{}, updateFn func(*models.Permit)) {
	db := repository.DB.WithContext(ctx)
	permit := GetOrCreatePermit(ctx, customsPostID, sourceID)
	if permit == nil {
		return
	}

	// Update existing
	updateFn(permit)
	permit.LastActivityAt = time.Now()
	if err := db.Save(permit).Error; err != nil {
		slog.Error("Failed to update permit", "permit_id", permit.ID, "error", err)
		return
	}
	slog.Info("Processed event for permit", "permit_id", permit.ID)

	// Refresh Redis TTL
	key := fmt.Sprintf("active_permit:customs_post:%d", customsPostID)
	timeout := getPermitTimeout(ctx)
	repository.RDB.Set(ctx, key, permit.ID, timeout)

	// Link Event to Permit
	switch e := event.(type) {
	case *models.PlateEvent:
		e.PermitID = &permit.ID
		db.Save(e)
	case *models.WeightEvent:
		e.PermitID = &permit.ID
		db.Save(e)
	}
}

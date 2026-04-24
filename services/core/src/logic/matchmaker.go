package logic

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/pkg/notify"
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
	val := repository.GetSystemSetting(ctx, "permit_timeout_seconds")
	if val == "" {
		return defaultPermitTimeout
	}
	sec, err := strconv.Atoi(val)
	if err != nil {
		return defaultPermitTimeout
	}
	slog.Debug("Permit timeout", "timeout", sec)
	return time.Duration(sec) * time.Second
}

func GetOrCreatePermit(ctx context.Context, customsPostID uint, currentSourceID string) (*models.Permit, bool) {
	db := repository.DB.WithContext(ctx)
	key := fmt.Sprintf("active_permit:customs_post:%d", customsPostID)
	lockKey := fmt.Sprintf("lock:active_permit:customs_post:%d", customsPostID)

	// checkExisting returns (permit, shouldReuse) — reuse=true means caller should return &permit, false
	checkExisting := func() (*models.Permit, bool) {
		permitID, _ := repository.RDB.Get(ctx, key).Uint64()
		if permitID == 0 {
			return nil, false
		}
		var p models.Permit
		if err := db.Preload("PlateEvents").Preload("WeightEvents").First(&p, permitID).Error; err != nil {
			return nil, false
		}
		if p.IsClosed || p.IsVoid || p.VerifiedBy != nil {
			return nil, false
		}

		alreadySentSources := make(map[string]bool)
		for _, pe := range p.PlateEvents {
			alreadySentSources[pe.CameraSourceID] = true
		}
		for _, we := range p.WeightEvents {
			alreadySentSources[we.ScaleSourceID] = true
		}

		var totalCams, totalScales int64
		db.Model(&models.CameraConfig{}).Where("customs_post_id = ?", customsPostID).Count(&totalCams)
		db.Model(&models.ScaleConfig{}).Where("customs_post_id = ?", customsPostID).Count(&totalScales)
		totalEvents := int64(len(p.PlateEvents) + len(p.WeightEvents))

		slog.Debug("Checking permit reuse",
			"permit_id", p.ID,
			"source", currentSourceID,
			"already_sent", alreadySentSources[currentSourceID],
			"count", totalEvents,
			"limit", totalCams+totalScales,
		)

		if alreadySentSources[currentSourceID] {
			slog.Info("Same source sent duplicate event, starting new permit", "permit_id", p.ID, "source", currentSourceID)
			return nil, false
		}
		if totalEvents >= (totalCams + totalScales) {
			slog.Info("Permit is physically full, starting new", "permit_id", p.ID, "events", totalEvents, "limit", totalCams+totalScales)
			return nil, false
		}

		slog.Debug("Reusing active permit", "permit_id", p.ID)
		return &p, true
	}

	// Fast path: check without lock first
	if p, reuse := checkExisting(); reuse {
		return p, false
	}

	// Need to create a new permit — acquire a short-lived distributed lock to prevent
	// duplicate creation when two events arrive simultaneously for the same post.
	const lockTTL = 10 * time.Second
	acquired, err := repository.RDB.SetNX(ctx, lockKey, 1, lockTTL).Result()
	if err != nil {
		slog.Error("Failed to acquire permit creation lock", "error", err)
		// Proceed anyway — worst case we create a duplicate, which is recoverable.
	}
	if !acquired {
		// Another goroutine is creating; wait briefly then re-check.
		time.Sleep(150 * time.Millisecond)
		if p, reuse := checkExisting(); reuse {
			return p, false
		}
		// Still no reusable permit — fall through to create (lock may have expired).
	} else {
		defer repository.RDB.Del(ctx, lockKey)
	}

	// Re-check under lock: between our fast-path check and acquiring the lock,
	// another goroutine may have already created a permit.
	if p, reuse := checkExisting(); reuse {
		return p, false
	}

	permit := models.Permit{
		CustomsPostID:  &customsPostID,
		Code:           fmt.Sprintf("%02d%06d", customsPostID, 0),
		EntryTime:      time.Now(),
		LastActivityAt: time.Now(),
	}

	if err := db.Create(&permit).Error; err != nil {
		slog.Error("Failed to create permit", "error", err)
		return nil, false
	}

	permit.Code = fmt.Sprintf("%02d%06d", customsPostID, permit.ID)
	if err := db.Save(&permit).Error; err != nil {
		slog.Error("Failed to update permit code", "error", err)
		return nil, false
	}

	// Publish the new permit to Redis before releasing the lock so any concurrent
	// goroutine that wakes after our sleep sees it.
	timeout := getPermitTimeout(ctx)
	repository.RDB.Set(ctx, key, permit.ID, timeout)

	permitsCreatedCounter.Add(ctx, 1)
	slog.Info("Created new permit", "permit_id", permit.ID, "code", permit.Code)

	return &permit, true
}

func MatchPlateEvent(ctx context.Context, event *models.PlateEvent) {
	ctx, span := tracer.Start(ctx, "MatchPlateEvent",
		trace.WithAttributes(attribute.String("plate", event.Plate)))
	defer span.End()

	db := repository.DB
	slog.Info("MatchPlateEvent: started", "event_id", event.ID, "plate", event.Plate)
	if err := db.WithContext(ctx).Preload("Camera").First(event).Error; err != nil {
		slog.Error("MatchPlateEvent: failed to load event", "error", err, "event_id", event.ID)
		span.RecordError(err)
		return
	}
	slog.Info("MatchPlateEvent: event loaded", "camera_id", event.CameraID)

	// Camera may have been deleted (FK SET NULL)
	if event.Camera == nil {
		slog.Warn("Camera was deleted, skipping permit matching", "camera_id", event.CameraID)
		return
	}

	if !event.Camera.MatchPermit {
		slog.Debug("Camera not configured to match permits", "camera_id", event.CameraID)
		return
	}

	if event.Camera.CustomsPostID == nil {
		slog.Warn("Camera not assigned to a CustomsPost", "camera_id", event.CameraID)
		return
	}
	customsPostID := *event.Camera.CustomsPostID

	sourceID := ""
	if event.CameraID != nil {
		sourceID = *event.CameraID
	}

	processEvent(ctx, customsPostID, sourceID, event, func(permit *models.Permit) {
		plate := event.Plate

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

	// Scale may have been deleted (FK SET NULL)
	if event.Scale == nil {
		slog.Warn("Scale was deleted, skipping permit matching", "scale_id", event.ScaleID)
		return
	}

	if !event.Scale.MatchPermit {
		slog.Debug("Scale not configured to match permits", "scale_id", event.ScaleID)
		return
	}

	if event.Scale.CustomsPostID == nil {
		slog.Warn("Scale not assigned to a CustomsPost", "scale_id", event.ScaleID)
		return
	}
	customsPostID := *event.Scale.CustomsPostID

	sourceID := ""
	if event.ScaleID != nil {
		sourceID = *event.ScaleID
	}

	processEvent(ctx, customsPostID, sourceID, event, func(permit *models.Permit) {
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
	slog.Info("processEvent: started", "source_id", sourceID)
	db := repository.DB.WithContext(ctx)
	permit, isNew := GetOrCreatePermit(ctx, customsPostID, sourceID)
	if permit == nil {
		slog.Warn("processEvent: GetOrCreatePermit returned nil")
		return
	}
	slog.Info("processEvent: permit obtained", "permit_id", permit.ID, "is_new", isNew)

	oldFront := permit.PlateFront
	oldBack := permit.PlateBack
	oldWeight := permit.TotalWeight

	// Update existing
	updateFn(permit)
	permit.LastActivityAt = time.Now()
	if err := db.Save(permit).Error; err != nil {
		slog.Error("Failed to update permit", "permit_id", permit.ID, "error", err)
		return
	}
	slog.Info("Processed event for permit", "permit_id", permit.ID)

	changes := make(map[string]interface{})
	if permit.PlateFront != oldFront {
		changes["plate_front"] = permit.PlateFront
	}
	if permit.PlateBack != oldBack {
		changes["plate_back"] = permit.PlateBack
	}
	if permit.TotalWeight != oldWeight {
		changes["total_weight"] = permit.TotalWeight
	}

	if isNew {
		repository.LogSystemPermitAudit(ctx, permit.ID, "create", changes, "Перепустка ініційована камерою/вагами")
		// Notify operators at this post
		plate := permit.PlateFront
		if plate == "" {
			plate = permit.PlateBack
		}
		slog.Info("Notifying operators about new auto-created permit", "permit_id", permit.ID, "post_id", customsPostID)
		notify.Global.PublishPermit(customsPostID, permit.ID, permit.Code, plate)
	} else if len(changes) > 0 {
		repository.LogSystemPermitAudit(ctx, permit.ID, "update", changes, "Дані оновлено подією з камери/ваг")
	}

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

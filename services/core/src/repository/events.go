package repository

import (
	"context"

	"github.com/truckguard/core/src/models"
)

func GetPlateEvents(ctx context.Context, limit, offset int, plate, from, to string) ([]models.PlateEvent, int64, error) {
	var events []models.PlateEvent
	var total int64

	query := DB.WithContext(ctx).Model(&models.PlateEvent{}).Scopes(models.ScopeByPost(ctx, models.PlateEvent{}, "read"))

	if plate != "" {
		query = query.Where("plate LIKE ?", "%"+plate+"%")
	}
	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func GetPlateEventByID(ctx context.Context, id string) (models.PlateEvent, error) {
	var event models.PlateEvent
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.PlateEvent{}, "read")).First(&event, id).Error; err != nil {
		return event, err
	}
	return event, nil
}

func GetWeightEvents(ctx context.Context, limit, offset int, from, to string) ([]models.WeightEvent, int64, error) {
	var events []models.WeightEvent
	var total int64

	query := DB.WithContext(ctx).Model(&models.WeightEvent{}).Scopes(models.ScopeByPost(ctx, models.WeightEvent{}, "read"))

	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func GetWeightEventByID(ctx context.Context, id string) (models.WeightEvent, error) {
	var event models.WeightEvent
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.WeightEvent{}, "read")).First(&event, id).Error; err != nil {
		return event, err
	}
	return event, nil
}

func GetSystemEvents(ctx context.Context, limit, offset int, eventType, from, to string) ([]models.SystemEvent, int64, error) {
	var events []models.SystemEvent
	var total int64

	query := DB.WithContext(ctx).Model(&models.SystemEvent{})

	if eventType != "" {
		query = query.Where("type = ?", eventType)
	}
	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}

	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func GetSystemEventByID(ctx context.Context, id string) (models.SystemEvent, error) {
	var event models.SystemEvent
	if err := DB.WithContext(ctx).First(&event, id).Error; err != nil {
		return event, err
	}
	return event, nil
}

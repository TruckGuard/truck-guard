package repository

import (
	"context"
	"log/slog"

	"github.com/truckguard/core/src/models"
)

func GetCameras(ctx context.Context, limit, offset int) ([]models.CameraConfig, int64, error) {
	var configs []models.CameraConfig
	var total int64

	query := DB.WithContext(ctx).Model(&models.CameraConfig{}).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "read"))
	slog.Debug("Get cameras", "query", query)

	if err := query.Count(&total).Error; err != nil {
		slog.Error("Failed to count cameras", "error", err)
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&configs).Error; err != nil {
		slog.Error("Failed to find cameras", "error", err)
		return nil, 0, err
	}

	return configs, total, nil
}

func GetCameraByID(ctx context.Context, id string) (models.CameraConfig, error) {
	var config models.CameraConfig
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "read")).Where("id = ?", id).First(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}

func GetCameraBySourceID(ctx context.Context, sourceID string) (models.CameraConfig, error) {
	var config models.CameraConfig
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "read")).Where("camera_id = ?", sourceID).First(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}

func CreateCamera(ctx context.Context, config *models.CameraConfig) error {
	return DB.WithContext(ctx).Create(config).Error
}

func UpdateCamera(ctx context.Context, config *models.CameraConfig) error {
	return DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "manage")).Save(config).Error
}

func DeleteCamera(ctx context.Context, config *models.CameraConfig) error {
	return DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "manage")).Delete(config).Error
}

func GetScales(ctx context.Context, limit, offset int) ([]models.ScaleConfig, int64, error) {
	var configs []models.ScaleConfig
	var total int64

	query := DB.WithContext(ctx).Model(&models.ScaleConfig{}).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "read"))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func GetScaleByID(ctx context.Context, id string) (models.ScaleConfig, error) {
	var config models.ScaleConfig
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "read")).Where("id = ?", id).First(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}

func GetScaleBySourceID(ctx context.Context, sourceID string) (models.ScaleConfig, error) {
	var config models.ScaleConfig
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "read")).Where("scale_id = ?", sourceID).First(&config).Error; err != nil {
		return config, err
	}
	return config, nil
}

func CreateScale(ctx context.Context, config *models.ScaleConfig) error {
	return DB.WithContext(ctx).Create(config).Error
}

func UpdateScale(ctx context.Context, config *models.ScaleConfig) error {
	return DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "manage")).Save(config).Error
}

func DeleteScale(ctx context.Context, config *models.ScaleConfig) error {
	return DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "manage")).Unscoped().Delete(config).Error
}

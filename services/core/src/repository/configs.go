package repository

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/truckguard/core/src/models"
)

func GetCameras(ctx context.Context, limit, offset int, search string) ([]models.CameraConfig, int64, error) {
	var configs []models.CameraConfig
	var total int64

	query := DB.WithContext(ctx).Model(&models.CameraConfig{}).Preload("CustomsPost").Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "read"))
	
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR camera_id ILIKE ?", searchPattern, searchPattern)
	}

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

func GetCameraByIdentifier(ctx context.Context, identifier string) (models.CameraConfig, error) {
	var config models.CameraConfig
	query := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.CameraConfig{}, "read"))
	
	idInt, err := strconv.Atoi(identifier)
	if err == nil {
		err = query.Where("id = ? OR camera_id = ?", idInt, identifier).First(&config).Error
	} else {
		err = query.Where("camera_id = ?", identifier).First(&config).Error
	}
	return config, err
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

func GetScales(ctx context.Context, limit, offset int, search string) ([]models.ScaleConfig, int64, error) {
	var configs []models.ScaleConfig
	var total int64

	query := DB.WithContext(ctx).Model(&models.ScaleConfig{}).Preload("CustomsPost").Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "read"))

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR scale_id ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func GetScaleByIdentifier(ctx context.Context, identifier string) (models.ScaleConfig, error) {
	var config models.ScaleConfig
	query := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.ScaleConfig{}, "read"))
	
	idInt, err := strconv.Atoi(identifier)
	if err == nil {
		err = query.Where("id = ? OR scale_id = ?", idInt, identifier).First(&config).Error
	} else {
		err = query.Where("scale_id = ?", identifier).First(&config).Error
	}
	return config, err
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

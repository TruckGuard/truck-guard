package repository

import (
	"context"
	"log/slog"

	"github.com/truckguard/core/src/models"
)

// GetSystemSetting returns the value of a system setting by key.
// Returns an empty string if not found.
func GetSystemSetting(ctx context.Context, key string) string {
	var setting models.SystemSetting
	if err := DB.WithContext(ctx).Where("key = ?", key).First(&setting).Error; err != nil {
		slog.Warn("System setting not found", "key", key)
		return ""
	}
	return setting.Value
}

// UpsertSystemSetting updates an existing setting or creates a new one.
func UpsertSystemSetting(ctx context.Context, key, value string) error {
	var setting models.SystemSetting
	if err := DB.WithContext(ctx).Where("key = ?", key).First(&setting).Error; err != nil {
		// Create
		setting = models.SystemSetting{
			Key:   key,
			Value: value,
		}
		return DB.WithContext(ctx).Create(&setting).Error
	}

	// Update
	setting.Value = value
	return DB.WithContext(ctx).Save(&setting).Error
}

// ListSystemSettings returns all system settings.
func ListSystemSettings(ctx context.Context) ([]models.SystemSetting, error) {
	var settings []models.SystemSetting
	if err := DB.WithContext(ctx).Find(&settings).Order("key ASC").Error; err != nil {
		return nil, err
	}
	return settings, nil
}

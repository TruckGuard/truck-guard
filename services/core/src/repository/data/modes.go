package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
	"gorm.io/gorm"
)

type ModeFilter struct {
	Code string
	Name string
}

func ListModes(ctx context.Context, limit, offset int, filter ModeFilter) ([]models.CustomsMode, int64, error) {
	var modes []models.CustomsMode
	var total int64

	query := repository.DB.WithContext(ctx).Model(&models.CustomsMode{})
	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("code asc").Find(&modes).Error; err != nil {
		return nil, 0, err
	}
	return modes, total, nil
}

func GetMode(ctx context.Context, id string) (models.CustomsMode, error) {
	var mode models.CustomsMode
	err := repository.DB.WithContext(ctx).First(&mode, id).Error
	return mode, err
}

func CreateMode(ctx context.Context, input *models.CustomsMode) error {
	return repository.DB.WithContext(ctx).Create(input).Error
}

func UpdateMode(ctx context.Context, id string, input *models.CustomsMode) (models.CustomsMode, error) {
	var mode models.CustomsMode
	if err := repository.DB.WithContext(ctx).First(&mode, id).Error; err != nil {
		return mode, err
	}
	mode.Name = input.Name
	mode.Code = input.Code
	mode.Description = input.Description
	mode.RequiredFields = input.RequiredFields
	if err := repository.DB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(&mode).Error; err != nil {
		return mode, err
	}
	return mode, nil
}

func DeleteMode(ctx context.Context, id string) error {
	return repository.DB.WithContext(ctx).Delete(&models.CustomsMode{}, id).Error
}

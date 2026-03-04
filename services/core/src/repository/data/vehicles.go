package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

type VehicleTypeFilter struct {
	Code string
	Name string
}

func ListVehicleTypes(ctx context.Context, limit, offset int, filter VehicleTypeFilter) ([]models.VehicleType, int64, error) {
	var types []models.VehicleType
	var total int64

	query := repository.DB.WithContext(ctx).Model(&models.VehicleType{})
	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("code asc").Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

func GetVehicleType(ctx context.Context, id string) (models.VehicleType, error) {
	var vt models.VehicleType
	err := repository.DB.WithContext(ctx).First(&vt, id).Error
	return vt, err
}

func CreateVehicleType(ctx context.Context, input *models.VehicleType) error {
	return repository.DB.WithContext(ctx).Create(input).Error
}

func UpdateVehicleType(ctx context.Context, id string, input *models.VehicleType) (models.VehicleType, error) {
	vt, err := GetVehicleType(ctx, id)
	if err != nil {
		return vt, err
	}
	vt.Name = input.Name
	vt.Code = input.Code
	vt.Description = input.Description
	vt.EntryPrice = input.EntryPrice
	vt.DailyPrice = input.DailyPrice
	vt.Color = input.Color
	if err := repository.DB.WithContext(ctx).Save(&vt).Error; err != nil {
		return vt, err
	}
	return vt, nil
}

func DeleteVehicleType(ctx context.Context, id string) error {
	return repository.DB.WithContext(ctx).Delete(&models.VehicleType{}, id).Error
}

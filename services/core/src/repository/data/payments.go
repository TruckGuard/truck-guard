package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

type PaymentTypeFilter struct {
	Code     string
	IsActive string
}

func ListPaymentTypes(ctx context.Context, limit, offset int, filter PaymentTypeFilter) ([]models.PaymentType, int64, error) {
	var types []models.PaymentType
	var total int64

	query := repository.DB.WithContext(ctx).Model(&models.PaymentType{})
	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.IsActive != "" {
		query = query.Where("is_active = ?", filter.IsActive == "true")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("code asc").Find(&types).Error; err != nil {
		return nil, 0, err
	}
	return types, total, nil
}

func GetPaymentType(ctx context.Context, id string) (models.PaymentType, error) {
	var pt models.PaymentType
	err := repository.DB.WithContext(ctx).First(&pt, id).Error
	return pt, err
}

func CreatePaymentType(ctx context.Context, input *models.PaymentType) error {
	return repository.DB.WithContext(ctx).Create(input).Error
}

func UpdatePaymentType(ctx context.Context, id string, input *models.PaymentType) (models.PaymentType, error) {
	pt, err := GetPaymentType(ctx, id)
	if err != nil {
		return pt, err
	}
	pt.Name = input.Name
	pt.Code = input.Code
	pt.Description = input.Description
	pt.IsActive = input.IsActive
	pt.Icon = input.Icon
	if err := repository.DB.WithContext(ctx).Save(&pt).Error; err != nil {
		return pt, err
	}
	return pt, nil
}

func DeletePaymentType(ctx context.Context, id string) error {
	return repository.DB.WithContext(ctx).Delete(&models.PaymentType{}, id).Error
}

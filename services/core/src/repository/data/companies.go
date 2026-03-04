package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

type CompanyFilter struct {
	Name   string
	EDRPOU string
}

func ListCompanies(ctx context.Context, limit, offset int, filter CompanyFilter) ([]models.Company, int64, error) {
	var companies []models.Company
	var total int64

	query := repository.DB.WithContext(ctx).Model(&models.Company{})
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.EDRPOU != "" {
		query = query.Where("edrpou ILIKE ?", "%"+filter.EDRPOU+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Limit(limit).Offset(offset).Order("name asc").Find(&companies).Error; err != nil {
		return nil, 0, err
	}
	return companies, total, nil
}

func GetCompany(ctx context.Context, id string) (models.Company, error) {
	var company models.Company
	err := repository.DB.WithContext(ctx).First(&company, id).Error
	return company, err
}

func CreateCompany(ctx context.Context, input *models.Company) error {
	return repository.DB.WithContext(ctx).Create(input).Error
}

func UpdateCompany(ctx context.Context, id string, input *models.Company) (models.Company, error) {
	company, err := GetCompany(ctx, id)
	if err != nil {
		return company, err
	}
	company.Name = input.Name
	company.EDRPOU = input.EDRPOU
	company.Details = input.Details
	if err := repository.DB.WithContext(ctx).Save(&company).Error; err != nil {
		return company, err
	}
	return company, nil
}

func DeleteCompany(ctx context.Context, id string) error {
	return repository.DB.WithContext(ctx).Delete(&models.Company{}, id).Error
}

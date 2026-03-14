package datarepo

import (
	"context"

	"github.com/truckguard/core/src/models"
	"github.com/truckguard/core/src/repository"
)

func IsPlateExcluded(ctx context.Context, plate string) (bool, string) {
	var excluded models.ExcludedPlate
	if err := repository.DB.WithContext(ctx).Where("plate = ?", plate).First(&excluded).Error; err == nil {
		return true, excluded.Comment
	}
	return false, ""
}

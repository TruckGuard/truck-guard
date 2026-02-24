package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/truckguard/core/src/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func GetPermits(ctx context.Context, limit, offset int, plate string) ([]models.Permit, int64, error) {
	var permits []models.Permit
	var total int64

	query := DB.WithContext(ctx).Model(&models.Permit{}).Scopes(models.ScopeByPost(ctx, models.Permit{}, "read"))

	if plate != "" {
		query = query.Where("plate_front = ? OR plate_back = ?", plate, plate)
	}

	var countQuery = query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").
		Preload("CustomsPost").
		Preload("PlateEvents").
		Preload("WeightEvents").
		Find(&permits).Error; err != nil {
		return nil, 0, err
	}

	return permits, total, nil
}

func GetPermitByID(ctx context.Context, id string) (models.Permit, error) {
	var permit models.Permit
	if err := DB.WithContext(ctx).
		Scopes(models.ScopeByPost(ctx, models.Permit{}, "read")).
		Preload("CustomsPost").
		Preload("PlateEvents").
		Preload("WeightEvents").
		First(&permit, id).Error; err != nil {
		return permit, err
	}
	return permit, nil
}

func CreatePermit(ctx context.Context, input *models.Permit, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)
	if user.ID != 0 {
		input.CreatedBy = &user.ID
		if input.CustomsPostID == nil && user.CustomsPostID != nil {
			input.CustomsPostID = user.CustomsPostID
		}
	}

	input.EntryTime = time.Now()
	input.LastActivityAt = time.Now()

	if err := DB.WithContext(ctx).Create(input).Error; err != nil {
		return err
	}

	LogPermitAudit(ctx, input.ID, user.ID, "create", map[string]interface{}{"source": "manual"}, "Створено вручну")
	return nil
}

func UpdatePermit(ctx context.Context, id string, input map[string]interface{}, authID string) (models.Permit, error) {
	var permit models.Permit
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.Permit{}, "update")).Preload("CustomsPost").First(&permit, id).Error; err != nil {
		return permit, err
	}

	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	if err := DB.WithContext(ctx).Model(&permit).Updates(input).Error; err != nil {
		return permit, err
	}

	if val, ok := input["verified_by"]; ok && val != nil {
		now := time.Now()
		DB.WithContext(ctx).Model(&permit).Update("verified_at", now)
	}

	LogPermitAudit(ctx, permit.ID, user.ID, "update", input, "Оновлено оператором")

	return permit, nil
}

func ValidatePermit(ctx context.Context, id string, authID string) (models.Permit, error) {
	var permit models.Permit
	var user models.User

	if err := DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user).Error; err != nil {
		return permit, err
	}

	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.Permit{}, "validate")).First(&permit, id).Error; err != nil {
		return permit, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"verified_by": user.ID,
		"verified_at": now,
	}

	if err := DB.WithContext(ctx).Model(&permit).Updates(updates).Error; err != nil {
		return permit, err
	}

	LogPermitAudit(ctx, permit.ID, user.ID, "validate", updates, "Валідовано користувачем")

	DB.WithContext(ctx).Preload("Verifier").First(&permit, id)
	return permit, nil
}

func LogPermitAudit(ctx context.Context, permitID uint, userID uint, action string, changes interface{}, comment string) {
	jsonBytes, _ := json.Marshal(changes)
	audit := models.PermitAudit{
		PermitID: permitID,
		UserID:   &userID,
		Action:   action,
		Changes:  datatypes.JSON(jsonBytes),
		Comment:  comment,
	}
	DB.WithContext(ctx).Create(&audit)
}

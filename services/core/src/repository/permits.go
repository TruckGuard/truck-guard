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
		Preload("Verifier").
		Preload("VehicleType").
		Preload("CustomsMode").
		Preload("PaymentType").
		Preload("Payers").
		Preload("Creator").
		Preload("ResponsibleUser").
		Preload("AuditEvents", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("AuditEvents.User").
		// Preload("Payers.Company").
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
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.Permit{}, "manage")).Preload("CustomsPost").First(&permit, id).Error; err != nil {
		return permit, err
	}

	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	// Copy changes for audit trace before modification
	auditChanges := make(map[string]interface{})
	for k, v := range input {
		auditChanges[k] = v
	}

	action := "update"
	comment := "Оновлено оператором"

	// Recalculate financials if exit time or status changes
	if isClosed, ok := input["is_closed"].(bool); ok && isClosed {
		action = "close"
		comment = "Перепустку закрито"
		if permit.ExitTime == nil {
			now := time.Now()
			input["exit_time"] = now
			permit.ExitTime = &now
		}
		if err := CalculateFinancials(ctx, &permit, input); err != nil {
			return permit, err
		}
	} else if _, hasExitTime := input["exit_time"]; hasExitTime {
		if err := CalculateFinancials(ctx, &permit, input); err != nil {
			return permit, err
		}
	}

	// Handle nested CustomsData in the input map by flattening it for GORM
	if cdMap, ok := input["customs_data"].(map[string]interface{}); ok {
		// Only flatten fields that exist in the PermitCustomsData struct
		allowedFields := map[string]bool{
			"declarant":  true,
			"goods":      true,
			"sender":     true,
			"receiver":   true,
			"vmd_number": true,
		}
		for k, v := range cdMap {
			if allowedFields[k] {
				input["customs_"+k] = v
			}
		}
		delete(input, "customs_data")
	}

	// Financials are already calculated above if needed

	if len(input) > 0 {
		if err := DB.WithContext(ctx).Model(&permit).Updates(input).Error; err != nil {
			return permit, err
		}
	}

	if val, ok := input["verified_by"]; ok && val != nil {
		now := time.Now()
		DB.WithContext(ctx).Model(&permit).Update("verified_at", now)
	}

	if input["responsible_user_id"] != nil {
		DB.WithContext(ctx).Model(&permit).Update("responsible_user_id", input["responsible_user_id"])
	}

	LogPermitAudit(ctx, permit.ID, user.ID, action, auditChanges, comment)

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

func LogSystemPermitAudit(ctx context.Context, permitID uint, action string, changes interface{}, comment string) {
	jsonBytes, _ := json.Marshal(changes)
	audit := models.PermitAudit{
		PermitID: permitID,
		UserID:   nil,
		Action:   action,
		Changes:  datatypes.JSON(jsonBytes),
		Comment:  comment,
	}
	DB.WithContext(ctx).Create(&audit)
}

func CalculateFinancials(ctx context.Context, permit *models.Permit, updates map[string]interface{}) error {
	// 1. Get ExitTime
	var exitTime time.Time
	if etVal, ok := updates["exit_time"]; ok {
		switch v := etVal.(type) {
		case time.Time:
			exitTime = v
		case *time.Time:
			if v != nil {
				exitTime = *v
			}
		case string:
			parsed, _ := time.Parse(time.RFC3339, v)
			exitTime = parsed
		}
	} else if permit.ExitTime != nil {
		exitTime = *permit.ExitTime
	} else {
		return nil // Cannot calculate without exit time
	}

	// 2. Load necessary relations
	if err := DB.WithContext(ctx).Preload("VehicleType").Preload("Payers.Company").First(permit, permit.ID).Error; err != nil {
		return err
	}

	// 3. Calculate Days In Zone (each partial day counts as full day)
	duration := exitTime.Sub(permit.EntryTime)
	days := int(duration.Hours() / 24)
	if duration.Hours() > float64(days*24) || duration <= 0 {
		days++
	}
	if days < 1 {
		days = 1 // Minimum 1 day
	}
	updates["days_in_zone"] = days

	// 4. Calculate Exist Fee (Sum 2)
	dailyPrice := 0.0
	if permit.VehicleType != nil {
		dailyPrice = permit.VehicleType.DailyPrice
	}
	exitFee := float64(days) * dailyPrice
	updates["exit_fee"] = exitFee

	// 5. Calculate Discount (use Payer 1 if exists, otherwise 0)
	discountPct := 0.0
	discountFixed := 0.0
	if len(permit.Payers) > 0 {
		for _, payer := range permit.Payers {
			if payer.SlotIndex == 1 && payer.Company != nil {
				discountPct = payer.Company.DiscountPercentage
				discountFixed = payer.Company.DiscountFixed
				break
			}
		}
	}

	// Sum 1 is EntryFee, Sum 2 is ExitFee
	subtotal := permit.EntryFee + exitFee

	// Default to applying fixed first, then percentage, or vice versa depending on logic.
	// For now: Total = (Subtotal - Fixed) * (1 - Pct/100)
	discountAmount := 0.0

	// apply fixed
	if discountFixed > 0 {
		if subtotal >= discountFixed {
			discountAmount += discountFixed
			subtotal -= discountFixed
		} else {
			discountAmount += subtotal
			subtotal = 0
		}
	}

	// apply percentage
	if discountPct > 0 && discountPct <= 100 {
		pctDiscount := subtotal * (discountPct / 100.0)
		discountAmount += pctDiscount
		subtotal -= pctDiscount
	}

	updates["discount_amount"] = discountAmount
	updates["total_sum"] = subtotal

	return nil
}

func GetPermitAuditEvents(ctx context.Context, permitID string) ([]models.PermitAudit, error) {
	var audits []models.PermitAudit
	err := DB.WithContext(ctx).
		Preload("User").
		Where("permit_id = ?", permitID).
		Order("created_at desc").
		Find(&audits).Error
	return audits, err
}

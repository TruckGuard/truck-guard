package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/truckguard/core/src/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CustomFilter represents a single dynamic filter condition
type CustomFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, neq, gt, lt, gte, lte, contains
	Value    string `json:"value"`
}

// PermitQueryParams holds all optional server-side query parameters for GetPermits.
type PermitQueryParams struct {
	Plate             string
	IsClosed          *bool
	IsVoid            *bool
	SortField         string // validated against allowedSortFields allowlist
	SortOrder         string // "asc" or "desc"
	FilterFrom        string // RFC3339
	FilterTo          string // RFC3339
	FilterPostID      string
	FilterVehicleType string
	FilterPaymentType string
	FilterPayer       string // search by company name or edrpou
	Search            string // search across code, plate_front, plate_back
	CustomFilters     string // JSON encoded array of CustomFilter
}

// allowedSortFields maps safe UI field name to actual SQL column — prevents injection.
var allowedSortFields = map[string]string{
	"id":                 "permits.id",
	"entry_time":         "permits.entry_time",
	"exit_time":          "permits.exit_time",
	"plate_front":        "permits.plate_front",
	"customs_post_id":    "permits.customs_post_id",
	"total_weight":       "permits.total_weight",
	"vehicle_type_id":    "permits.vehicle_type_id",
	"payment_type_id":    "permits.payment_type_id",
	"is_closed":          "permits.is_closed",
	"customs_mode":       "permits.customs_mode_code",
	"declaration_number": "permits.declaration_number",
	"verified_at":        "permits.verified_at",
	"entry_fee":          "permits.entry_fee",
	"exit_fee":           "permits.exit_fee",
	"total_sum":          "permits.total_sum",
	"days_in_zone":       "permits.days_in_zone",
}

// allowedFilterFields prevents SQL injection by restricting what fields can be filtered dynamically
var allowedFilterFields = map[string]string{
	"id":                            "permits.id",
	"entry_time":                    "permits.entry_time",
	"exit_time":                     "permits.exit_time",
	"plate_front":                   "permits.plate_front",
	"plate_back":                    "permits.plate_back",
	"customs_post_id":               "permits.customs_post_id",
	"total_weight":                  "permits.total_weight",
	"vehicle_type_id":               "permits.vehicle_type_id",
	"payment_type_id":               "permits.payment_type_id",
	"is_closed":                     "permits.is_closed",
	"customs_mode_id":               "permits.customs_mode_id",
	"customs_mode":                  "permits.customs_mode_code",
	"declaration_number":            "permits.declaration_number",
	"customs_declarant_name":        "permits.customs_declarant_name",
	"customs_commodity_description": "permits.customs_commodity_description",
	"customs_vmd_number":            "permits.customs_vmd_number",
	"customs_sender":                "permits.customs_sender",
	"customs_receiver":              "permits.customs_receiver",
	"verified_at":                   "permits.verified_at",
	"entry_fee":                     "permits.entry_fee",
	"exit_fee":                      "permits.exit_fee",
	"total_sum":                     "permits.total_sum",
	"days_in_zone":                  "permits.days_in_zone",
}

func GetPermits(ctx context.Context, limit, offset int, params PermitQueryParams) ([]models.Permit, int64, error) {
	var permits []models.Permit
	var total int64

	query := DB.WithContext(ctx).Model(&models.Permit{}).Scopes(models.ScopeByPost(ctx, models.Permit{}, "read"))

	// --- Filters ---

	if params.Plate != "" {
		query = query.Where("permits.plate_front = ? OR permits.plate_back = ?", params.Plate, params.Plate)
	}

	slog.Debug("Params", "params", params)
	if params.IsClosed != nil {
		query = query.Where("permits.is_closed = ?", *params.IsClosed)
	}
	if params.IsVoid != nil {
		query = query.Where("permits.is_void = ?", *params.IsVoid)
	}

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("permits.code ILIKE ? OR permits.plate_front ILIKE ? OR permits.plate_back ILIKE ?", search, search, search)
	}

	if params.FilterPostID != "" {
		query = query.Where("permits.customs_post_id = ?", params.FilterPostID)
	}

	if params.FilterVehicleType != "" {
		query = query.Where("permits.vehicle_type_id = ?", params.FilterVehicleType)
	}

	if params.FilterFrom != "" {
		if t, err := time.Parse(time.RFC3339, params.FilterFrom); err == nil {
			query = query.Where("permits.entry_time >= ?", t)
		}
	}

	if params.FilterTo != "" {
		if t, err := time.Parse(time.RFC3339, params.FilterTo); err == nil {
			query = query.Where("permits.entry_time <= ?", t)
		}
	}

	if params.FilterPaymentType != "" {
		query = query.Where("permits.payment_type_id = ?", params.FilterPaymentType)
	}

	if params.FilterPayer != "" {
		// JOIN against permit_payers and companies to filter by payer name/edrpou.
		// The join is scoped to avoid duplicates via a subquery.
		payerSearch := "%" + params.FilterPayer + "%"
		query = query.Where(
			"permits.id IN (SELECT pp.permit_id FROM permit_payers pp JOIN companies c ON c.id = pp.company_id WHERE pp.deleted_at IS NULL AND (c.name ILIKE ? OR c.edrpou ILIKE ?))",
			payerSearch, payerSearch,
		)
	}

	// --- Custom Filters (Advanced Array) ---
	if params.CustomFilters != "" {
		var customFilters []CustomFilter
		slog.Debug("custom filters", "custom_filters", params.CustomFilters)
		if err := json.Unmarshal([]byte(params.CustomFilters), &customFilters); err == nil {
			for _, f := range customFilters {
				// Prevent SQL injection by strictly matching allowed fields
				if col, ok := allowedFilterFields[f.Field]; ok {
					var val interface{} = f.Value

					// Type casting based on the database column type to prevent Postgres operator errors
					if col == "permits.is_closed" {
						val = f.Value == "true"
					} else if col == "permits.customs_post_id" || col == "permits.vehicle_type_id" || col == "permits.payment_type_id" || col == "permits.verified_by" || col == "EXTRACT(DAY FROM (NOW() - permits.entry_time))" {
						if vInt, err := strconv.Atoi(f.Value); err == nil {
							val = vInt
						} else if f.Operator != "isnull" && f.Operator != "notnull" {
							continue
						}
					} else if col == "permits.total_weight" || col == "permits.entry_fee" || col == "permits.exit_fee" || col == "permits.total_sum" {
						if vFloat, err := strconv.ParseFloat(f.Value, 64); err == nil {
							val = vFloat
						} else if f.Operator != "isnull" && f.Operator != "notnull" {
							continue
						}
					}

					switch f.Operator {
					case "eq":
						if col == "permits.entry_time" || col == "permits.exit_time" {
							strVal := fmt.Sprintf("%v", val)
							if len(strVal) == 10 {
								// Format YYYY-MM-DD -> match the entire day
								query = query.Where(fmt.Sprintf("%s >= ? AND %s <= ?", col, col), strVal+" 00:00:00", strVal+" 23:59:59")
							} else if len(strVal) == 16 {
								// Format YYYY-MM-DDTHH:MM -> match the entire minute
								// Postgres accepts 'T' or space
								baseStr := strVal[:10] + " " + strVal[11:]
								query = query.Where(fmt.Sprintf("%s >= ? AND %s <= ?", col, col), baseStr+":00", baseStr+":59")
							} else {
								query = query.Where(fmt.Sprintf("%s = ?", col), val)
							}
						} else {
							query = query.Where(fmt.Sprintf("%s = ?", col), val)
						}
					case "neq":
						query = query.Where(fmt.Sprintf("%s != ?", col), val)
					case "gt":
						query = query.Where(fmt.Sprintf("%s > ?", col), val)
					case "lt":
						query = query.Where(fmt.Sprintf("%s < ?", col), val)
					case "gte":
						query = query.Where(fmt.Sprintf("%s >= ?", col), val)
					case "lte":
						query = query.Where(fmt.Sprintf("%s <= ?", col), val)
					case "contains":
						query = query.Where(fmt.Sprintf("%s ILIKE ?", col), "%"+f.Value+"%")
					case "isnull":
						query = query.Where(fmt.Sprintf("%s IS NULL", col))
					case "notnull":
						query = query.Where(fmt.Sprintf("%s IS NOT NULL", col))
					}
				}
			}
		} else {
			fmt.Printf("Error unmarshalling CustomFilters: %v\\n", err)
		}
	}

	// --- Count (before ORDER / LIMIT) ---

	var countQuery = query.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// --- Sorting ---

	orderClause := "permits.entry_time DESC" // default
	if col, ok := allowedSortFields[params.SortField]; ok {
		dir := "ASC"
		if params.SortOrder == "desc" {
			dir = "DESC"
		}
		orderClause = fmt.Sprintf("%s %s", col, dir)
	}

	if err := query.
		Limit(limit).Offset(offset).
		Order(orderClause).
		Preload("CustomsPost").
		Preload("VehicleType").
		Preload("CustomsMode").
		Preload("PaymentType").
		Preload("ResponsibleUser").
		Preload("Verifier").
		Preload("Payers", func(db *gorm.DB) *gorm.DB {
			return db.Order("slot_index asc")
		}).
		Preload("Payers.Company").
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
		Preload("ClosedBy").
		Preload("VoidedBy").
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

	// Generate code: {PostID}{PermitID}
	postID := uint(0)
	if input.CustomsPostID != nil {
		postID = *input.CustomsPostID
	}
	input.Code = fmt.Sprintf("%02d%06d", postID, input.ID)
	DB.WithContext(ctx).Model(input).Update("code", input.Code)

	LogPermitAudit(ctx, input.ID, user.ID, "create", map[string]interface{}{"source": "manual", "code": input.Code}, "Створено вручну")
	return nil
}

func UpdatePermit(ctx context.Context, id string, input map[string]interface{}, authID string) (models.Permit, error) {
	var permit models.Permit
	if err := DB.WithContext(ctx).Scopes(models.ScopeByPost(ctx, models.Permit{}, "manage")).Preload("CustomsPost").First(&permit, id).Error; err != nil {
		return permit, err
	}

	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	// If vehicle_type_id is provided, automatically populate entry_fee and daily_fee
	if vtIDRaw, ok := input["vehicle_type_id"]; ok && vtIDRaw != nil {
		var vtID uint
		switch v := vtIDRaw.(type) {
		case float64:
			vtID = uint(v)
		case int:
			vtID = uint(v)
		case string:
			if id, err := strconv.Atoi(v); err == nil {
				vtID = uint(id)
			}
		}

		if vtID > 0 {
			var vt models.VehicleType
			if err := DB.WithContext(ctx).First(&vt, vtID).Error; err == nil {
				// Only populate if not explicitly provided in the input, or always? 
				// User said "also price for entry and cost per day are filled". 
				// Usually this means we want to override with defaults from VehicleType.
				input["entry_fee"] = vt.EntryPrice
				input["daily_fee"] = vt.DailyPrice
			}
		}
	}

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
		if user.ID != 0 {
			input["closed_by_id"] = user.ID
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
	// Handle Payers association explicitly before Updates()
	if payersVal, ok := input["payers"]; ok {
		var inputPayers []models.PermitPayer
		// Convert from interface to specific slice
		payerBytes, _ := json.Marshal(payersVal)
		if err := json.Unmarshal(payerBytes, &inputPayers); err == nil {
			// Start transaction for atomic replacement
			tx := DB.WithContext(ctx).Begin()

			// We MUST use Unscoped() here because PermitPayer has a uniqueIndex on (permit_id, slot_index).
			// Regular Delete() only soft-deletes (sets deleted_at), which causes the uniqueIndex to collide
			// when we try to Create() new records with the same slot indices.
			if err := tx.Unscoped().Where("permit_id = ?", permit.ID).Delete(&models.PermitPayer{}).Error; err != nil {
				tx.Rollback()
				return permit, err
			}

			for i := range inputPayers {
				inputPayers[i].ID = 0 // Ensure it's treated as a new record
				inputPayers[i].PermitID = permit.ID
				if err := tx.Create(&inputPayers[i]).Error; err != nil {
					tx.Rollback()
					return permit, err
				}
			}

			if err := tx.Commit().Error; err != nil {
				return permit, err
			}
		}
		delete(input, "payers")
	}

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
	if permit.DailyFee > 0 {
		dailyPrice = permit.DailyFee
	} else if permit.VehicleType != nil {
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

// VoidUnverifiedPermits voids non-verified permits older than specified hours.
func VoidUnverifiedPermits(ctx context.Context, hours int) (int64, error) {
	if hours <= 0 {
		return 0, nil
	}

	var permits []models.Permit
	err := DB.WithContext(ctx).
		Where("is_void = false AND verified_by IS NULL AND created_at < ?", time.Now().Add(-time.Duration(hours)*time.Hour)).
		Find(&permits).Error
	if err != nil {
		return 0, err
	}

	if len(permits) == 0 {
		return 0, nil
	}

	var ids []uint
	for _, p := range permits {
		ids = append(ids, p.ID)
	}

	res := DB.WithContext(ctx).Model(&models.Permit{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"is_void":    true,
			"updated_at": time.Now(),
		})

	if res.Error == nil && res.RowsAffected > 0 {
		for _, p := range permits {
			LogSystemPermitAudit(ctx, p.ID, "void", nil, "Перепустку анульовано автоматично (таймаут)")
		}
	}

	return res.RowsAffected, res.Error
}

// CleanupAuditLogs deletes audit logs for closed permits older than specified days.
func CleanupAuditLogs(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, nil
	}
	res := DB.WithContext(ctx).
		Where("permit_id IN (SELECT id FROM permits WHERE is_closed = true) AND created_at < ?", time.Now().Add(-time.Duration(days)*24*time.Hour)).
		Delete(&models.PermitAudit{})
	return res.RowsAffected, res.Error
}

// RestorePermit restores a voided permit.
func RestorePermit(ctx context.Context, id string, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	res := DB.WithContext(ctx).Model(&models.Permit{}).
		Where("id = ? AND is_void = true", id).
		Updates(map[string]interface{}{
			"is_void":    false,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}

	permitID, _ := strconv.Atoi(id)
	LogPermitAudit(ctx, uint(permitID), user.ID, "restore", nil, "Перепустку відновлено")
	return nil
}

// VoidPermit marks a permit as voided.
func VoidPermit(ctx context.Context, id string, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	res := DB.WithContext(ctx).Model(&models.Permit{}).
		Where("id = ? AND is_void = false", id).
		Updates(map[string]interface{}{
			"is_void":        true,
			"voided_by_id":   user.ID,
			"updated_at":     time.Now(),
		})
	if res.Error != nil {
		return res.Error
	}

	permitID, _ := strconv.Atoi(id)
	LogPermitAudit(ctx, uint(permitID), user.ID, "void", nil, "Перепустку анульовано вручну")
	return nil
}

// DeletePermit permanently deletes a permit and all its related events/data.
func DeletePermit(ctx context.Context, id string, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Delete associated event links (Audit, PlateEvents, WeightEvents, Payers)
		// Note: We are doing a hard delete (Unscoped) for these as well if the parent is hard deleted
		if err := tx.Unscoped().Where("permit_id = ?", id).Delete(&models.PermitAudit{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("permit_id = ?", id).Delete(&models.PermitPayer{}).Error; err != nil {
			return err
		}

		// For PlateEvents and WeightEvents, we might want to just un-link them 
		// OR delete them if they were created specifically for this permit.
		// Given user's request "events from permit also deleted", we delete them.
		if err := tx.Unscoped().Where("permit_id = ?", id).Delete(&models.PlateEvent{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("permit_id = ?", id).Delete(&models.WeightEvent{}).Error; err != nil {
			return err
		}

		// 2. Finally hard delete the permit
		if err := tx.Unscoped().Delete(&models.Permit{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// CleanupVoidedPermits permanently deletes voided permits older than specified days.
func CleanupVoidedPermits(ctx context.Context, days int) (int64, error) {
	if days <= 0 {
		return 0, nil
	}
	
	now := time.Now()
	threshold := now.Add(-time.Duration(days) * 24 * time.Hour)

	var ids []uint
	if err := DB.WithContext(ctx).Model(&models.Permit{}).
		Where("is_void = true AND updated_at < ?", threshold).
		Pluck("id", &ids).Error; err != nil {
		return 0, err
	}

	if len(ids) == 0 {
		return 0, nil
	}

	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Unscoped().Where("permit_id IN ?", ids).Delete(&models.PermitAudit{})
		tx.Unscoped().Where("permit_id IN ?", ids).Delete(&models.PermitPayer{})
		tx.Unscoped().Where("permit_id IN ?", ids).Delete(&models.PlateEvent{})
		tx.Unscoped().Where("permit_id IN ?", ids).Delete(&models.WeightEvent{})
		return tx.Unscoped().Where("id IN ?", ids).Delete(&models.Permit{}).Error
	})

	return int64(len(ids)), err
}

// LinkPermitEvent links a permit to a plate or weight event.
func LinkPermitEvent(ctx context.Context, permitID uint, eventID uint, eventType string, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var permit models.Permit
		if err := tx.First(&permit, permitID).Error; err != nil {
			return fmt.Errorf("permit not found: %w", err)
		}

		var action, comment string
		changes := map[string]interface{}{
			"event_id":   eventID,
			"event_type": eventType,
		}

		if eventType == "plate" {
			var event models.PlateEvent
			if err := tx.First(&event, eventID).Error; err != nil {
				return fmt.Errorf("plate event not found: %w", err)
			}
			if event.PermitID != nil && *event.PermitID == permitID {
				return nil // Already linked
			}
			event.PermitID = &permitID
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
			action = "link_plate"
			comment = fmt.Sprintf("Прив'язано подію камери #%d", eventID)
			
			// If permit has no plate, update it from event
			if permit.PlateFront == "" {
				tx.Model(&permit).Update("plate_front", event.Plate)
			}
		} else if eventType == "weight" {
			var event models.WeightEvent
			if err := tx.First(&event, eventID).Error; err != nil {
				return fmt.Errorf("weight event not found: %w", err)
			}
			if event.PermitID != nil && *event.PermitID == permitID {
				return nil // Already linked
			}
			event.PermitID = &permitID
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
			action = "link_weight"
			comment = fmt.Sprintf("Прив'язано подію ваг #%d", eventID)

			// If permit has no weight, update it from event
			if permit.TotalWeight == 0 {
				tx.Model(&permit).Update("total_weight", event.Weight)
			}
		} else {
			return fmt.Errorf("invalid event type: %s", eventType)
		}

		LogPermitAudit(ctx, permitID, user.ID, action, changes, comment)
		return nil
	})
}

// UnlinkPermitEvent removes the link between a permit and an event.
func UnlinkPermitEvent(ctx context.Context, permitID uint, eventID uint, eventType string, authID string) error {
	var user models.User
	DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var action, comment string
		changes := map[string]interface{}{
			"event_id":   eventID,
			"event_type": eventType,
		}

		if eventType == "plate" {
			var event models.PlateEvent
			if err := tx.First(&event, eventID).Error; err != nil {
				return fmt.Errorf("plate event not found: %w", err)
			}
			if event.PermitID == nil || *event.PermitID != permitID {
				return fmt.Errorf("event is not linked to this permit")
			}
			event.PermitID = nil
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
			action = "unlink_plate"
			comment = fmt.Sprintf("Відв'язано подію камери #%d", eventID)
		} else if eventType == "weight" {
			var event models.WeightEvent
			if err := tx.First(&event, eventID).Error; err != nil {
				return fmt.Errorf("weight event not found: %w", err)
			}
			if event.PermitID == nil || *event.PermitID != permitID {
				return fmt.Errorf("event is not linked to this permit")
			}
			event.PermitID = nil
			if err := tx.Save(&event).Error; err != nil {
				return err
			}
			action = "unlink_weight"
			comment = fmt.Sprintf("Відв'язано подію ваг #%d", eventID)
		} else {
			return fmt.Errorf("invalid event type: %s", eventType)
		}

		LogPermitAudit(ctx, permitID, user.ID, action, changes, comment)
		return nil
	})
}

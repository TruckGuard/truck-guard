package repository

import (
	"context"
	"encoding/json"

	"github.com/truckguard/core/src/models"
	"gorm.io/datatypes"
)

// LogSystemAudit logs a general system action.
func LogSystemAudit(ctx context.Context, authID string, action string, target string, changes interface{}, comment string) {
	var user models.User
	if authID != "" {
		DB.WithContext(ctx).Where("auth_id = ?", authID).First(&user)
	}

	jsonBytes, _ := json.Marshal(changes)
	audit := models.SystemAudit{
		Action:  action,
		Target:  target,
		Changes: datatypes.JSON(jsonBytes),
		Comment: comment,
	}

	if user.ID != 0 {
		audit.UserID = &user.ID
	}

	DB.WithContext(ctx).Create(&audit)
}

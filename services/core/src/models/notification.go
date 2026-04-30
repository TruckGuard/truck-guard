package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Notification represents a persistent system notification sent to the UI
type Notification struct {
	gorm.Model
	PostID  uint           `json:"post_id" gorm:"index"`
	Type    string         `json:"type" gorm:"index"` // e.g., "new_permit", "fuzzy_match"
	Message string         `json:"message"`           // Short human-readable summary
	Payload datatypes.JSON `json:"payload" gorm:"type:jsonb"` // Raw SSE event JSON
	IsRead  bool           `json:"is_read" gorm:"default:false;index"`
}

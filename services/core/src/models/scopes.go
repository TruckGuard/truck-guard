package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type PostScopedModel interface {
	GetPostFilterPath() string
	GetResourceName() string
}

func ScopeByPost(ctx context.Context, model PostScopedModel, action string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		perms, _ := ctx.Value("user_permissions").(string)
		postIDStr, _ := ctx.Value("user_post_id").(string)

		resourceName := model.GetResourceName()
		log.Println("ScopeByPost", "perms", perms, "postIDStr", postIDStr, resourceName, action, strings.Contains(perms, fmt.Sprintf("%s:%s:all", action, resourceName)))

		if strings.Contains(perms, fmt.Sprintf("%s:%s:all", action, resourceName)) {
			return db
		}

		if strings.Contains(perms, fmt.Sprintf("%s:%s", action, resourceName)) {
			if postIDStr == "" || postIDStr == "nil" {
				return db.Where("1 = 0")
			}

			path := model.GetPostFilterPath()

			if strings.Contains(path, ".") {
				parts := strings.Split(path, ".")
				relation := parts[0]
				column := parts[1]
				return db.Joins(relation).Where(fmt.Sprintf("\"%s\".%s = ?", relation, column), postIDStr)
			}
			return db.Where(fmt.Sprintf("%s = ?", path), postIDStr)
		}

		return db.Where("1 = 0")
	}
}

func (Permit) GetPostFilterPath() string { return "customs_post_id" }
func (Permit) GetResourceName() string   { return "permits" }

func (PlateEvent) GetPostFilterPath() string { return "Camera.customs_post_id" }
func (PlateEvent) GetResourceName() string   { return "events" }

func (WeightEvent) GetPostFilterPath() string { return "Scale.customs_post_id" }
func (WeightEvent) GetResourceName() string   { return "events" }

func (SystemEvent) GetResourceName() string { return "events" }

func (CameraConfig) GetPostFilterPath() string { return "customs_post_id" }
func (CameraConfig) GetResourceName() string   { return "cameras" }

func (ScaleConfig) GetPostFilterPath() string { return "customs_post_id" }
func (ScaleConfig) GetResourceName() string   { return "scales" }

func (User) GetPostFilterPath() string { return "customs_post_id" }
func (User) GetResourceName() string   { return "users" }

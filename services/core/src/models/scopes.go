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

func HasScopePermission(userPermsStr, requiredAction, requiredResource, requiredScope string) bool {
	if userPermsStr == "" {
		return false
	}
	perms := strings.Split(userPermsStr, ",")
	for _, p := range perms {
		p = strings.TrimSpace(p)
		if p == "admin" {
			return true
		}
		parts := strings.Split(p, ":")
		if len(parts) < 2 {
			continue
		}
		actionUser := parts[0]
		resourceUser := parts[1]
		scopeUser := ""
		if len(parts) > 2 {
			scopeUser = parts[2]
		}

		if resourceUser != "*" && resourceUser != requiredResource {
			continue
		}

		if scopeUser != "all" && scopeUser != requiredScope {
			continue
		}

		if actionUser == "manage" || actionUser == "admin" {
			return true
		}

		if requiredAction == actionUser {
			return true
		}

		switch requiredAction {
		case "read":
			return true
		case "create":
			if actionUser == "update" || actionUser == "delete" {
				return true
			}
		case "update", "validate":
			if actionUser == "delete" {
				return true
			}
		}
	}
	return false
}

func ScopeByPost(ctx context.Context, model PostScopedModel, action string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		perms, _ := ctx.Value("user_permissions").(string)
		postIDStr, _ := ctx.Value("user_post_id").(string)

		resourceName := model.GetResourceName()
		log.Println("ScopeByPost", "perms", perms, "postIDStr", postIDStr, "resource", resourceName, "action", action)

		if HasScopePermission(perms, action, resourceName, "all") {
			return db
		}

		if HasScopePermission(perms, action, resourceName, "") {
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

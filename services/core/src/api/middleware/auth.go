package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/repository"
)

func HasPermission(userPerms []string, required string) bool {
	if required == "" {
		return true
	}

	for _, p := range userPerms {
		if p == "admin" || p == required {
			return true
		}

		partsUser := strings.Split(p, ":")
		partsRequired := strings.Split(required, ":")

		// Check for wildcard permissions, e.g., "read:*" allows "read:users", "read:roles"
		if len(partsUser) == 2 && partsUser[1] == "*" && partsUser[0] == partsRequired[0] {
			return true
		}
	}
	return false
}

func RequireCorePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsHeader := c.GetHeader("X-Permissions")
		if permsHeader == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "Missing permission: " + required})
			return
		}

		userPerms := strings.Split(permsHeader, ",")
		for i, p := range userPerms {
			userPerms[i] = strings.TrimSpace(p)
		}

		if !HasPermission(userPerms, required) {
			c.AbortWithStatusJSON(403, gin.H{"error": "Missing permission: " + required})
			return
		}
		c.Next()
	}
}

func AuthContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		permissions := c.GetHeader("X-Permissions")

		var postIDStr string
		if userID != "" {
			if cachedPostID, err := repository.GetUserPostIDCached(c.Request.Context(), userID); err == nil {
				postIDStr = cachedPostID
			}
		}

		ctx := context.WithValue(c.Request.Context(), "user_post_id", postIDStr)
		ctx = context.WithValue(ctx, "user_permissions", permissions)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

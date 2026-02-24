package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/truckguard/core/src/repository"
)

func RequireCorePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms := c.GetHeader("X-Permissions")
		if perms == "" || !strings.Contains(perms, required) {
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
